package models

// Tweet represents a Twitter tweet with its basic information
type Tweet struct {
	ID              string
	AuthorUsername  string
	AuthorID        string
	Text            string
	CreatedAt       string
	LikeCount       int
	RetweetCount    int
	QuoteCount      int
	ReplyCount      int
	IsLiked         bool
	IsRetweeted     bool
	IsQuoted        bool
	IsReply         bool
	ConversationID  string
	InReplyToUserID string
}

// TweetGraphQLResponse represents the GraphQL response for a tweet action
type TweetGraphQLResponse struct {
	Data struct {
		CreateTweet struct {
			TweetResults struct {
				Result struct {
					RestID string `json:"rest_id"`
				} `json:"result"`
			} `json:"tweet_results"`
		} `json:"create_tweet"`
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

// MediaUploadResponse represents the response from media upload
type MediaUploadResponse struct {
	MediaIDString string `json:"media_id_string"`
	Size          int    `json:"size"`
	ExpiresAfter  int    `json:"expires_after_secs"`
}

// RetweetGraphQLResponse represents the GraphQL response for a retweet action
type RetweetGraphQLResponse struct {
	Data struct {
		CreateRetweet struct {
			RetweetResults struct {
				Result struct {
					RestID string `json:"rest_id"`
					Legacy struct {
						FullText string `json:"full_text"`
					} `json:"legacy"`
				} `json:"result"`
			} `json:"retweet_results"`
		} `json:"create_retweet"`
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

// UnretweetGraphQLResponse represents the GraphQL response for an unretweet action
type UnretweetGraphQLResponse struct {
	Data struct {
		Unretweet struct {
			SourceTweetResults struct {
				Result struct {
					RestID string `json:"rest_id"`
					Legacy struct {
						FullText string `json:"full_text"`
					} `json:"legacy"`
				} `json:"result"`
			} `json:"source_tweet_results"`
		} `json:"unretweet"`
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
