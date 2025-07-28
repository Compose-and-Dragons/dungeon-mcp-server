package handlers

import (
	"slices"
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
)

func MoveToRoomHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	log.Printf("🟢 MoveToRoomHandler called with arguments: %v", args)

	targetRoomValue, exists := args["target_room"]
	if !exists {
		return mcp.NewToolResultText("Missing required parameter: target_room"), nil
	}

	targetRoom, ok := targetRoomValue.(string)
	if !ok {
		return mcp.NewToolResultText("Invalid parameter type: target_room must be a string"), nil
	}

	if CrystalCavernsDungeon == nil {
		return mcp.NewToolResultText("Dungeon data not loaded"), nil
	}

	if CurrentPlayer == nil {
		return mcp.NewToolResultText("Player not initialized"), nil
	}

	targetLocation, exists := CrystalCavernsDungeon.Locations[targetRoom]
	if !exists {
		return mcp.NewToolResultText(fmt.Sprintf("Room '%s' does not exist", targetRoom)), nil
	}

	if CurrentPlayer.CurrentLocation == targetRoom {
		return mcp.NewToolResultText(fmt.Sprintf("Player is already in room '%s'", targetRoom)), nil
	}

	currentLocation, exists := CrystalCavernsDungeon.Locations[CurrentPlayer.CurrentLocation]
	if !exists {
		return mcp.NewToolResultText(fmt.Sprintf("Current player location '%s' is invalid", CurrentPlayer.CurrentLocation)), nil
	}

	connected := slices.Contains(currentLocation.Connections, targetRoom)

	/*
		for _, connection := range currentLocation.Connections {
			if connection == targetRoom {
				connected = true
				break
			}
		}


	*/

	if !connected {
		return mcp.NewToolResultText(fmt.Sprintf("Cannot move to '%s' - not connected to current room '%s'", targetRoom, CurrentPlayer.CurrentLocation)), nil
	}

	CurrentPlayer.CurrentLocation = targetRoom
	CurrentPlayer.Coordinates = targetLocation.Coordinates

	result := fmt.Sprintf("Player %s moved to %s at coordinates [%d, %d]",
		CurrentPlayer.Name, targetRoom, targetLocation.Coordinates[0], targetLocation.Coordinates[1])

	return mcp.NewToolResultText(result), nil
}
