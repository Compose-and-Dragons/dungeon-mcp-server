package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"

	"mcp-dungeon/handlers"
	"mcp-dungeon/models"
	myserver "mcp-dungeon/server"
	"mcp-dungeon/storage"
)

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
		err := storage.GeneratePlayerSample(outputFile)
		if err != nil {
			return fmt.Errorf("failed to generate player sample: %v", err)
		}
		log.Printf("Generated sample player file: %s", outputFile)
		return nil
	}

	// Load player from file or create default
	if playerFile != "" {
		var err error
		handlers.CurrentPlayer, err = storage.LoadPlayerFromYAML(playerFile)
		if err != nil {
			return fmt.Errorf("failed to load player: %v", err)
		}
		log.Printf("Loaded player: %s", handlers.CurrentPlayer.Name)
	} else {
		handlers.CurrentPlayer = &models.Player{
			Name:            "Bob",
			Avatar:          "😝",
			Type:            "adventurer",
			CurrentLocation: "entrance_cave",
		}
	}

	// Load dungeon data from YAML file
	var err error
	handlers.CrystalCavernsDungeon, err = storage.LoadDungeonFromYAML(dungeonFile)
	if err != nil {
		return fmt.Errorf("failed to load dungeon: %v", err)
	}

	log.Printf("Loaded dungeon: %s", handlers.CrystalCavernsDungeon.Name)
	log.Printf("Dungeon size: %dx%d", handlers.CrystalCavernsDungeon.Size.Width, handlers.CrystalCavernsDungeon.Size.Height)
	log.Printf("Number of locations: %d", len(handlers.CrystalCavernsDungeon.Locations))

	// Initialize player coordinates
	if entranceRoom, exists := handlers.CrystalCavernsDungeon.Locations[handlers.CurrentPlayer.CurrentLocation]; exists {
		handlers.CurrentPlayer.Coordinates = entranceRoom.Coordinates
		log.Printf("Player %s starting at %s [%d, %d]", handlers.CurrentPlayer.Name, handlers.CurrentPlayer.CurrentLocation,
			handlers.CurrentPlayer.Coordinates[0], handlers.CurrentPlayer.Coordinates[1])
	}

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
	s.AddTool(sayHello, handlers.SayHelloHandler)

	rollDices := mcp.NewTool("rool_dices",
		mcp.WithDescription("Roll some dices"),
		mcp.WithNumber("nb_dices",
			mcp.Required(),
			mcp.Description("Number of dices to roll"),
		),
		mcp.WithNumber("nb_sides",
			mcp.Required(),
			mcp.Description("Number of sides of the dices"),
		),
	)

	s.AddTool(rollDices, handlers.RollDicesHandler)

	getRoomDetails := mcp.NewTool("get_room_details_by_name",
		mcp.WithDescription(`Get detailed information about a room by its name/ID.`),
		mcp.WithString("room_name",
			mcp.Required(),
			mcp.Description("The name/ID of the room to get details for."),
		),
	)
	s.AddTool(getRoomDetails, handlers.GetRoomDetailsHandler)

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
	s.AddTool(getRoomByCoords, handlers.GetRoomDetailsByCoordinatesHandler)

	moveToRoom := mcp.NewTool("move_to_room_by_name",
		mcp.WithDescription(`Move the player to a specified room by name. Only allows movement to connected rooms.`),
		mcp.WithString("target_room",
			mcp.Required(),
			mcp.Description("The name/ID of the room to move to."),
		),
	)
	s.AddTool(moveToRoom, handlers.MoveToRoomHandler)

	getPlayerStatus := mcp.NewTool("get_player_status",
		mcp.WithDescription(`Get the current status and information of the player.`),
	)
	s.AddTool(getPlayerStatus, handlers.GetPlayerStatusHandler)

	displayDungeonMap := mcp.NewTool("display_dungeon_map",
		mcp.WithDescription(`Display an ASCII map of the entire dungeon showing rooms, corridors, and the player's current position.`),
	)
	s.AddTool(displayDungeonMap, handlers.DisplayDungeonMapHandler)

	// Start the HTTP server
	httpPort := port
	if httpPort == "" {
		httpPort = "9090"
	}

	log.Println("MCP StreamableHTTP server is running on port", httpPort)

	// Create a custom mux to handle both MCP and health endpoints
	mux := http.NewServeMux()

	// Add healthcheck endpoint
	mux.HandleFunc("/health", myserver.HealthCheckHandler)

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
