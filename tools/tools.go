package tools

import "fmt"

// IsInList returns a boolean based on if the instance is in the provided list
func IsInList(instance string, array []string) bool {
	for _, entry := range array {
		if entry == instance {
			return true
		}
	}
	return false
}

// ByteCountDecimal converts a byte count to the human readable form
func ByteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
