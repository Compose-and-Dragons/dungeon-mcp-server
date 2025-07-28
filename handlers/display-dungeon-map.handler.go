package handlers

import (
	"context"
	"log"
	"mcp-dungeon/game"

	"github.com/mark3labs/mcp-go/mcp"
)

func DisplayDungeonMapHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("🟢 DisplayDungeonMapHandler called")
	if CrystalCavernsDungeon == nil {
		return mcp.NewToolResultText("Dungeon data not loaded"), nil
	}

	if CurrentPlayer == nil {
		return mcp.NewToolResultText("Player not initialized"), nil
	}

	mapString := game.GenerateDungeonMap(CrystalCavernsDungeon, CurrentPlayer)
	return mcp.NewToolResultText(mapString), nil
}
