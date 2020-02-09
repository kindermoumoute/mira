package mira

import (
	"net/http"
	"time"
)

type Reddit struct {
	Token    string  `json:"access_token"`
	Duration float64 `json:"expires_in"`
	Creds    Credentials
	Chain    []ChainVals
	Stream   Streaming
	Values   RedditVals
	Client   *http.Client
}

type Streaming struct {
	CommentListInterval         time.Duration
	CommentListNoResultInterval time.Duration
	PostListInterval            time.Duration
	PostListNoResultInterva     time.Duration
	PostListSlice               int
}

type RedditVals struct {
	GetSubmissionFromCommentTries int
}

type ChainVals struct {
	Name string
	Type string
}
