package models

import (
	http "github.com/bogdanfinn/fhttp"
)

// Account represents a Twitter account with all necessary credentials and information
type Account struct {
	Ct0       string // CSRF token
	AuthToken string // auth_token cookie
	Proxy     string // Format: "ip:port" or "user:pass@ip:port"

	// Account Info
	Username      string
	DisplayName   string
	UserID        string
	Email         string
	PhoneNumber   string
	IsVerified    bool
	CreatedAt     string
	FollowCount   int
	FollowerCount int

	// Session Data
	Cookies []*http.Cookie

	// Optional fields
	ProfileImageURL string
	Bio             string
	Location        string
	Website         string

	Suspended bool
	Locked    bool
}
