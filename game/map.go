package game

import (
	"fmt"
	"mcp-dungeon/models"
)

func GenerateVisualMap(dungeon *models.Dungeon, player *models.Player) string {
	grid := make([][]string, dungeon.Size.Height)
	for i := range grid {
		grid[i] = make([]string, dungeon.Size.Width)
		for j := range grid[i] {
			grid[i][j] = " . "
		}
	}

	for _, location := range dungeon.Locations {
		x, y := location.Coordinates[0], location.Coordinates[1]
		if x >= 0 && x < dungeon.Size.Width && y >= 0 && y < dungeon.Size.Height {
			symbol := " . "
			switch {
			case location.ID == dungeon.EntranceRoom:
				symbol = "[E]"
			case location.ID == dungeon.ExitRoom:
				symbol = "[X]"
			case location.Type == "room":
				symbol = "[R]"
			case location.Type == "corridor":
				symbol = "[C]"
			}

			if player != nil && location.ID == player.CurrentLocation {
				switch symbol {
				case "[E]":
					symbol = "{E}"
				case "[X]":
					symbol = "{X}"
				case "[R]":
					symbol = "{R}"
				case "[C]":
					symbol = "{C}"
				}
			}

			grid[y][x] = symbol
		}
	}

	var result string
	result += "## Visual Map\n\n```\n"
	result += "   "
	for x := 0; x < dungeon.Size.Width; x++ {
		result += fmt.Sprintf("%d  ", x)
	}
	result += "\n"

	for y := 0; y < dungeon.Size.Height; y++ {
		result += fmt.Sprintf("%d  ", y)
		for x := 0; x < dungeon.Size.Width; x++ {
			result += grid[y][x]
		}
		result += "\n"
	}
	result += "```\n\n"

	result += "## Legend\n\n"
	result += "- [R] = Room [C] = Corridor [E] = Entrance [X] = Exit\n"
	result += "- {P} = Player position\n\n"

	result += fmt.Sprintf("Player: %s\n", player.Name)
	result += fmt.Sprintf("Current Location: %s Coordinates: [%d, %d]\n", player.CurrentLocation, player.Coordinates[0], player.Coordinates[1])
	result += fmt.Sprintf("Connections: %v\n", dungeon.Locations[player.CurrentLocation].Connections)

	return result
}

func GenerateDungeonMap(dungeon *models.Dungeon, player *models.Player) string {
	var report string
	
	report += "Dungeon: " + dungeon.Name + "\n"
	report += "Size: " + fmt.Sprintf("%dx%d", dungeon.Size.Width, dungeon.Size.Height) + "\n"
	report += "Entrance Room: " + dungeon.EntranceRoom + "\n"
	report += "Exit Room: " + dungeon.ExitRoom + "\n\n"

	report += GenerateVisualMap(dungeon, player) + "\n"

	// report += "=== ROOMS ===\n"
	// for id, location := range dungeon.Locations {
	// 	report += fmt.Sprintf("Room ID: %s\n", id)
	// 	report += fmt.Sprintf("  Type: %s\n", location.Type)
	// 	report += fmt.Sprintf("  Coordinates: [%d, %d]\n", location.Coordinates[0], location.Coordinates[1])
	// 	report += fmt.Sprintf("  Description: %s\n", location.Description)

	// 	if len(location.Connections) > 0 {
	// 		report += "  Connections: "
	// 		for i, conn := range location.Connections {
	// 			if i > 0 {
	// 				report += ", "
	// 			}
	// 			report += conn
	// 		}
	// 		report += "\n"
	// 	}
	// 	report += "\n"
	// }

	return report
}
