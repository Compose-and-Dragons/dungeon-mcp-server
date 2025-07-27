package models

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