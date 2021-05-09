package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Success bool
	Data    map[string]string
	Message string
}

func errorResponse(c echo.Context, message string) error {
	resp := new(Response)
	resp.Message = message
	resp.Success = false
	return c.JSON(http.StatusBadRequest, resp)
}

func successResponse(c echo.Context, data map[string]string, message string) error {
	resp := new(Response)
	resp.Data = data
	resp.Message = message
	resp.Success = true
	return c.JSON(http.StatusOK, resp)
}
