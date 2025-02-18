package models

// LikeStatus represents the status of a like action
type LikeStatus int

const (
	LikeStatusSuccess LikeStatus = iota
	LikeStatusAlreadyLiked
	LikeStatusTweetNotFound
	LikeStatusRateLimited
	LikeStatusAuthError
	LikeStatusUnknown
)

// LikeResponse represents the response from a like action
type LikeResponse struct {
	Success bool
	Error   error
	Status  LikeStatus
}

// LikeGraphQLResponse represents the GraphQL response for a like action
type LikeGraphQLResponse struct {
	Data struct {
		FavoriteTweet string `json:"favorite_tweet"`
	} `json:"data"`
	Errors []struct {
		Message   string `json:"message"`
		Locations []struct {
			Line   int `json:"line"`
			Column int `json:"column"`
		} `json:"locations"`
		Path       []string `json:"path"`
		Extensions struct {
			Name    string `json:"name"`
			Source  string `json:"source"`
			Code    int    `json:"code"`
			Kind    string `json:"kind"`
			Tracing struct {
				TraceID string `json:"trace_id"`
			} `json:"tracing"`
		} `json:"extensions"`
	} `json:"errors"`
}

// UnlikeGraphQLResponse represents the GraphQL response for an unlike action
type UnlikeGraphQLResponse struct {
	Data struct {
		UnfavoriteTweet string `json:"unfavorite_tweet"`
	} `json:"data"`
	Errors []struct {
		Message   string `json:"message"`
		Locations []struct {
			Line   int `json:"line"`
			Column int `json:"column"`
		} `json:"locations"`
		Path       []string `json:"path"`
		Extensions struct {
			Name    string `json:"name"`
			Source  string `json:"source"`
			Code    int    `json:"code"`
			Kind    string `json:"kind"`
			Tracing struct {
				TraceID string `json:"trace_id"`
			} `json:"tracing"`
		} `json:"extensions"`
	} `json:"errors"`
}

// AlreadyLikedResponse represents the response when tweet is already liked
type AlreadyLikedResponse struct {
	Errors []struct {
		Message   string `json:"message"`
		Locations []struct {
			Line   int `json:"line"`
			Column int `json:"column"`
		} `json:"locations"`
		Path       []string `json:"path"`
		Extensions struct {
			Name    string `json:"name"`
			Source  string `json:"source"`
			Code    int    `json:"code"`
			Kind    string `json:"kind"`
			Tracing struct {
				TraceID string `json:"trace_id"`
			} `json:"tracing"`
		} `json:"extensions"`
	} `json:"errors"`
}
