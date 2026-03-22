package ui

// Icon set for UI elements. Uses plain Unicode symbols instead of emojis
// for consistent rendering across terminals and a cleaner aesthetic.

// Type icons (displayed next to issue titles)
const (
	IconBug     = "!"  // Bug/fix needed
	IconFeature = "+"  // New feature
	IconTask    = "○"  // Task item
	IconEpic    = "◆"  // Epic/milestone
	IconChore   = "~"  // Maintenance/chore
	IconDefault = "•"  // Unknown type
)

// Status indicators
const (
	IconOpen       = "○"  // Open/unstarted
	IconInProgress = "►"  // In progress
	IconBlocked    = "✗"  // Blocked
	IconDeferred   = "◎"  // Deferred/frozen
	IconPinned     = "◆"  // Pinned
	IconHooked     = "○"  // Hooked to agent
	IconReview     = "?"  // Awaiting review
	IconClosed     = "✓"  // Done/closed
	IconTombstone  = "×"  // Deleted
	IconDraft      = "○"  // Draft (same as open but dimmed)
)

// Triage/recommendation indicators
const (
	IconStar      = "★"  // Top pick
	IconUnblocks  = "→"  // Unblocks other items
	IconAltUnblk  = "↳"  // Alternative unblock indicator
	IconStale     = "⏳" // Stale/no activity
	IconAvailable = "✓"  // Ready to work
	IconWorking   = "►"  // Already being worked
	IconWarning   = "!"  // Warning/attention needed
)

// Priority indicators
const (
	IconCritical = "!!" // P0 - Critical
	IconHigh     = "!"  // P1 - High
	IconMedium   = "-"  // P2 - Medium
	IconLow      = "."  // P3 - Low
	IconBacklog  = "_"  // P4 - Backlog
)

// UI element icons
const (
	IconList    = "›" // List/header
	IconDetail  = "»" // Detail view
	IconCommand = "$" // Command line
	IconInfo    = "i" // Info text
)

// GetTypeIcon returns the icon for an issue type.
func GetTypeIcon(typ string) string {
	switch typ {
	case "bug":
		return IconBug
	case "feature":
		return IconFeature
	case "task":
		return IconTask
	case "epic":
		return IconEpic
	case "chore":
		return IconChore
	default:
		return IconDefault
	}
}
