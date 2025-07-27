package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	"mcp-dungeon/game"
	"mcp-dungeon/models"
)

var (
	CrystalCavernsDungeon *models.Dungeon
	CurrentPlayer         *models.Player
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

func GetRoomDetailsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

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

func GetRoomDetailsByCoordinatesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

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

func MoveToRoomHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

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

	connected := false
	for _, connection := range currentLocation.Connections {
		if connection == targetRoom {
			connected = true
			break
		}
	}

	if !connected {
		return mcp.NewToolResultText(fmt.Sprintf("Cannot move to '%s' - not connected to current room '%s'", targetRoom, CurrentPlayer.CurrentLocation)), nil
	}

	CurrentPlayer.CurrentLocation = targetRoom
	CurrentPlayer.Coordinates = targetLocation.Coordinates

	result := fmt.Sprintf("Player %s moved to %s at coordinates [%d, %d]",
		CurrentPlayer.Name, targetRoom, targetLocation.Coordinates[0], targetLocation.Coordinates[1])

	return mcp.NewToolResultText(result), nil
}

func GetPlayerStatusHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if CurrentPlayer == nil {
		return mcp.NewToolResultText("Player not initialized"), nil
	}

	jsonData, err := json.MarshalIndent(CurrentPlayer, "", "  ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error serializing player data: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

func DisplayDungeonMapHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if CrystalCavernsDungeon == nil {
		return mcp.NewToolResultText("Dungeon data not loaded"), nil
	}

	if CurrentPlayer == nil {
		return mcp.NewToolResultText("Player not initialized"), nil
	}

	mapString := game.GenerateDungeonMap(CrystalCavernsDungeon, CurrentPlayer)
	return mcp.NewToolResultText(mapString), nil
}