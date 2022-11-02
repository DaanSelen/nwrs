package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	initApp()
}

func initApp() {
	fmt.Println("Starting Nerthus WRS SQLite Connection")
	nerthusDB, _ := sql.Open("sqlite3", "./nwrs.db")
	getMaxID(nerthusDB)
	nerthusDB.Exec("CREATE TABLE IF NOT EXISTS user(id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, username TEXT, passwd TEXT);")
	nerthusDB.Exec("UPDATE SQLITE_SEQUENCE SET SEQ=0 WHERE NAME='user';")

	fmt.Println("Starting Nerthus WRS REST API-Server")
	NWRS := mux.NewRouter().StrictSlash(true)

	NWRS.HandleFunc("/", rootEndpoint).Methods("GET")
	NWRS.HandleFunc("/nwrs/user", manipulateUser("CREATE", nerthusDB)).Methods("POST")
	NWRS.HandleFunc("/nwrs/user", manipulateUser("DELETE", nerthusDB)).Methods("DELETE")
	NWRS.HandleFunc("/nwrs/container", manipulateContainer("CREATE", nerthusDB)).Methods("POST")
	NWRS.HandleFunc("/nwrs/container", manipulateContainer("DELETE", nerthusDB)).Methods("DELETE")
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
				if !checkDuplicate(uQuery[0], nerthusDB) {
					executeBash("/usr/local/nwrs/scripts/createUser.sh -u "+strings.ToLower(uQuery[0])+" -p "+pQuery[0], true)
					manipulateData("CREATE", strings.ToLower(uQuery[0]), pQuery[0], nerthusDB)
					json.NewEncoder(w).Encode("CREATE USER")
				} else {
					w.WriteHeader(400)
				}
			case "DELETE":
				if checkAuth(strings.ToLower(uQuery[0]), pQuery[0], nerthusDB) {
					executeBash("/usr/local/nwrs/scripts/removeUser.sh -u "+strings.ToLower(uQuery[0]), true)
					manipulateData("REMOVE", strings.ToLower(uQuery[0]), pQuery[0], nerthusDB)
					json.NewEncoder(w).Encode("REMOVE USER")
				} else {
					w.WriteHeader(401)
				}
			}
		} else {
			w.WriteHeader(400)
		}
	}
}

func manipulateContainer(command string, nerthusDB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uQuery, ok1 := r.URL.Query()["user"]
		pQuery, ok2 := r.URL.Query()["pass"]
		if ok1 || len(uQuery) > 0 && ok2 || len(pQuery) > 0 {
			if checkAuth(strings.ToLower(uQuery[0]), pQuery[0], nerthusDB) {
				switch command {
				case "CREATE":
					executeBash("/usr/local/nwrs/scripts/createContainer.sh -u "+strings.ToLower(uQuery[0]), true)
					json.NewEncoder(w).Encode("Creating Container.")
				case "DELETE":
					executeBash("/usr/local/nwrs/scripts/removeContainer.sh -u "+strings.ToLower(uQuery[0]), true)
					json.NewEncoder(w).Encode("DELETE CONTAINER")
				}
			} else {
				w.WriteHeader(401)
			}
		} else {
			w.WriteHeader(400)
		}
	}
}

func resetPort(w http.ResponseWriter, _ *http.Request) {
	output := executeBash("/usr/local/nwrs/scripts/resetPort.sh", false)
	json.NewEncoder(w).Encode(("PortReset finished, status: " + output))
	fmt.Println("Done")
}

func manipulateData(command, username, passwd string, db *sql.DB) {
	if command == "CREATE" {
		_, err := db.Exec("INSERT INTO user(username, passwd) VALUES('" + username + "', '" + passwd + "')")
		if err != nil {
			panic(err)
		}
	} else if command == "REMOVE" {
		_, err := db.Exec("DELETE FROM user WHERE username = '" + username + "'")
		if err != nil {
			panic(err)
		}
	}
}

func checkAuth(username, tryPasswd string, db *sql.DB) bool {
	var passwd string
	err := db.QueryRow("SELECT passwd FROM user WHERE username = '" + username + "'").Scan(&passwd)
	if err != nil {
		panic(err)
	}
	fmt.Println(tryPasswd == passwd)
	if tryPasswd == passwd && passwd != "" {
		return true
	} else {
		return false
	}
}

func executeBash(path string, special bool) string {
	var out []byte
	var err error
	if special {
		out, err = exec.Command("/bin/bash", "-c", path).Output()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		out, err = exec.Command("/bin/bash", path).Output()
		if err != nil {
			fmt.Println(err)
		}
	}
	return (string(out))
}

func checkDuplicate(user string, db *sql.DB) bool {
	var duplicateAmount int
	db.QueryRow("SELECT COUNT(*) FROM user WHERE EXISTS (SELECT username FROM user WHERE username == '" + user + "');").Scan(&duplicateAmount)
	log.Println(duplicateAmount)
	if duplicateAmount == 0 {
		return false
	} else {
		return true
	}
}

func getMaxID(db *sql.DB) int {
	var maxID int
	err := db.QueryRow("SELECT MAX(id) FROM user;").Scan(&maxID)
	switch err {
	case nil:
		fmt.Println(maxID)
	default:
		maxID = 0
	}
	return maxID
}
