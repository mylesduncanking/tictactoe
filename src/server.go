package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORS())

	// Bearer token required
	r := e.Group("")
	r.Use(middleware.JWT([]byte("secret")))

	// Routes
	// - Auth
	e.POST("/auth/register", authRegister)
	e.POST("/auth/login", authLogin)
	r.GET("/auth/check", authCheck)

	// Friends
	f := r.Group("friends")
	f.GET("", friendList)
	f.POST("/:friend", friendAdd)
	f.DELETE("/:friend", friendDelete)
	f.GET("/:friend", friendGames)

	// Gameplay
	g := r.Group("games")
	g.POST("", gameStart)
	g.GET("/:game/accept", gameAccept)
	g.GET("/:game", gameStatus)
	g.POST("/:game", gameMove)

	// Leaderboard
	l := r.Group("leaderboard")
	l.GET("", leaderboardGlobal)
	l.GET("/friends", leaderboardFriend)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
