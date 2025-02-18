package utils

// IsNumeric checks if a string contains only numeric characters.
// Used to distinguish between numeric IDs and usernames.
// 
// Parameters:
//   - s: the string to check
//
// Returns:
//   - bool: true if string contains only digits
//
// Example:
//
//	IsNumeric("123456") // returns true
//	IsNumeric("user123") // returns false
//	IsNumeric("") // returns false
func IsNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
} 