package main

import (
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func leaderboardGlobal(c echo.Context) error {
	/*
		SELECT
			players.username,
			(
				SELECT COUNT(*)
				FROM games as gamewins
				WHERE (gamewins.player1_id = games.player1_id AND status = 'player1_win')
					OR (gamewins.player2_id = games.player1_id AND status = 'player2_win')
			) as `wins`
		FROM games
		LEFT JOIN players
			ON games.player1_id = players.id
		GROUP BY games.player1_id
		ORDER BY wins DESC
		LIMIT 3
	*/
	query := "SELECT players.username, ( SELECT COUNT(*) FROM games as gamewins WHERE (gamewins.player1_id = games.player1_id AND status = 'player1_win') OR (gamewins.player2_id = games.player1_id AND status = 'player2_win') ) as `wins` FROM games LEFT JOIN players ON games.player1_id = players.id GROUP BY games.player1_id ORDER BY wins DESC LIMIT 3"
	dbWinners, err := dbSelect(query)

	if err != nil {
		return errorResponse(c, err.Error())
	}

	defer dbWinners.Close()
	winners := map[string]string{}
	for dbWinners.Next() {
		var username string
		var wins string
		dbWinners.Scan(&username, &wins)
		winners[username] = wins
	}

	return successResponse(c, winners, "Global leaderboard")
}

func leaderboardFriend(c echo.Context) error {
	player := getPlayerFromToken(c)

	dbFriends, err := dbSelect("SELECT player1_id, player2_id FROM friends WHERE player1_id = ? OR player2_id = ?", player.Id, player.Id)
	if err != nil {
		return errorResponse(c, err.Error())
	}

	friendIDs := []string{}
	defer dbFriends.Close()
	for dbFriends.Next() {
		var player1ID string
		var player2ID string
		dbFriends.Scan(&player1ID, &player2ID)

		if player1ID == strconv.Itoa(player.Id) {
			friendIDs = append(friendIDs, player2ID)
		} else {
			friendIDs = append(friendIDs, player1ID)
		}
	}

	// Add current player into set too
	friendIDs = append(friendIDs, strconv.Itoa(player.Id))

	friendIDsString := strings.Join(friendIDs, ",")

	/*
		SELECT
			players.username,
			(
				SELECT COUNT(*)
				FROM games as gamewins
				WHERE (gamewins.player1_id = games.player1_id AND status = 'player1_win')
					OR (gamewins.player2_id = games.player1_id AND status = 'player2_win')
			) as `wins`
		FROM games
		LEFT JOIN players
			ON games.player1_id = players.id
		WHERE games.player1_id IN ("+friendIDsString+")
		GROUP BY games.player1_id
		ORDER BY wins DESC
		LIMIT 3
	*/
	query := "SELECT players.username, ( SELECT COUNT(*) FROM games as gamewins WHERE (gamewins.player1_id = games.player1_id AND status = 'player1_win') OR (gamewins.player2_id = games.player1_id AND status = 'player2_win') ) as `wins` FROM games LEFT JOIN players ON games.player1_id = players.id WHERE games.player1_id IN (" + friendIDsString + ") GROUP BY games.player1_id ORDER BY wins DESC LIMIT 3"
	dbWinners, err := dbSelect(query)

	if err != nil {
		return errorResponse(c, err.Error())
	}

	defer dbWinners.Close()
	winners := map[string]string{}
	for dbWinners.Next() {
		var username string
		var wins string
		dbWinners.Scan(&username, &wins)
		winners[username] = wins
	}

	return successResponse(c, winners, "Friends only leaderboard")
}
