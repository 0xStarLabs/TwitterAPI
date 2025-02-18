package addons

import (
	"fmt"
	"strings"

	"github.com/0xStarLabs/TwitterAPI/utils"
)

// ExtractTweetID extracts the numeric tweet ID from a tweet URL or ID string.
// 
// Parameters:
//   - tweetLink: can be a full tweet URL or just the ID
//   - username: account username for logging purposes
//   - logger: logger instance for error reporting
//
// Returns:
//   - string: the extracted tweet ID
//   - error: any error that occurred during extraction
//
// Example:
//
//	// Extract from URL
//	id, err := ExtractTweetID("https://twitter.com/user/status/1234567890", username, logger)
//	
//	// Extract from ID string
//	id, err := ExtractTweetID("1234567890", username, logger)
func ExtractTweetID(tweetLink string, username string, logger utils.Logger) (string, error) {
	tweetLink = strings.TrimSpace(tweetLink)

	var tweetID string
	if strings.Contains(tweetLink, "tweet_id=") {
		parts := strings.Split(tweetLink, "tweet_id=")
		tweetID = strings.Split(parts[1], "&")[0]
	} else if strings.Contains(tweetLink, "?") {
		parts := strings.Split(tweetLink, "status/")
		tweetID = strings.Split(parts[1], "?")[0]
	} else if strings.Contains(tweetLink, "status/") {
		parts := strings.Split(tweetLink, "status/")
		tweetID = parts[1]
	} else {
		logger.Error("%s | Failed to get tweet ID from your link: %s", username, tweetLink)
		return "", fmt.Errorf("failed to get tweet ID from your link: %s", tweetLink)
	}

	return tweetID, nil
}