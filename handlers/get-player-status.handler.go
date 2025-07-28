package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
)

func GetPlayerStatusHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("🟢 GetPlayerStatusHandler called")
	if CurrentPlayer == nil {
		return mcp.NewToolResultText("Player not initialized"), nil
	}

	jsonData, err := json.MarshalIndent(CurrentPlayer, "", "  ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error serializing player data: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}
