# docbase-mcp-server

A Model Context Protocol(MCP) server for Docbase.

## Features

- Search Posts
- Get Post
- Create Post

## Usage

```
$ git clone https://github.com/shogo-ma/docbase-mcp-server.git
$ cd docbase-mcp-server
$ go build -o docbase-mcp-server
```

Please write in .cursor/mcp.json as follows:

```
{
    "mcpServers": {
        "docbase-mcp-server": {
            "command": "your build path"
        },
        "env": {
            "DOCBASE_API_DOMAIN": "your docbase domain",
            "DOCBASE_API_TOKEN": "your api key"
        }
    }
}
```
