package main

import (
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func friendList(c echo.Context) error {
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

	friendIDsString := strings.Join(friendIDs, ",")
	dbPlayers, _ := dbSelect("SELECT id, username FROM players WHERE id IN (" + friendIDsString + ") LIMIT " + strconv.Itoa(len(friendIDs)))

	friends := map[string]string{}

	defer dbPlayers.Close()
	var friendID string
	var friendUsername string
	for dbPlayers.Next() {
		dbPlayers.Scan(&friendID, &friendUsername)
		friends[friendID] = friendUsername
	}

	return successResponse(c, friends, friendIDsString)
}

func friendAdd(c echo.Context) error {
	player := getPlayerFromToken(c)
	friend := c.Param("friend")

	// Check is valid ID
	var friendID int
	friendIDQuery := dbSelectRow("SELECT id FROM players WHERE username = ? LIMIT 1", friend)
	friendIDQuery.Scan(&friendID)
	if friendID <= 0 {
		return errorResponse(c, "Invalid username")
	}

	// Check not already friends
	var count int
	countQuery := dbSelectRow("SELECT count(player1_id) as count FROM friends WHERE (player1_id = ? AND player2_id = ?) OR (player1_id = ? AND player2_id = ?) LIMIT 1", player.Id, friendID, friendID, player.Id)
	countQuery.Scan(&count)
	if count > 0 {
		return errorResponse(c, "Already friends")
	}

	// Add friend
	dbInsert("INSERT INTO friends(player1_id, player2_id) VALUES (?, ?)", player.Id, friendID)

	return successResponse(c, nil, "Friend added")
}

func friendGames(c echo.Context) error {
	player := getPlayerFromToken(c)
	friend := c.Param("friend")

	dbGames, err := dbSelect("SELECT id, status FROM games WHERE (player1_id = ? AND player2_id = ?) OR (player1_id = ? AND player2_id = ?)", player.Id, friend, friend, player.Id)

	if err != nil {
		return errorResponse(c, err.Error())
	}

	games := map[string]string{}
	defer dbGames.Close()
	for dbGames.Next() {
		var gameID string
		var gameStatus string
		dbGames.Scan(&gameID, &gameStatus)
		games[gameID] = gameStatus
	}

	return successResponse(c, games, "")
}

func friendDelete(c echo.Context) error {
	player := getPlayerFromToken(c)
	friend := c.Param("friend")

	// Delete friend
	dbDelete("DELETE FROM friends WHERE (player1_id = ? AND player2_id = ?) OR (player1_id = ? AND player2_id = ?)", player.Id, friend, friend, player.Id)

	return successResponse(c, nil, "Friend deleted")
}
