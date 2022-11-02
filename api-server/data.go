package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

func initDB() {
	nerthusDB, err := sql.Open("sqlite3", "./nwrs.db") // Open the created SQLite File
	if err != nil {
		log.Println("2", err.Error())
	}
	nerthusDB.Exec("CREATE TABLE IF NOT EXISTS user(id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, username TEXT, passwd TEXT);")
	log.Println(getMaxID(nerthusDB))
}

func getMaxID(db *sql.DB) int {
	var id string
	rows, err := db.Query("SELECT id FROM user;")
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		rows.Scan(&id)
		fmt.Println(id)
	}
	return 0
}
