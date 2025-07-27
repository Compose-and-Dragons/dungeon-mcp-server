package game

import (
	"fmt"
	"strings"

	"mcp-dungeon/models"
)

func GenerateDungeonMap(dungeon *models.Dungeon, player *models.Player) string {
	if dungeon == nil || player == nil {
		return "Dungeon or player data not available"
	}

	width := dungeon.Size.Width
	height := dungeon.Size.Height
	mapGrid := make([][]string, height)
	for i := range mapGrid {
		mapGrid[i] = make([]string, width)
		for j := range mapGrid[i] {
			mapGrid[i][j] = "   "
		}
	}

	for _, location := range dungeon.Locations {
		x, y := location.Coordinates[0], location.Coordinates[1]
		mapY := height - 1 - y
		mapX := x

		var symbol string
		switch location.Type {
		case "room":
			if location.ID == dungeon.EntranceRoom {
				symbol = "[E]"
			} else if location.ID == dungeon.ExitRoom {
				symbol = "[X]"
			} else if location.Monster != nil {
				symbol = "[M]"
			} else if location.NPC != nil {
				symbol = "[N]"
			} else if location.Treasure != nil {
				symbol = "[T]"
			} else {
				symbol = "[R]"
			}
		case "corridor":
			symbol = " - "
		default:
			symbol = " ? "
		}

		if location.ID == player.CurrentLocation {
			symbol = "[P]"
		}

		mapGrid[mapY][mapX] = symbol
	}

	var mapBuilder strings.Builder
	mapBuilder.WriteString(fmt.Sprintf("\n=== %s Map ===\n", dungeon.Name))
	mapBuilder.WriteString(fmt.Sprintf("Size: %dx%d\n\n", width, height))

	mapBuilder.WriteString("    ")
	for x := 0; x < width; x++ {
		mapBuilder.WriteString(fmt.Sprintf("%2d ", x))
	}
	mapBuilder.WriteString("\n")

	for y := 0; y < height; y++ {
		displayY := height - 1 - y
		mapBuilder.WriteString(fmt.Sprintf("%2d  ", displayY))
		for x := 0; x < width; x++ {
			mapBuilder.WriteString(mapGrid[y][x])
		}
		mapBuilder.WriteString("\n")
	}

	mapBuilder.WriteString("\nLegend:\n")
	mapBuilder.WriteString("[P] - Player Position\n")
	mapBuilder.WriteString("[E] - Entrance\n")
	mapBuilder.WriteString("[X] - Exit\n")
	mapBuilder.WriteString("[M] - Monster\n")
	mapBuilder.WriteString("[N] - NPC\n")
	mapBuilder.WriteString("[T] - Treasure\n")
	mapBuilder.WriteString("[R] - Room\n")
	mapBuilder.WriteString(" -  - Corridor\n")

	mapBuilder.WriteString(fmt.Sprintf("\nPlayer: %s at %s [%d, %d]\n",
		player.Name, player.CurrentLocation,
		player.Coordinates[0], player.Coordinates[1]))

	return mapBuilder.String()
}