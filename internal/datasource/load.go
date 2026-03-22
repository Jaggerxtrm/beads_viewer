package datasource

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/Dicklesworthstone/beads_viewer/pkg/loader"
	"github.com/Dicklesworthstone/beads_viewer/pkg/model"
)

// LoadIssues performs smart multi-source detection and loading.
// It discovers all available sources (Dolt, SQLite, JSONL), validates them, selects
// the best valid source, and loads issues from it. Dolt is preferred over SQLite
// which is preferred over JSONL when both exist at comparable freshness.
//
// Falls back to legacy JSONL-only loading via loader.LoadIssues if smart
// detection finds no valid sources.
func LoadIssues(repoPath string) ([]model.Issue, error) {
	beadsDir, err := loader.GetBeadsDir(repoPath)
	if err != nil {
		return nil, err
	}

	issues, smartErr := loadSmart(beadsDir, repoPath)
	if smartErr == nil {
		return issues, nil
	}

	// Fall back to legacy JSONL-only loading
	return loader.LoadIssues(repoPath)
}

// LoadIssuesFromDir performs smart source detection within a known beads directory.
// This is useful when the caller already knows the .beads path.
func LoadIssuesFromDir(beadsDir string) ([]model.Issue, error) {
	issues, smartErr := loadSmart(beadsDir, "")
	if smartErr == nil {
		return issues, nil
	}

	// Fall back to JSONL
	jsonlPath, err := loader.FindJSONLPath(beadsDir)
	if err != nil {
		return nil, err
	}
	return loader.LoadIssuesFromFile(jsonlPath)
}

// loadSmart discovers sources, validates, selects the best, and loads from it.
func loadSmart(beadsDir, repoPath string) ([]model.Issue, error) {
	sources, err := DiscoverSources(DiscoveryOptions{
		BeadsDir:               beadsDir,
		RepoPath:               repoPath,
		ValidateAfterDiscovery: true,
		IncludeInvalid:         false,
	})
	if err != nil {
		return nil, err
	}
	if len(sources) == 0 {
		return nil, fmt.Errorf("no valid sources discovered")
	}

	// Filter out sources with 0 issues - they're not useful for triage
	var nonEmptySources []DataSource
	for _, s := range sources {
		if s.IssueCount > 0 {
			nonEmptySources = append(nonEmptySources, s)
		}
	}
	
	// If all sources are empty, fall back to JSONL loader
	if len(nonEmptySources) == 0 {
		return nil, fmt.Errorf("no sources with issues found")
	}

	// Select best source with priority over freshness
	// This ensures dolt is preferred, then canonical beads.jsonl/issues.jsonl
	best, err := SelectBestSourceWithOptions(nonEmptySources, SelectionOptions{
		PreferFreshest:      false, // Prefer priority over freshness
		MinimumValidSources: 1,
	})
	if err != nil {
		return nil, err
	}

	return LoadFromSource(best)
}

// LoadFromSource loads issues from a specific DataSource, dispatching to the
// appropriate reader based on source type.
func LoadFromSource(source DataSource) ([]model.Issue, error) {
	switch source.Type {
	case SourceTypeDolt:
		return loadFromDolt(source)
	case SourceTypeSQLite:
		reader, err := NewSQLiteReader(source)
		if err != nil {
			return nil, fmt.Errorf("failed to open SQLite source %s: %w", source.Path, err)
		}
		defer reader.Close()
		return reader.LoadIssues()

	case SourceTypeJSONLLocal, SourceTypeJSONLWorktree:
		return loader.LoadIssuesFromFile(source.Path)

	default:
		return nil, fmt.Errorf("unknown source type: %s", source.Type)
	}
}

// loadFromDolt loads issues from a dolt-backed beads database by running bd export
func loadFromDolt(source DataSource) ([]model.Issue, error) {
	// Find bd command
	bdPath, err := exec.LookPath("bd")
	if err != nil {
		return nil, fmt.Errorf("bd command not found: %w", err)
	}

	// Run bd export to get JSONL output
	cmd := exec.Command(bdPath, "export")
	cmd.Env = append(os.Environ(), fmt.Sprintf("BEADS_DIR=%s", source.BeadsDir))
	
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("bd export failed: %w (stderr: %s)", err, string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("bd export failed: %w", err)
	}

	// Parse the JSONL output using the existing loader
	return loader.ParseIssues(bytes.NewReader(output))
}
