package models

import (
	"time"

	"github.com/0xStarLabs/TwitterAPI/utils"
)

// TwitterConstants holds all Twitter-specific constants
type TwitterConstants struct {
	// API Constants
	BearerToken string
	UserAgent   string

	// Query IDs
	QueryID struct {
		Like    string
		Unlike    string
		Retweet   string
		Unretweet string
		Tweet     string
	}
}

// Config holds Twitter client configuration
type Config struct {
	// HTTP Client settings
	MaxRetries      int
	Timeout         time.Duration
	FollowRedirects bool

	// Logging options
	LogLevel utils.LogLevel // Level of logging detail

	// Twitter Constants
	Constants TwitterConstants
}

// NewConfig returns a Config with default settings
func NewConfig() *Config {
	return &Config{
		MaxRetries:      3,
		Timeout:         30 * time.Second,
		FollowRedirects: true,
		LogLevel:        utils.LogLevelError, // By default, only log errors
		Constants: TwitterConstants{
			UserAgent:   UserAgent,
			BearerToken: BearerToken,
			QueryID: struct {
				Like    string
				Unlike  string
				Retweet string
				Unretweet string
				Tweet   string
			}{
				Like:      QueryIDLike,
				Unlike:    QueryIDUnlike,
				Retweet:   QueryIDRetweet,
				Unretweet: QueryIDUnretweet,
				Tweet:     QueryIDTweet,
			},
		},
	}
}
