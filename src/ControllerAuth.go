package main

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func authRegister(c echo.Context) error {
	username := c.FormValue("username")
	password := generatePasswordHash(c.FormValue("password"))

	id, err := dbInsert("INSERT INTO players(username, password) VALUES (?, ?)", username, password)
	if err != nil {
		return errorResponse(c, err.Error())
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	data := map[string]string{}
	data["token"] = t

	return successResponse(c, data, "Account created")
}

func authLogin(c echo.Context) error {
	username := c.FormValue("username")
	player := dbSelectRow("SELECT id, password FROM players WHERE username = ? LIMIT 1", username)

	var id int
	var password string

	err := player.Scan(&id, &password)
	if err != nil || !checkHash(password, c.FormValue("password")) {
		return errorResponse(c, "Invalid username and/or password")
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	data := map[string]string{}
	data["token"] = t

	return successResponse(c, data, "Welcome back")
}

func authCheck(c echo.Context) error {
	data := map[string]string{}

	player := getPlayerFromToken(c)

	if player.Id <= 0 {
		return errorResponse(c, "Invalid token")
	}

	data["id"] = strconv.Itoa(player.Id)
	return successResponse(c, data, "Logged in")
}
