package handlers

import (
	"context"
	"log"
	"math/rand"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
)

func RollDicesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	nbDices := request.GetInt("nb_dices", 1)
	sides := request.GetInt("nb_sides", 6)

	log.Printf("🎲 Rolling %d dice(s) with %d sides each...\n", nbDices, sides)

	roll := func(n, x int) int {
		if n <= 0 || x <= 0 {
			return 0
		}

		results := make([]int, n)
		sum := 0

		for i := range n {
			roll := rand.Intn(x) + 1
			results[i] = roll
			sum += roll
		}

		return sum
	}

	// Simulate rolling dice
	result := roll(nbDices, sides)

	return mcp.NewToolResultText("Result: " + strconv.Itoa(nbDices) + " dices with " + strconv.Itoa(sides) + " sides: " + strconv.Itoa(result)), nil

}
