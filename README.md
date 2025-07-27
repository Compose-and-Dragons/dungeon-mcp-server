# Dungeon MCP Server
> 🚧 this is a work in progress
## Overview

The **Dungeon MCP Server** is a Model Context Protocol (MCP) server that provides a text-based dungeon crawling adventure system. It allows players to explore mystical dungeons, interact with NPCs, collect treasures, and battle monsters through a set of MCP tools.

## Features

- **Dungeon Exploration**: Navigate through interconnected rooms and corridors
- **Interactive Elements**: NPCs, treasures, and monsters in various locations
- **Player Management**: Load/save player data from YAML files
- **Visual Map Display**: ASCII art representation of the entire dungeon
- **RESTful API**: HTTP server with MCP protocol support
- **Customizable**: Support for custom dungeon and player configurations


## Usage

### Command Line Interface

The server supports several CLI parameters using the Fang library for enhanced user experience:

```bash
./mcp-dungeon [OPTIONS]
```

#### CLI Parameters

| Parameter | Description | Default | Required |
|-----------|-------------|---------|----------|
| `--dungeon-file` | Path to the dungeon YAML file | `crystal_caverns.yaml` | No |
| `--player-file` | Path to the player YAML file | None (uses default player) | No |
| `--port` | HTTP server port | `9090` | No |
| `--generate-player` | Generate a sample player YAML file | `false` | No |
| `--help`, `-h` | Show help information | | No |
| `--version`, `-v` | Show version information | | No |

#### Examples

```bash
# Start server with default settings
./mcp-dungeon

# Start with custom player and port
./mcp-dungeon --player-file my_player.yaml --port 8080

# Use a different dungeon file
./mcp-dungeon --dungeon-file my_dungeon.yaml

# Generate a sample player file
./mcp-dungeon --generate-player --player-file hero.yaml

# Show help
./mcp-dungeon --help
```

### Server Endpoints

Once started, the server provides the following endpoints:

- **MCP Endpoint**: `http://localhost:PORT/mcp` - Main MCP protocol endpoint
- **Health Check**: `http://localhost:PORT/health` - Server health status

## Player Configuration

### Player YAML Format

Look for the player configuration in `templates/player_sample.yaml`. 

### Generating Sample Player

To create a sample player configuration:

```bash
./mcp-dungeon --generate-player
```

This creates `player_sample.yaml` with default values.

## Dungeon Configuration

### Dungeon YAML Format

The dungeon configuration is defined in a YAML file, typically `templates/crystal_caverns.yaml`. 

## MCP Tools

The server provides the following MCP tools:

### 1. say_hello

Say hello to a user.

**Parameters:**
- `name` (string, required): The name of the user to greet

**Example:**
```json
{
  "name": "say_hello",
  "arguments": {
    "name": "Alice"
  }
}
```

### 2. get_room_details_by_name

Get detailed information about a room by its name/ID.

**Parameters:**
- `room_name` (string, required): The name/ID of the room

**Example:**
```json
{
  "name": "get_room_details_by_name",
  "arguments": {
    "room_name": "entrance_cave"
  }
}
```

### 3. get_room_details_by_coordinates

Get detailed information about a room by its coordinates.

**Parameters:**
- `x` (number, required): The X coordinate
- `y` (number, required): The Y coordinate

**Example:**
```json
{
  "name": "get_room_details_by_coordinates",
  "arguments": {
    "x": 2,
    "y": 5
  }
}
```

### 4. move_to_room_by_name

Move the player to a specified room by name. Only allows movement to connected rooms.

**Parameters:**
- `target_room` (string, required): The name/ID of the target room

**Example:**
```json
{
  "name": "move_to_room_by_name",
  "arguments": {
    "target_room": "corridor_1"
  }
}
```

### 5. get_player_status

Get the current status and information of the player.

**Parameters:** None

**Example:**
```json
{
  "name": "get_player_status",
  "arguments": {}
}
```

### 6. display_dungeon_map
> To be implemented
Display an ASCII map of the entire dungeon showing rooms, corridors, and the player's current position.

**Parameters:** None

**Example:**
```json
{
  "name": "display_dungeon_map",
  "arguments": {}
}
```



## Game Mechanics

### Movement Rules

- Players can only move to rooms that are directly connected to their current location
- Movement is validated against the dungeon's connection graph
- Player coordinates are automatically updated when moving

### Room Types

1. **Rooms**: Can contain NPCs, treasures, monsters, or items
2. **Corridors**: Connecting passages between rooms
3. **Entrance**: Starting location for players
4. **Exit**: Goal location to complete the dungeon

### Location Elements

- **NPCs**: Non-player characters with names and descriptions
- **Treasures**: Valuable items with gold values
- **Monsters**: Enemies with difficulty levels and hit points
- **Items**: Consumables like healing potions



### Testing

Test scripts are provided for manual testing:

```bash
cd tests
# Test room details
./room.tool.call.sh

# Test player movement
./move.tool.call.sh

# Test player status
./status.tool.call.sh
```

## Troubleshooting

### Common Issues

1. **Session ID Missing**: Ensure you include the `Mcp-Session-Id` header
2. **Invalid Room Movement**: Check room connections in dungeon YAML
3. **Port Already in Use**: Change the port with `--port` parameter
4. **File Not Found**: Verify file paths for dungeon and player files

### Health Check

Verify the server is running:

```bash
curl http://localhost:9090/health
```

Expected response:
```json
{
  "status": "healthy",
  "dungeon_name": "Dungeon of Crystal Caverns"
}
```
