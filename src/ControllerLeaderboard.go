package main

import (
	"github.com/labstack/echo/v4"
)

func leaderboardGlobal(c echo.Context) error {
	return successResponse(c, nil, "Global leaderboard")
}

func leaderboardFriend(c echo.Context) error {
	return successResponse(c, nil, "Friend leaderboard")
}
