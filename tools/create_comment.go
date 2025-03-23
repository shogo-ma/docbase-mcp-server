package tools

import (
	"context"
	"docbase-mcp-server/docbase"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func NewCreateCommentTool() server.ServerTool {
	return server.ServerTool{
		Tool:    newCreateCommentTool(),
		Handler: handleCreateCommentRequest,
	}
}

func newCreateCommentTool() mcp.Tool {
	return mcp.NewTool(
		"create_comment",
		mcp.WithDescription("Add a comment to a DocBase post"),
		mcp.WithString(
			"post_id",
			mcp.Required(),
			mcp.Description("The ID of the post to comment on"),
		),
		mcp.WithString(
			"body",
			mcp.Required(),
			mcp.Description("The body content of the comment"),
		),
		mcp.WithBoolean(
			"notice",
			mcp.Description("Whether to send notification or not (default is true)"),
		),
	)
}

func handleCreateCommentRequest(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := docbase.NewDocBaseClient(
		os.Getenv("DOCBASE_API_DOMAIN"),
		os.Getenv("DOCBASE_API_TOKEN"),
	)

	// 投稿IDは必須
	postIDStr, ok := request.Params.Arguments["post_id"].(string)
	if !ok || postIDStr == "" {
		return nil, errors.New("post_id is required")
	}

	// 数値に変換
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		return nil, errors.New("post_id must be a valid number")
	}

	// コメント本文は必須
	body, ok := request.Params.Arguments["body"].(string)
	if !ok || body == "" {
		return nil, errors.New("body is required")
	}

	// 通知設定（デフォルトはtrue）
	notice := true
	if noticeParam, ok := request.Params.Arguments["notice"].(bool); ok {
		notice = noticeParam
	}

	// コメント作成パラメータの設定
	commentParam := docbase.CreateCommentParam{
		Body:   body,
		Notice: notice,
	}

	// コメント作成APIの呼び出し
	comment, err := client.CreateComment(ctx, postID, commentParam)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Comment created successfully!\nID: %d\nBody: %s", comment.ID, comment.Body)), nil
}
