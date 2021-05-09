package main

import (
	"github.com/labstack/echo/v4"
	"github.com/pusher/pusher-http-go"
)

func gameStart(c echo.Context) error {
	player := getPlayerFromToken(c)
	opponent := c.FormValue("player_id")

	// Game against AI
	if opponent == "" {
		return errorResponse(c, "In development") // Needs building!
	}

	// Check friends
	var count int
	countQuery := dbSelectRow("SELECT count(player1_id) as count FROM friends WHERE (player1_id = ? AND player2_id = ?) OR (player1_id = ? AND player2_id = ?) LIMIT 1", player.Id, opponent, opponent, player.Id)
	countQuery.Scan(&count)
	if count <= 0 {
		return errorResponse(c, "Opponent not found")
	}

	// Check game not already in play
	countQuery = dbSelectRow("SELECT count(player1_id) as count FROM games WHERE ((player1_id = ? AND player2_id = ?) OR (player1_id = ? AND player2_id = ?)) AND status NOT LIKE '%_wins' LIMIT 1", player.Id, opponent, opponent, player.Id)
	countQuery.Scan(&count)
	if count >= 1 {
		return errorResponse(c, "Game already in play")
	}

	// Create new game
	dbInsert("INSERT INTO games(player1_id, player2_id) VALUES (?, ?)", player.Id, opponent)

	// Send request
	pusherClient := pusher.Client{
		AppID:   "1201036",
		Key:     "f372001f569f6df60783",
		Secret:  "b6f7fef8155c288e1e8d",
		Cluster: "eu",
		Secure:  true,
	}

	data := map[string]string{"message": "Game requested"}
	pusherClient.Trigger("games."+opponent, "request", data)

	return successResponse(c, nil, "Game request sent")
}

func gameAccept(c echo.Context) error {
	return successResponse(c, nil, "Game accepted")
}

func gameMove(c echo.Context) error {
	return successResponse(c, nil, "Move added")
}

func gameStatus(e echo.Context) error {
	return successResponse(e, nil, "Game data")
}
