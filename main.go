package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/fang"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var crystalCavernsDungeon *Dungeon
var currentPlayer *Player

type Size struct {
	Width  int `yaml:"width"`
	Height int `yaml:"height"`
}

type NPC struct {
	Type        string `yaml:"type"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

type Item struct {
	Type         string `yaml:"type"`
	HealingLevel int    `yaml:"healing_level,omitempty"`
	Quantity     int    `yaml:"quantity"`
}

type Treasure struct {
	Type  string `yaml:"type"`
	Value int    `yaml:"value"`
}

type Monster struct {
	Type            string   `yaml:"type"`
	Name            string   `yaml:"name"`
	Description     string   `yaml:"description"`
	DifficultyLevel int      `yaml:"difficulty_level"`
	HitPoints       int      `yaml:"hit_points"`
	Treasure        Treasure `yaml:"treasure"`
}

type Location struct {
	ID          string    `yaml:"id"`
	Type        string    `yaml:"type"`
	Coordinates [2]int    `yaml:"coordinates"`
	Description string    `yaml:"description"`
	Connections []string  `yaml:"connections"`
	NPC         *NPC      `yaml:"npc,omitempty"`
	Items       []Item    `yaml:"items,omitempty"`
	Treasure    *Treasure `yaml:"treasure,omitempty"`
	Monster     *Monster  `yaml:"monster,omitempty"`
}

type Dungeon struct {
	Name         string              `yaml:"name"`
	Description  string              `yaml:"description"`
	Size         Size                `yaml:"size"`
	EntranceRoom string              `yaml:"entrance_room"`
	ExitRoom     string              `yaml:"exit_room"`
	Locations    map[string]Location `yaml:"locations"`
}

// TODO: remove some fields to simplify
type Player struct {
	Name            string `json:"name" yaml:"name"`
	Avatar          string `json:"avatar" yaml:"avatar"`
	Type            string `json:"type" yaml:"type"`
	Level           int    `json:"level" yaml:"level"`
	HitPoints       int    `json:"hit_points" yaml:"hit_points"`
	MaxHitPoints    int    `json:"max_hit_points" yaml:"max_hit_points"`
	AttackPower     int    `json:"attack_power" yaml:"attack_power"`
	Defense         int    `json:"defense" yaml:"defense"`
	Experience      int    `json:"experience" yaml:"experience"`
	Gold            int    `json:"gold" yaml:"gold"`
	CurrentLocation string `json:"current_location" yaml:"current_location"`
	Coordinates     [2]int `json:"coordinates" yaml:"coordinates"`
	Inventory       []Item `json:"inventory" yaml:"inventory"`
	Status          string `json:"status" yaml:"status"`
}

func loadDungeonFromYAML(filename string) (*Dungeon, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var dungeon Dungeon
	err = yaml.Unmarshal(data, &dungeon)
	if err != nil {
		return nil, err
	}

	return &dungeon, nil
}

func loadPlayerFromYAML(filename string) (*Player, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var player Player
	err = yaml.Unmarshal(data, &player)
	if err != nil {
		return nil, err
	}

	return &player, nil
}

func savePlayerToYAML(player *Player, filename string) error {
	data, err := yaml.Marshal(player)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

func generatePlayerSample(filename string) error {
	samplePlayer := &Player{
		Name:            "Hero",
		Avatar:          "🗡️",
		Type:            "warrior",
		Level:           1,
		HitPoints:       100,
		MaxHitPoints:    100,
		AttackPower:     15,
		Defense:         10,
		Experience:      0,
		Gold:            50,
		CurrentLocation: "entrance_cave",
		Coordinates:     [2]int{0, 0},
		Inventory:       []Item{{Type: "potion", HealingLevel: 25, Quantity: 2}},
		Status:          "healthy",
	}

	return savePlayerToYAML(samplePlayer, filename)
}

var (
	dungeonFile string
	playerFile  string
	port        string
	generate    bool
)

func runServer(cmd *cobra.Command, args []string) error {
	// Generate sample player file if requested
	if generate {
		outputFile := playerFile
		if outputFile == "" {
			outputFile = "player_sample.yaml"
		}
		err := generatePlayerSample(outputFile)
		if err != nil {
			return fmt.Errorf("failed to generate player sample: %v", err)
		}
		log.Printf("Generated sample player file: %s", outputFile)
		return nil
	}

	// Load player from file or create default
	if playerFile != "" {
		var err error
		currentPlayer, err = loadPlayerFromYAML(playerFile)
		if err != nil {
			return fmt.Errorf("failed to load player: %v", err)
		}
		log.Printf("Loaded player: %s", currentPlayer.Name)
	} else {
		currentPlayer = &Player{
			Name:            "Bob",
			Avatar:          "😝",
			Type:            "adventurer",
			CurrentLocation: "entrance_cave",
		}
	}

	// Load dungeon data from YAML file
	var err error
	crystalCavernsDungeon, err = loadDungeonFromYAML(dungeonFile)
	if err != nil {
		return fmt.Errorf("failed to load dungeon: %v", err)
	}

	log.Printf("Loaded dungeon: %s", crystalCavernsDungeon.Name)
	log.Printf("Dungeon size: %dx%d", crystalCavernsDungeon.Size.Width, crystalCavernsDungeon.Size.Height)
	log.Printf("Number of locations: %d", len(crystalCavernsDungeon.Locations))

	// Initialize player coordinates
	if entranceRoom, exists := crystalCavernsDungeon.Locations[currentPlayer.CurrentLocation]; exists {
		currentPlayer.Coordinates = entranceRoom.Coordinates
		log.Printf("Player %s starting at %s [%d, %d]", currentPlayer.Name, currentPlayer.CurrentLocation,
			currentPlayer.Coordinates[0], currentPlayer.Coordinates[1])
	}

	//fmt.Println("Dungeon Details:", crystalCavernsDungeon)
	// Create MCP server
	s := server.NewMCPServer(
		"mcp-dungeon",
		"0.0.0",
	)

	// =================================================
	// TOOLS:
	// =================================================
	sayHello := mcp.NewTool("say_hello",
		mcp.WithDescription(`Say hello to the user.`),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("The name of the user to greet."),
		),
	)
	s.AddTool(sayHello, sayHelloHandler)

	getRoomDetails := mcp.NewTool("get_room_details_by_name",
		mcp.WithDescription(`Get detailed information about a room by its name/ID.`),
		mcp.WithString("room_name",
			mcp.Required(),
			mcp.Description("The name/ID of the room to get details for."),
		),
	)
	s.AddTool(getRoomDetails, getRoomDetailsHandler)

	getRoomByCoords := mcp.NewTool("get_room_details_by_coordinates",
		mcp.WithDescription(`Get detailed information about a room by its coordinates.`),
		mcp.WithNumber("x",
			mcp.Required(),
			mcp.Description("The X coordinate of the room."),
		),
		mcp.WithNumber("y",
			mcp.Required(),
			mcp.Description("The Y coordinate of the room."),
		),
	)
	s.AddTool(getRoomByCoords, getRoomDetailsByCoordinatesHandler)

	moveToRoom := mcp.NewTool("move_to_room_by_name",
		mcp.WithDescription(`Move the player to a specified room by name. Only allows movement to connected rooms.`),
		mcp.WithString("target_room",
			mcp.Required(),
			mcp.Description("The name/ID of the room to move to."),
		),
	)
	s.AddTool(moveToRoom, moveToRoomHandler)

	getPlayerStatus := mcp.NewTool("get_player_status",
		mcp.WithDescription(`Get the current status and information of the player.`),
	)
	s.AddTool(getPlayerStatus, getPlayerStatusHandler)

	displayDungeonMap := mcp.NewTool("display_dungeon_map",
		mcp.WithDescription(`Display an ASCII map of the entire dungeon showing rooms, corridors, and the player's current position.`),
	)
	s.AddTool(displayDungeonMap, displayDungeonMapHandler)

	// Start the HTTP server
	httpPort := port
	if httpPort == "" {
		httpPort = "9090"
	}

	log.Println("MCP StreamableHTTP server is running on port", httpPort)

	// Create a custom mux to handle both MCP and health endpoints
	mux := http.NewServeMux()

	// Add healthcheck endpoint
	mux.HandleFunc("/health", healthCheckHandler)

	// Add MCP endpoint
	httpServer := server.NewStreamableHTTPServer(s,
		server.WithEndpointPath("/mcp"),
	)

	// Register MCP handler with the mux
	mux.Handle("/mcp", httpServer)

	// Start the HTTP server with custom mux
	log.Fatal(http.ListenAndServe(":"+httpPort, mux))
	return nil
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "mcp-dungeon",
		Short: "MCP Dungeon Server",
		Long:  "A Model Context Protocol server for dungeon crawling adventures",
		RunE:  runServer,
	}

	rootCmd.Flags().StringVar(&dungeonFile, "dungeon-file", "crystal_caverns.yaml", "Path to the dungeon YAML file")
	rootCmd.Flags().StringVar(&playerFile, "player-file", "", "Path to the player YAML file")
	rootCmd.Flags().StringVar(&port, "port", "9090", "HTTP server port")
	rootCmd.Flags().BoolVar(&generate, "generate-player", false, "Generate a sample player YAML file")

	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		os.Exit(1)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status":       "healthy",
		"dungeon_name": "Dungeon of Crystal Caverns",
	}
	json.NewEncoder(w).Encode(response)
}

func sayHelloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

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

func getRoomDetailsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	roomNameValue, exists := args["room_name"]
	if !exists {
		return mcp.NewToolResultText("Missing required parameter: room_name"), nil
	}

	roomName, ok := roomNameValue.(string)
	if !ok {
		return mcp.NewToolResultText("Invalid parameter type: room_name must be a string"), nil
	}

	if crystalCavernsDungeon == nil {
		return mcp.NewToolResultText("Dungeon data not loaded"), nil
	}

	location, exists := crystalCavernsDungeon.Locations[roomName]
	if !exists {
		return mcp.NewToolResultText(fmt.Sprintf("Room '%s' not found", roomName)), nil
	}

	jsonData, err := json.MarshalIndent(location, "", "  ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error serializing room data: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

func getRoomDetailsByCoordinatesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	if crystalCavernsDungeon == nil {
		return mcp.NewToolResultText("Dungeon data not loaded"), nil
	}

	for _, location := range crystalCavernsDungeon.Locations {
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

func moveToRoomHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	targetRoomValue, exists := args["target_room"]
	if !exists {
		return mcp.NewToolResultText("Missing required parameter: target_room"), nil
	}

	targetRoom, ok := targetRoomValue.(string)
	if !ok {
		return mcp.NewToolResultText("Invalid parameter type: target_room must be a string"), nil
	}

	if crystalCavernsDungeon == nil {
		return mcp.NewToolResultText("Dungeon data not loaded"), nil
	}

	if currentPlayer == nil {
		return mcp.NewToolResultText("Player not initialized"), nil
	}

	// Check if target room exists
	targetLocation, exists := crystalCavernsDungeon.Locations[targetRoom]
	if !exists {
		return mcp.NewToolResultText(fmt.Sprintf("Room '%s' does not exist", targetRoom)), nil
	}

	// Check if player is already in the target room
	if currentPlayer.CurrentLocation == targetRoom {
		return mcp.NewToolResultText(fmt.Sprintf("Player is already in room '%s'", targetRoom)), nil
	}

	// Get current room to check connections
	currentLocation, exists := crystalCavernsDungeon.Locations[currentPlayer.CurrentLocation]
	if !exists {
		return mcp.NewToolResultText(fmt.Sprintf("Current player location '%s' is invalid", currentPlayer.CurrentLocation)), nil
	}

	// Check if target room is connected to current room
	connected := false
	for _, connection := range currentLocation.Connections {
		if connection == targetRoom {
			connected = true
			break
		}
	}

	if !connected {
		return mcp.NewToolResultText(fmt.Sprintf("Cannot move to '%s' - not connected to current room '%s'", targetRoom, currentPlayer.CurrentLocation)), nil
	}

	// Move the player
	currentPlayer.CurrentLocation = targetRoom
	currentPlayer.Coordinates = targetLocation.Coordinates

	result := fmt.Sprintf("Player %s moved to %s at coordinates [%d, %d]",
		currentPlayer.Name, targetRoom, targetLocation.Coordinates[0], targetLocation.Coordinates[1])

	return mcp.NewToolResultText(result), nil
}

func getPlayerStatusHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if currentPlayer == nil {
		return mcp.NewToolResultText("Player not initialized"), nil
	}

	jsonData, err := json.MarshalIndent(currentPlayer, "", "  ")
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Error serializing player data: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

func displayDungeonMapHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if crystalCavernsDungeon == nil {
		return mcp.NewToolResultText("Dungeon data not loaded"), nil
	}

	if currentPlayer == nil {
		return mcp.NewToolResultText("Player not initialized"), nil
	}

	// Create a 2D grid to represent the dungeon map
	width := crystalCavernsDungeon.Size.Width
	height := crystalCavernsDungeon.Size.Height
	mapGrid := make([][]string, height)
	for i := range mapGrid {
		mapGrid[i] = make([]string, width)
		// Initialize with empty spaces
		for j := range mapGrid[i] {
			mapGrid[i][j] = "   " // 3 spaces for better formatting
		}
	}

	// Place locations on the map
	for _, location := range crystalCavernsDungeon.Locations {
		x, y := location.Coordinates[0], location.Coordinates[1]
		// Convert coordinates to map grid (y-axis is flipped for display)
		mapY := height - 1 - y
		mapX := x

		var symbol string
		switch location.Type {
		case "room":
			if location.ID == crystalCavernsDungeon.EntranceRoom {
				symbol = "[E]" // Entrance
			} else if location.ID == crystalCavernsDungeon.ExitRoom {
				symbol = "[X]" // Exit
			} else if location.Monster != nil {
				symbol = "[M]" // Monster
			} else if location.NPC != nil {
				symbol = "[N]" // NPC
			} else if location.Treasure != nil {
				symbol = "[T]" // Treasure
			} else {
				symbol = "[R]" // Regular room
			}
		case "corridor":
			symbol = " - " // Corridor
		default:
			symbol = " ? " // Unknown
		}

		// Mark player position
		if location.ID == currentPlayer.CurrentLocation {
			symbol = "[P]" // Player
		}

		mapGrid[mapY][mapX] = symbol
	}

	// Build the map string
	var mapBuilder strings.Builder
	mapBuilder.WriteString(fmt.Sprintf("\n=== %s Map ===\n", crystalCavernsDungeon.Name))
	mapBuilder.WriteString(fmt.Sprintf("Size: %dx%d\n\n", width, height))

	// Add coordinate headers
	mapBuilder.WriteString("    ")
	for x := 0; x < width; x++ {
		mapBuilder.WriteString(fmt.Sprintf("%2d ", x))
	}
	mapBuilder.WriteString("\n")

	// Draw the map (from top to bottom, but label with correct coordinates)
	for y := 0; y < height; y++ {
		displayY := height - 1 - y // Y coordinate to display (5 at top, 0 at bottom)
		mapBuilder.WriteString(fmt.Sprintf("%2d  ", displayY))
		for x := 0; x < width; x++ {
			mapBuilder.WriteString(mapGrid[y][x])
		}
		mapBuilder.WriteString("\n")
	}

	// Add legend
	mapBuilder.WriteString("\nLegend:\n")
	mapBuilder.WriteString("[P] - Player Position\n")
	mapBuilder.WriteString("[E] - Entrance\n")
	mapBuilder.WriteString("[X] - Exit\n")
	mapBuilder.WriteString("[M] - Monster\n")
	mapBuilder.WriteString("[N] - NPC\n")
	mapBuilder.WriteString("[T] - Treasure\n")
	mapBuilder.WriteString("[R] - Room\n")
	mapBuilder.WriteString(" -  - Corridor\n")

	// Add current player info
	mapBuilder.WriteString(fmt.Sprintf("\nPlayer: %s at %s [%d, %d]\n",
		currentPlayer.Name, currentPlayer.CurrentLocation,
		currentPlayer.Coordinates[0], currentPlayer.Coordinates[1]))

	return mcp.NewToolResultText(mapBuilder.String()), nil
}
