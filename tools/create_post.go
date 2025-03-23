package tools

import (
	"context"
	"docbase-mcp-server/docbase"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func NewCreatePostTool() server.ServerTool {
	return server.ServerTool{
		Tool:    newCreatePostTool(),
		Handler: handleCreatePostRequest,
	}
}

func newCreatePostTool() mcp.Tool {
	return mcp.NewTool(
		"create_post",
		mcp.WithDescription("Create a new post in DocBase"),
		mcp.WithString(
			"title",
			mcp.Required(),
			mcp.Description("The title of the post"),
		),
		mcp.WithString(
			"body",
			mcp.Required(),
			mcp.Description("The body content of the post"),
		),
		mcp.WithBoolean(
			"draft",
			mcp.Description("Whether the post is a draft or not (default is false)"),
		),
		mcp.WithBoolean(
			"notice",
			mcp.Description("Whether to send notification or not (default is true)"),
		),
		mcp.WithString(
			"tags",
			mcp.Description("Comma-separated list of tags"),
		),
		mcp.WithString(
			"scope",
			mcp.Description("Scope of the post: 'everyone', 'group', or 'private' (default is 'everyone')"),
		),
		mcp.WithString(
			"groups",
			mcp.Description("Comma-separated list of group IDs (required if scope is 'group')"),
		),
	)
}

func handleCreatePostRequest(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := docbase.NewDocBaseClient(
		os.Getenv("DOCBASE_API_DOMAIN"),
		os.Getenv("DOCBASE_API_TOKEN"),
	)

	title, ok := request.Params.Arguments["title"].(string)
	if !ok || title == "" {
		return nil, errors.New("title is required")
	}

	body, ok := request.Params.Arguments["body"].(string)
	if !ok || body == "" {
		return nil, errors.New("body is required")
	}

	draft, _ := request.Params.Arguments["draft"].(bool)
	notice, _ := request.Params.Arguments["notice"].(bool)

	var tags []string
	if tagsStr, ok := request.Params.Arguments["tags"].(string); ok && tagsStr != "" {
		tags = strings.Split(tagsStr, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
	}

	scopeStr, _ := request.Params.Arguments["scope"].(string)
	var scope docbase.Scope
	switch scopeStr {
	case "everyone":
		scope = docbase.ScopeAll
	case "group":
		scope = docbase.ScopeGroup
	case "private":
		scope = docbase.ScopePrivate
	default: // デフォルトはPrivateにする
		scope = docbase.ScopePrivate
	}

	var groups []int
	if groupsStr, ok := request.Params.Arguments["groups"].(string); ok && groupsStr != "" && scope == docbase.ScopeGroup {
		groupStrs := strings.Split(groupsStr, ",")
		for _, groupStr := range groupStrs {
			var groupID int
			if _, err := fmt.Sscanf(strings.TrimSpace(groupStr), "%d", &groupID); err == nil {
				groups = append(groups, groupID)
			}
		}
	}

	createParam := docbase.CreatePostParam{
		Title:  title,
		Body:   body,
		Draft:  draft,
		Notice: notice,
		Tags:   tags,
		Scope:  scope,
		Groups: groups,
	}

	post, err := client.CreatePost(ctx, createParam)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Post created successfully!\nTitle: %s\nID: %d", post.Title, post.PostID)), nil
}
