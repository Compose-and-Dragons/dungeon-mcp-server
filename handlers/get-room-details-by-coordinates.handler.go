package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
)

func GetRoomDetailsByCoordinatesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	log.Printf("🟢 GetRoomDetailsByCoordinatesHandler called with arguments: %v", args)

	xValue, exists := args["x"]
	if !exists {
		return mcp.NewToolResultText("Missing required parameter: x"), nil
	}

	xFloat, ok := xValue.(float64)
	if !ok {
		return mcp.NewToolResultText("Invalid parameter type: x must be a number"), nil
	}

	yValue, exists := args["y"]
	if !exists {
		return mcp.NewToolResultText("Missing required parameter: y"), nil
	}

	yFloat, ok := yValue.(float64)
	if !ok {
		return mcp.NewToolResultText("Invalid parameter type: y must be a number"), nil
	}

	x := int(xFloat)
	y := int(yFloat)

	if CrystalCavernsDungeon == nil {
		return mcp.NewToolResultText("Dungeon data not loaded"), nil
	}

	for _, location := range CrystalCavernsDungeon.Locations {
		if location.Coordinates[0] == x && location.Coordinates[1] == y {
			jsonData, err := json.MarshalIndent(location, "", "  ")
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("Error serializing room data: %v", err)), nil
			}
			return mcp.NewToolResultText(string(jsonData)), nil
		}
	}

	return mcp.NewToolResultText(fmt.Sprintf("No room found at coordinates [%d, %d]", x, y)), nil
}
