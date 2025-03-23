package main

import (
	"docbase-mcp-server/tools"
	"log"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"docbase-mcp-server",
		"0.0.1",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	s.AddTools(
		tools.NewCreatePostTool(),
		tools.NewGetPostTool(),
		tools.NewSearchPostsTool(),
		tools.NewUpdatePostTool(),
	)

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
