package tools

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"docbase-mcp-server/docbase"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func NewGetPostTool() server.ServerTool {
	return server.ServerTool{
		Tool:    newGetPostTool(),
		Handler: handleGetPostRequest,
	}
}

func newGetPostTool() mcp.Tool {
	return mcp.NewTool(
		"get_post_by_post_id",
		mcp.WithDescription("Get post from docbase by post ID"),
		mcp.WithString(
			"post_id",
			mcp.Required(),
			mcp.Description("The ID of the post to get"),
		),
	)
}

func handleGetPostRequest(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := docbase.NewDocBaseClient(
		os.Getenv("DOCBASE_API_DOMAIN"),
		os.Getenv("DOCBASE_API_TOKEN"),
	)

	postIDString, ok := request.Params.Arguments["post_id"]
	if !ok {
		return nil, errors.New("post_id is required")
	}

	postID, err := strconv.Atoi(postIDString.(string))
	if err != nil {
		return nil, errors.New("post_id is not a number")
	}

	post, err := client.GetPost(ctx, int64(postID))
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Title: %s\nBody: %s\n", post.Title, post.Body)), nil
}
