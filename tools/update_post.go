package tools

import (
	"context"
	"docbase-mcp-server/docbase"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func NewUpdatePostTool() server.ServerTool {
	return server.ServerTool{
		Tool:    newUpdatePostTool(),
		Handler: handleUpdatePostRequest,
	}
}

func newUpdatePostTool() mcp.Tool {
	return mcp.NewTool(
		"update_post",
		mcp.WithDescription("Update an existing post in DocBase"),
		mcp.WithString(
			"post_id",
			mcp.Required(),
			mcp.Description("The ID of the post to update"),
		),
		mcp.WithString(
			"title",
			mcp.Description("The title of the post"),
		),
		mcp.WithString(
			"body",
			mcp.Description("The body content of the post"),
		),
		mcp.WithBoolean(
			"draft",
			mcp.Description("Whether the post is a draft or not"),
		),
		mcp.WithBoolean(
			"notice",
			mcp.Description("Whether to send notification or not"),
		),
		mcp.WithString(
			"tags",
			mcp.Description("Comma-separated list of tags"),
		),
		mcp.WithString(
			"scope",
			mcp.Description("Scope of the post: 'everyone', 'group', or 'private'"),
		),
		mcp.WithString(
			"groups",
			mcp.Description("Comma-separated list of group IDs (required if scope is 'group')"),
		),
	)
}

func handleUpdatePostRequest(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := docbase.NewDocBaseClient(
		os.Getenv("DOCBASE_API_DOMAIN"),
		os.Getenv("DOCBASE_API_TOKEN"),
	)

	// post_idは必須
	postIDStr, ok := request.Params.Arguments["post_id"].(string)
	if !ok || postIDStr == "" {
		return nil, errors.New("post_id is required")
	}

	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		return nil, errors.New("post_id must be a valid number")
	}

	// UpdatePostParamを作成
	updateParam := docbase.UpdatePostParam{}

	// titleが指定されていれば設定
	if title, ok := request.Params.Arguments["title"].(string); ok && title != "" {
		updateParam.Title = title
	}

	// bodyが指定されていれば設定
	if body, ok := request.Params.Arguments["body"].(string); ok && body != "" {
		updateParam.Body = body
	}

	// draftが指定されていれば設定
	if draft, ok := request.Params.Arguments["draft"].(bool); ok {
		updateParam.Draft = &draft
	}

	// noticeが指定されていれば設定
	if notice, ok := request.Params.Arguments["notice"].(bool); ok {
		updateParam.Notice = &notice
	}

	// tagsが指定されていれば設定
	if tagsStr, ok := request.Params.Arguments["tags"].(string); ok && tagsStr != "" {
		tags := []string{}
		for _, tag := range strings.Split(tagsStr, ",") {
			tags = append(tags, strings.TrimSpace(tag))
		}
		updateParam.Tags = tags
	}

	// scopeが指定されていれば設定
	if scopeStr, ok := request.Params.Arguments["scope"].(string); ok && scopeStr != "" {
		var scope docbase.Scope
		switch scopeStr {
		case "everyone":
			scope = docbase.ScopeAll
		case "group":
			scope = docbase.ScopeGroup
		case "private":
			scope = docbase.ScopePrivate
		default:
			return nil, errors.New("scope must be one of 'everyone', 'group', or 'private'")
		}
		updateParam.Scope = scope

		// scopeがgroupの場合はgroupsパラメータが必要
		if scope == docbase.ScopeGroup {
			if groupsStr, ok := request.Params.Arguments["groups"].(string); ok && groupsStr != "" {
				groups := []int{}
				for _, groupStr := range strings.Split(groupsStr, ",") {
					var groupID int
					if _, err := fmt.Sscanf(strings.TrimSpace(groupStr), "%d", &groupID); err == nil {
						groups = append(groups, groupID)
					}
				}
				if len(groups) == 0 {
					return nil, errors.New("at least one valid group ID is required when scope is 'group'")
				}
				updateParam.Groups = groups
			} else {
				return nil, errors.New("groups parameter is required when scope is 'group'")
			}
		}
	}

	// UpdatePost APIを呼び出し
	post, err := client.UpdatePost(ctx, postID, updateParam)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Post updated successfully!\nTitle: %s\nID: %d", post.Title, post.PostID)), nil
}
