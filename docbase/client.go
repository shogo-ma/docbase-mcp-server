package docbase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

// SearchPostsResponse は検索結果のレスポンスを表します
type SearchPostsResponse struct {
	Posts []GetPostResponse `json:"posts"`
	Meta  Meta              `json:"meta"`
}

type Meta struct {
	PreviousPage *int `json:"previous_page"`
	NextPage     *int `json:"next_page"`
	Total        int  `json:"total"`
}

// SearchQuery は検索クエリのパラメータを表します
type SearchQuery struct {
	Q       string // 検索クエリ
	Page    int    // ページ番号 (1-indexed)
	PerPage int    // 1ページあたりの結果数
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

func (c *DocBaseClient) SearchPosts(ctx context.Context, query SearchQuery) (*SearchPostsResponse, error) {
	baseURL := fmt.Sprintf("%s/posts", c.BaseURL)
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	q := u.Query()

	if query.Q != "" {
		q.Set("q", query.Q)
	}

	if query.Page > 0 {
		q.Set("page", fmt.Sprintf("%d", query.Page))
	}

	if query.PerPage > 0 {
		q.Set("per_page", fmt.Sprintf("%d", query.PerPage))
	}

	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
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

	var searchResp SearchPostsResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &searchResp, nil
}
