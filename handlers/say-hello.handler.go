package handlers

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

func SayHelloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	nameValue, exists := args["name"]
	if !exists {
		return mcp.NewToolResultText("Missing required parameter: name"), nil
	}

	name, ok := nameValue.(string)
	if !ok {
		return mcp.NewToolResultText("Invalid parameter type: name must be a string"), nil
	}

	return mcp.NewToolResultText("👋 Hello " + name + " 🙂"), nil
}
