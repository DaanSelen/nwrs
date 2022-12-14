package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() {
	log.Println("Starting Nerthus WRS SQLite Connection")
	db, _ = sql.Open("sqlite3", "./nwrs.db")

	db.Exec("CREATE TABLE IF NOT EXISTS user(id INTEGER NOT NULL PRIMARY KEY, username TEXT, passwd TEXT);")
	db.Exec("CREATE TABLE IF NOT EXISTS cont(id INTEGER NOT NULL PRIMARY KEY, owner TEXT, seq INTEGER, port INTEGER);")
	db.Exec("CREATE TABLE IF NOT EXISTS port(id INTEGER NOT NULL PRIMARY KEY, number INTEGER);")

	db.Exec("INSERT INTO port VALUES('0', '10000')")
	log.Println("SQLITE Connection Succesful.")
}

func manipulateData(command, username, passwd string) {
	if command == "CREATE" {
		_, err := db.Exec("INSERT INTO user(id, username, passwd) VALUES(null, '" + username + "', '" + passwd + "');")
		if err != nil {
			log.Println(err)
		}
	} else if command == "REMOVE" {
		_, err := db.Exec("DELETE FROM user WHERE username = '" + username + "'")
		if err != nil {
			log.Println(err)
		}
	}
}

func manageContainer(command, username string, optnumber int) {
	switch command {
	case "CREATE":
		sequence := checkContainer(username)
		db.Exec("INSERT INTO cont VALUES(null, '" + username + "', '" + strconv.Itoa(sequence) + "', '" + strconv.Itoa(optnumber) + "')")
		log.Println("INSERT INTO cont VALUES(null, '" + username + "', '" + strconv.Itoa(sequence) + "', '" + strconv.Itoa(optnumber) + "')")
	case "DELETE":
		db.Exec("DELETE FROM cont WHERE owner =='" + username + "' AND seq == " + strconv.Itoa(optnumber))
	}
}

func checkContainer(username string) int {
	var sequence int
	db.QueryRow("SELECT MAX(seq) FROM cont WHERE owner == '" + username + "'").Scan(&sequence)
	return sequence + 1
}

func listContainers(username string) []ListedContainer {
	var listedContainers []ListedContainer
	data, err := db.Query("SELECT * FROM cont WHERE owner == '" + username + "'")
	if err != nil {
		fmt.Println(err)
	}
	defer data.Close()
	for data.Next() {
		var singleContainer ListedContainer
		data.Scan(&singleContainer.ID, &singleContainer.Owner, &singleContainer.Seq, &singleContainer.Port)
		listedContainers = append(listedContainers, singleContainer)
	}
	return listedContainers
}

func generateContainername(username string) string {
	var sequence int
	var owner string
	db.QueryRow("SELECT owner FROM cont WHERE owner == '" + username + "'").Scan(&owner)
	db.QueryRow("SELECT MAX(seq) FROM cont WHERE owner == '" + username + "'").Scan(&sequence)
	return (owner + "-web-" + strconv.Itoa(sequence))
}

func getPort(command string) int {
	var number int
	switch command {
	case "CURRENT":
		db.QueryRow("SELECT number FROM port;").Scan(&number)
		return number
	case "NEXT":
		db.QueryRow("SELECT number FROM port;").Scan(&number)
		setPort("SETNEW", (number + 1))
		return number + 1
	default:
		return 0
	}
}

func setPort(command string, newPort int) {
	switch command {
	case "SETNEW":
		db.Exec("UPDATE port SET number = " + strconv.Itoa(newPort) + " WHERE id = 0;")
	case "RESET":
		db.Exec("UPDATE port SET number = 10001 WHERE id = 0;")
	}
}

func checkAuth(username, tryPasswd string) bool {
	var passwd string
	var trialID int
	db.QueryRow("SELECT id, passwd FROM user WHERE username = '"+username+"'").Scan(&trialID, &passwd)
	hashedTryPasswd := hashWithSalt(tryPasswd, trialID)
	if hashedTryPasswd == passwd && len(passwd) != 0 {
		return true
	} else {
		return false
	}
}

func checkDupl(user string) bool {
	var duplicateAmount int
	db.QueryRow("SELECT COUNT(*) FROM user WHERE EXISTS (SELECT username FROM user WHERE username == '" + user + "');").Scan(&duplicateAmount)
	if duplicateAmount == 0 {
		return false
	} else {
		return true
	}
}

func getMaxID() int {
	var maxID int
	err := db.QueryRow("SELECT MAX(id) FROM user;").Scan(&maxID)
	switch err {
	case nil:
		return maxID
	default:
		return 0
	}
}

func hashWithSalt(passwd string, id int) string {
	hash := sha256.New()
	hash.Write([]byte((passwd + strconv.Itoa(id))))
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}
