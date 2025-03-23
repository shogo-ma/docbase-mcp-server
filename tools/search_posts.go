package tools

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strconv"

	"docbase-mcp-server/docbase"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func NewSearchPostsTool() (mcp.Tool, server.ToolHandlerFunc) {
	return newSearchPostsTool(), handleSearchPostsRequest
}

func newSearchPostsTool() mcp.Tool {
	return mcp.NewTool(
		"search_posts",
		mcp.WithDescription("Search posts in DocBase by query"),
		mcp.WithString(
			"query",
			mcp.Required(),
			mcp.Description("The query to search for"),
		),
		mcp.WithString(
			"page",
			mcp.Description("The page number (default is 1)"),
		),
		mcp.WithString(
			"per_page",
			mcp.Description("Number of results per page (default is 20, max is 100)"),
		),
	)
}

func handleSearchPostsRequest(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client := docbase.NewDocBaseClient(
		os.Getenv("DOCBASE_API_DOMAIN"),
		os.Getenv("DOCBASE_API_TOKEN"),
	)

	// Get query parameter
	queryStr, ok := request.Params.Arguments["query"]
	if !ok {
		return nil, errors.New("query is required")
	}

	searchQuery := docbase.SearchQuery{
		Q:       queryStr.(string),
		Page:    1,  // Default page is 1
		PerPage: 20, // Default is 20 results per page
	}

	if pageStr, ok := request.Params.Arguments["page"]; ok {
		page, err := strconv.Atoi(pageStr.(string))
		if err != nil {
			return nil, errors.New("page must be a number")
		}
		searchQuery.Page = page
	}

	if perPageStr, ok := request.Params.Arguments["per_page"]; ok {
		perPage, err := strconv.Atoi(perPageStr.(string))
		if err != nil {
			return nil, errors.New("per_page must be a number")
		}
		if perPage > 100 {
			perPage = 100 // API limit is 100
		}
		searchQuery.PerPage = perPage
	}

	result, err := client.SearchPosts(ctx, searchQuery)
	if err != nil {
		return nil, err
	}

	jsonResponse, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(jsonResponse)), nil
}
