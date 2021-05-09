package main

import "database/sql"

type Game struct {
	Code    string
	P1Moves sql.NullString
	P2Moves sql.NullString
	Status  string
}
