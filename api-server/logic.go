package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func main() {
	initDB()
	initHTTP()
}

func initHTTP() {
	log.Println("Starting Nerthus WRS REST API-Server")
	NWRS := mux.NewRouter().StrictSlash(true)

	NWRS.HandleFunc("/", rootEndpoint).Methods("GET")
	NWRS.HandleFunc("/nwrs/user", manipulateUser("CREATE", db)).Methods("POST")
	NWRS.HandleFunc("/nwrs/user", manipulateUser("DELETE", db)).Methods("DELETE")
	NWRS.HandleFunc("/nwrs/container", manipulateContainer("CREATE")).Methods("POST")
	NWRS.HandleFunc("/nwrs/container", manipulateContainer("DELETE")).Methods("DELETE")
	NWRS.HandleFunc("/nwrs/management/port", manipulatePort("GETPORT")).Methods("GET")
	NWRS.HandleFunc("/nwrs/management/port", manipulatePort("RESETPORT")).Methods("PATCH")

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
				if !checkDupl(uQuery[0]) {
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
		if ok1 || len(uQuery) > 0 && ok2 || len(pQuery) > 0 && checkAuth(strings.ToLower(uQuery[0]), pQuery[0]) {
			switch command {
			case "CREATE":
				executeBash("/usr/local/nwrs/scripts/createCont.sh -u "+strings.ToLower(uQuery[0])+" -port "+strconv.Itoa(getPort("NEXT")), true)
				manageContainer("CREATE", uQuery[0])
				json.NewEncoder(w).Encode("CREATING Container.")
			case "DELETE":
				executeBash("/usr/local/nwrs/scripts/removeCont.sh -u "+strings.ToLower(uQuery[0]), true)
				json.NewEncoder(w).Encode("DELETE CONTAINER")
			}
		} else {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode("ERROR: Missing one (or more) required query argument.")
		}
	}
}

func executeBash(path string, special bool) string {
	var out []byte
	if special {
		out, _ = exec.Command("/bin/bash", "-c", path).Output()
	} else {
		out, _ = exec.Command("/bin/bash", path).Output()
	}
	return (string(out))
}

func manipulatePort(command string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch command {
		case "GETPORT":
			json.NewEncoder(w).Encode(getPort("CURRENT"))
		case "RESETPORT":
			setPort("RESET", 0)
			json.NewEncoder(w).Encode(getPort("CURRENT"))
		}
	}
}
