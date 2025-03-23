package docbase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GetPostResponse struct {
	PostID    int64     `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Draft     bool      `json:"draft"`
	Archived  bool      `json:"archived"`
	Tags      []Tag     `json:"tags"`
	User      User      `json:"user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	UserID   int64  `json:"id"`
	UserName string `json:"name"`
}

type Tag struct {
	Name string `json:"name"`
}

type DocBaseClient struct {
	Client   *http.Client
	Domain   string
	APIToken string
	BaseURL  string
}

func NewDocBaseClient(domain, apiToken string) *DocBaseClient {
	return &DocBaseClient{
		Client:   &http.Client{},
		Domain:   domain,
		APIToken: apiToken,
		BaseURL:  fmt.Sprintf("https://api.docbase.io/teams/%s", domain),
	}
}

func (c *DocBaseClient) GetPost(ctx context.Context, postID int64) (*GetPostResponse, error) {
	url := fmt.Sprintf("%s/posts/%d", c.BaseURL, postID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-DocBaseToken", c.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var post GetPostResponse
	if err := json.NewDecoder(resp.Body).Decode(&post); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &post, nil
}
