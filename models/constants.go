package models

import "errors"

// User Agent
const (
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
)

// API Constants
const (
	BearerToken = "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA"
)

// Query IDs for different operations
const (
	QueryIDLike      = "lI07N6Otwv1PhnEgXILM7A"
	QueryIDUnlike    = "ZYKSe-w7KEslx3JhSIk5LA"
	QueryIDRetweet   = "ojPdsZsimiJrUGLR1sjUtA"
	QueryIDUnretweet = "iQtK4dl5hBmXewYZLkNG9A"
	QueryIDTweet     = "bDE2rBtZb3uyrczSZ_pI9g"
)

// Common error types for Twitter operations
var (
	ErrAccountLocked = errors.New("account is temporarily locked")
	ErrAuthFailed    = errors.New("authentication failed")
	ErrInvalidToken  = errors.New("invalid token")
	ErrUnknown       = errors.New("unable to complete operation")
)

// ActionStatus represents the status of any Twitter action (like, retweet, etc.)
type ActionStatus int

const (
	StatusSuccess     ActionStatus = iota
	StatusAlreadyDone              // Already liked, already retweeted, etc.
	StatusLocked                   // Account is locked
	StatusNotFound                 // Tweet/User not found
	StatusRateLimited              // Rate limit exceeded
	StatusAuthError                // Authentication error
	StatusInvalidToken             // Invalid token
	StatusUnknown                  // Unknown error
)

// ActionResponse represents the response from any Twitter action
type ActionResponse struct {
	Success bool
	Error   error
	Status  ActionStatus
}
