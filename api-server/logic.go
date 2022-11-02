package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var nerthusDB *sql.DB

func main() {
	initApp()
}

func initApp() {
	fmt.Println("Starting Nerthus WRS SQLite Connection")
	nerthusDB, _ = sql.Open("sqlite3", "./nwrs.db")
	nerthusDB.Exec("CREATE TABLE IF NOT EXISTS user(id INTEGER NOT NULL PRIMARY KEY, username TEXT, passwd TEXT);")

	fmt.Println("Starting Nerthus WRS REST API-Server")
	NWRS := mux.NewRouter().StrictSlash(true)

	NWRS.HandleFunc("/", rootEndpoint).Methods("GET")
	NWRS.HandleFunc("/nwrs/user", manipulateUser("CREATE", nerthusDB)).Methods("POST")
	NWRS.HandleFunc("/nwrs/user", manipulateUser("DELETE", nerthusDB)).Methods("DELETE")
	NWRS.HandleFunc("/nwrs/container", manipulateContainer("CREATE")).Methods("POST")
	NWRS.HandleFunc("/nwrs/container", manipulateContainer("DELETE")).Methods("DELETE")
	NWRS.HandleFunc("/nwrs/portcount", resetPort).Methods("PATCH")

	http.ListenAndServe((":1234"), NWRS)
}

func rootEndpoint(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	json.NewEncoder(w).Encode("Root endpoint hit.")
}

func manipulateUser(command string, nerthusDB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uQuery, ok1 := r.URL.Query()["user"]
		pQuery, ok2 := r.URL.Query()["pass"]
		if ok1 || len(uQuery) > 0 && ok2 || len(pQuery) > 0 {
			switch command {
			case "CREATE":
				if !checkDuplicate(uQuery[0]) {
					nextID := (getMaxID() + 1)
					hashedPasswd := hashWithSalt(pQuery[0], nextID)
					manipulateData("CREATE", strings.ToLower(uQuery[0]), hashedPasswd)
					executeBash("/usr/local/nwrs/scripts/createUser.sh -u "+strings.ToLower(uQuery[0])+" -p "+pQuery[0], true)
					json.NewEncoder(w).Encode("CREATING USER: " + uQuery[0] + " FINISHED")
				} else {
					w.WriteHeader(400)
					json.NewEncoder(w).Encode("ERROR: Duplicate Detected (User already exists).")
				}
			case "DELETE":
				if checkAuth(strings.ToLower(uQuery[0]), pQuery[0]) {
					manipulateData("REMOVE", strings.ToLower(uQuery[0]), pQuery[0])
					executeBash("/usr/local/nwrs/scripts/removeUser.sh -u "+strings.ToLower(uQuery[0]), true)
					json.NewEncoder(w).Encode("REMOVING USER: " + uQuery[0] + " FINISHED")
				} else {
					w.WriteHeader(401)
					json.NewEncoder(w).Encode("ERROR: Incorrect credentials or user does not exist.")
				}
			}
		} else {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode("ERROR: Missing one (or more) required query argument.")
		}
	}
}

func manipulateContainer(command string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uQuery, ok1 := r.URL.Query()["user"]
		pQuery, ok2 := r.URL.Query()["pass"]
		if ok1 || len(uQuery) > 0 && ok2 || len(pQuery) > 0 {
			if checkAuth(strings.ToLower(uQuery[0]), pQuery[0]) {
				switch command {
				case "CREATE":
					executeBash("/usr/local/nwrs/scripts/createContainer.sh -u "+strings.ToLower(uQuery[0]), true)
					json.NewEncoder(w).Encode("CREATING Container.")
				case "DELETE":
					executeBash("/usr/local/nwrs/scripts/removeContainer.sh -u "+strings.ToLower(uQuery[0]), true)
					json.NewEncoder(w).Encode("DELETE CONTAINER")
				}
			} else {
				w.WriteHeader(401)
				json.NewEncoder(w).Encode("ERROR: Duplicate Detected (User already exists).")
			}
		} else {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode("ERROR: Missing one (or more) required query argument.")
		}
	}
}

func manipulateData(command, username, passwd string) {
	if command == "CREATE" {
		_, err := nerthusDB.Exec("INSERT INTO user(id, username, passwd) VALUES('" + strconv.Itoa(getMaxID()+1) + "', '" + username + "', '" + passwd + "');")
		if err != nil {
			log.Println(err)
		}
	} else if command == "REMOVE" {
		_, err := nerthusDB.Exec("DELETE FROM user WHERE username = '" + username + "'")
		if err != nil {
			panic(err)
		}
	}
}
