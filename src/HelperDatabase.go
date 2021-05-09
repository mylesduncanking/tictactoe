package main

import (
	"database/sql"
	"fmt"
)

func dbConnect() *sql.DB {
	// DB Connection
	// db, err := sql.Open("mysql", "root:@tcp(host.docker.internal:3306)/tictactoe")
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/tictactoe")
	if err != nil {
		fmt.Println(err.Error())
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
	}

	return db
}

func dbInsert(query string, params ...interface{}) (int, error) {
	db := dbConnect()

	res, err := db.Exec(query, params...)

	if err != nil {
		return 0, err
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(lastId), nil
}

func dbSelect(query string, params ...interface{}) (*sql.Rows, error) {
	db := dbConnect()

	// Run query
	insert, err := db.Query(query, params...)

	// If there is an error inserting, handle it
	if err != nil {
		return nil, err
	}

	return insert, err
}

func dbSelectRow(query string, params ...interface{}) *sql.Row {
	db := dbConnect()
	return db.QueryRow(query, params...)
}

func dbDelete(query string, params ...interface{}) bool {
	db := dbConnect()
	_, err := db.Query(query, params...)
	return err == nil
}
