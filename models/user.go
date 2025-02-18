package models

// UserResponse represents the response from Twitter's user-related endpoints
type UserResponse struct {
	ID                int64  `json:"id"`
	IDStr             string `json:"id_str"`
	Name              string `json:"name"`
	ScreenName        string `json:"screen_name"`
	Location          string `json:"location"`
	Description       string `json:"description"`
	URL               any    `json:"url"`
	Protected         bool   `json:"protected"`
	FollowersCount    int    `json:"followers_count"`
	FriendsCount      int    `json:"friends_count"`
	ListedCount       int    `json:"listed_count"`
	CreatedAt         string `json:"created_at"`
	FavouritesCount   int    `json:"favourites_count"`
	Verified          bool   `json:"verified"`
	StatusesCount     int    `json:"statuses_count"`
	MediaCount        int    `json:"media_count"`
	Following         bool   `json:"following"`
	FollowRequestSent bool   `json:"follow_request_sent"`
	Notifications     bool   `json:"notifications"`
	Entities          struct {
		Description struct {
			Urls []any `json:"urls"`
		} `json:"description"`
	} `json:"entities"`
}

// AccountInfoResponse represents the response from Twitter's user lookup endpoint
type AccountInfoResponse struct {
	ID                int64  `json:"id"`
	IDStr             string `json:"id_str"`
	Name              string `json:"name"`
	ScreenName        string `json:"screen_name"`
	Location          string `json:"location"`
	Description       string `json:"description"`
	FollowersCount    int    `json:"followers_count"`
	FriendsCount      int    `json:"friends_count"`
	ListedCount       int    `json:"listed_count"`
	CreatedAt         string `json:"created_at"`
	FavouritesCount   int    `json:"favourites_count"`
	StatusesCount     int    `json:"statuses_count"`
	MediaCount        int    `json:"media_count"`
	Protected         bool   `json:"protected"`
	Verified          bool   `json:"verified"`
	Suspended         bool   `json:"suspended"`
	Following         bool   `json:"following"`
	FollowedBy        bool   `json:"followed_by"`
	FollowRequestSent bool   `json:"follow_request_sent"`
}
