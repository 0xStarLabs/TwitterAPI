package twitter_utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

/*
	GenerateCSRFToken generates a new CSRF token for Twitter requests.
	The token is a 32-character hexadecimal string, matching Twitter's format.

Example:
	token, err := GenerateCSRFToken()
	if err != nil {
		// handle error
	}
	// token = "1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p"

Returns:
  - string: 32-character hexadecimal CSRF token
  - error: if token generation fails
*/
func GenerateCSRFToken() (string, error) {
	// Create a byte array of length 16 (same as Uint8Array(16) in JS)
	bytes := make([]byte, 16)

	// Fill bytes with random values (equivalent to crypto.getRandomValues in JS)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Convert bytes to hex string (equivalent to toString(16) in JS)
	// hex.EncodeToString automatically handles padding, similar to padStart(2, '0')
	token := hex.EncodeToString(bytes)

	return token, nil
}
