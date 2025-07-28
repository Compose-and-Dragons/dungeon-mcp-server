package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
)

func GetRoomDetailsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	log.Printf("🟢 GetRoomDetailsHandler called with arguments: %v", args)

	roomNameValue, exists := args["room_name"]
	if !exists {
		return mcp.NewToolResultText("Missing required parameter: room_name"), nil
	}

	roomName, ok := roomNameValue.(string)
	if !ok {
		return mcp.NewToolResultText("Invalid parameter type: room_name must be a string"), nil
	}

	if CrystalCavernsDungeon == nil {
		return mcp.NewToolResultText("Dungeon data not loaded"), nil
	}

	location, exists := CrystalCavernsDungeon.Locations[roomName]
	if !exists {
		return mcp.NewToolResultText(fmt.Sprintf("Room '%s' not found", roomName)), nil
	}

	jsonData, err := json.MarshalIndent(location, "", "  ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error serializing room data: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}
