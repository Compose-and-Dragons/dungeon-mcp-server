package storage

import (
	"os"

	"gopkg.in/yaml.v3"

	"mcp-dungeon/models"
)

func LoadDungeonFromYAML(filename string) (*models.Dungeon, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var dungeon models.Dungeon
	err = yaml.Unmarshal(data, &dungeon)
	if err != nil {
		return nil, err
	}

	return &dungeon, nil
}

func LoadPlayerFromYAML(filename string) (*models.Player, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var player models.Player
	err = yaml.Unmarshal(data, &player)
	if err != nil {
		return nil, err
	}

	return &player, nil
}

func SavePlayerToYAML(player *models.Player, filename string) error {
	data, err := yaml.Marshal(player)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

func GeneratePlayerSample(filename string) error {
	samplePlayer := &models.Player{
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
		Inventory:       []models.Item{{Type: "potion", HealingLevel: 25, Quantity: 2}},
		Status:          "healthy",
	}

	return SavePlayerToYAML(samplePlayer, filename)
}