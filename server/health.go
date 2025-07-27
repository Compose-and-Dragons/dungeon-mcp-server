package server

import (
	"encoding/json"
	"net/http"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := map[string]interface{}{
		"status":       "healthy",
		"dungeon_name": "Dungeon of Crystal Caverns",
	}
	json.NewEncoder(w).Encode(response)
}