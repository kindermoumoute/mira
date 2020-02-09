package mira

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Credentials struct {
	ClientId     string
	ClientSecret string
	Username     string
	Password     string
	UserAgent    string
}

// Returns an access_token acquired using the provided credentials
func Authenticate(c *Credentials) (*Reddit, error) {
	// URL to get access_token
	auth_url := RedditBase + "api/v1/access_token"

	// Define the data to send in the request
	form := url.Values{}
	form.Add("grant_type", "password")
	form.Add("username", c.Username)
	form.Add("password", c.Password)

	// Encode the Authorization Header
	raw := c.ClientId + ":" + c.ClientSecret
	encoded := b64.StdEncoding.EncodeToString([]byte(raw))

	// Create a request to allow customised headers
	r, err := http.NewRequest("POST", auth_url, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	// Customise request headers
	r.Header.Set("User-Agent", c.UserAgent)
	r.Header.Set("Authorization", "Basic "+encoded)

	// Create client
	client := &http.Client{}

	// Run the request
	response, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)

	data := buf.Bytes()
	if err := findRedditError(data); err != nil {
		return nil, err
	}

	auth := Reddit{}
	json.Unmarshal(data, &auth)
	auth.Creds = *c
	return &auth, nil
}

// This goroutine reauthenticates the user
// every 45 minutes. It should be run with the go
// statement
func (c *Reddit) auto_refresh() {
	for {
		time.Sleep(45 * time.Minute)
		c.update_creds()
	}
}

// Reauthenticate and updates the object itself
func (c *Reddit) update_creds() {
	temp, _ := Authenticate(&c.Creds)
	// Just updated the token
	c.Token = temp.Token
}

func (c *Reddit) SetDefault() {
	c.Stream = Streaming{
		CommentListInterval:         8 * time.Second,
		CommentListNoResultInterval: 8 * time.Second,
		PostListInterval:            10 * time.Second,
		PostListNoResultInterval:    10 * time.Second,
		PostListSlice:               25,
	}
	c.Values = RedditVals{
		GetSubmissionFromCommentTries: 32,
	}
}

func (c *Reddit) SetClient(client *http.Client) {
	c.Client = client
}
