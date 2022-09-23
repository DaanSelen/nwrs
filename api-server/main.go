package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	ipaddr      = "localhost:3306"
	auth        = "nwrsuser:Welkom01"
	workingbase = "nwrs"
)

var (
	nwrsdb *sql.DB
	err    error
)

func main() {
	initDBConnection()
	initHttp()
}

func initHttp() {
	fmt.Println("Starting Nerthus WRS API-Server")
	WRS := mux.NewRouter().StrictSlash(true)

	WRS.HandleFunc("/", rootEndpoint).Methods("GET")
	WRS.HandleFunc("/wrs/user", manipulateUser("CREATE")).Methods("POST")
	WRS.HandleFunc("/wrs/user", manipulateUser("DELETE")).Methods("DELETE")
	WRS.HandleFunc("/wrs/container", manipulateContainer("CREATE")).Methods("POST")
	WRS.HandleFunc("/wrs/container", manipulateContainer("DELETE")).Methods("DELETE")
	WRS.HandleFunc("/wrs/portcount", resetPort).Methods("PATCH")

	http.ListenAndServe((":1234"), WRS)
}

func initDBConnection() {
	fmt.Println("Attempting MariaDB Connection")
	nwrsdb, err = sql.Open("mysql", auth+"@tcp("+ipaddr+")/"+workingbase)
	if err != nil {
		panic("[ERROR] CONNECTING TO DATABASE FAILED")
	}
	_, err := nwrsdb.Exec("CREATE TABLE IF NOT EXISTS user (id INTEGER PRIMARY KEY AUTO_INCREMENT, username TEXT, passwd TEXT)")
	if err != nil {
		panic(err)
	}
}

func rootEndpoint(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	json.NewEncoder(w).Encode("Root endpoint hit.")
}

func manipulateUser(command string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uQuery, ok1 := r.URL.Query()["user"]
		pQuery, ok2 := r.URL.Query()["pass"]
		if ok1 || len(uQuery) > 0 && ok2 || len(pQuery) > 0 {
			switch command {
			case "CREATE":
				executeBash("/usr/local/nwrs/scripts/createUser.sh -u "+uQuery[0]+" -p "+pQuery[0], true)
				manipulateData("CREATE", uQuery[0], pQuery[0])
				json.NewEncoder(w).Encode("CREATE USER")
			case "DELETE":
				if checkAuth(uQuery[0], pQuery[0]) {
					executeBash("/usr/local/nwrs/scripts/removeUser.sh -u "+uQuery[0], true)
					manipulateData("REMOVE", uQuery[0], pQuery[0])
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

func manipulateContainer(command string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uQuery, ok1 := r.URL.Query()["user"]
		pQuery, ok2 := r.URL.Query()["pass"]
		if ok1 || len(uQuery) > 0 && ok2 || len(pQuery) > 0 {
			if checkAuth(uQuery[0], pQuery[0]) {
				switch command {
				case "CREATE":
					executeBash("/usr/local/nwrs/scripts/createContainer.sh -u "+uQuery[0], true)
					json.NewEncoder(w).Encode("CREATE CONTAINER")
				case "DELETE":
					executeBash("/usr/local/nwrs/scripts/removeContainer.sh -u "+uQuery[0], true)
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

func manipulateData(command, username, passwd string) {
	if command == "CREATE" {
		_, err := nwrsdb.Exec("INSERT INTO user(username, passwd) VALUES('" + username + "', '" + passwd + "')")
		if err != nil {
			panic(err)
		}
	} else if command == "REMOVE" {
		_, err := nwrsdb.Exec("DELETE FROM user WHERE username = '" + username + "'")
		if err != nil {
			panic(err)
		}
	}
}

func checkAuth(username, tryPasswd string) bool {
	var passwd string
	err = nwrsdb.QueryRow("SELECT passwd FROM user WHERE username = '" + username + "'").Scan(&passwd)
	if err != nil {
		panic(err)
	}
	fmt.Println(tryPasswd == passwd)
	if tryPasswd == passwd {
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
