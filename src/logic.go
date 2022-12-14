package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type UserRequestForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type ContRequestForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Seq      int    `json:"seq"`
}
type ListedContainer struct {
	ID    int    `json:"id"`
	Owner string `json:"owner"`
	Seq   int    `json:"seq"`
	Port  int    `json:"port"`
}

func main() {
	initDB()
	initHTTP()
}

func initHTTP() {
	log.Println("Starting Nerthus WRS REST API-Server")
	NWRS := mux.NewRouter().StrictSlash(true)

	NWRS.HandleFunc("/", rootEndpoint).Methods("GET")
	NWRS.HandleFunc("/nwrs/user", manipulateUser("CREATE")).Methods("POST")
	NWRS.HandleFunc("/nwrs/user", manipulateUser("DELETE")).Methods("DELETE")
	NWRS.HandleFunc("/nwrs/container", manipulateContainer("CHECK")).Methods("GET")
	NWRS.HandleFunc("/nwrs/container", manipulateContainer("CREATE")).Methods("POST")
	NWRS.HandleFunc("/nwrs/container", manipulateContainer("DELETE")).Methods("DELETE")
	NWRS.HandleFunc("/nwrs/management/port", manipulatePort("GETPORT")).Methods("GET")
	NWRS.HandleFunc("/nwrs/management/port", manipulatePort("RESETPORT")).Methods("PATCH")

	http.ListenAndServe((":1234"), NWRS)
}

func rootEndpoint(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	json.NewEncoder(w).Encode("Root endpoint hit. Version 0.2")
}

func manipulateUser(command string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form UserRequestForm
		json.NewDecoder(r.Body).Decode(&form)
		if len(form.Username) > 0 && len(form.Password) > 0 {
			switch command {
			case "CREATE":
				if !checkDupl(form.Username) {
					nextID := (getMaxID() + 1)
					hashedPasswd := hashWithSalt(form.Password, nextID)
					manipulateData("CREATE", strings.ToLower(form.Username), hashedPasswd)
					executeBash("./scripts/createUser.sh -u "+strings.ToLower(form.Username)+" -p "+form.Password, true)
					json.NewEncoder(w).Encode("CREATING USER: " + form.Username + " FINISHED")
				} else {
					w.WriteHeader(400)
					json.NewEncoder(w).Encode("ERROR: Duplicate Detected (User already exists).")
				}
			case "DELETE":
				if checkAuth(strings.ToLower(form.Username), form.Password) {
					manipulateData("REMOVE", strings.ToLower(form.Username), form.Password)
					executeBash("./scripts/removeUser.sh -u "+strings.ToLower(form.Username), true)
					json.NewEncoder(w).Encode("REMOVING USER: " + form.Username + " FINISHED")
				} else {
					w.WriteHeader(401)
					json.NewEncoder(w).Encode("ERROR: Authentication failure, or non-existing user.")
					log.Println("Invalid or incorrect user DELETE request.")
				}
			}
		} else {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode("ERROR: Missing one (or more) required query argument(s).")
		}
	}
}

func manipulateContainer(command string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form ContRequestForm
		json.NewDecoder(r.Body).Decode(&form)
		if len(form.Username) > 0 && len(form.Password) > 0 {
			if checkAuth(strings.ToLower(form.Username), form.Password) {
				switch command {
				case "CREATE":
					port := getPort("NEXT")
					manageContainer("CREATE", form.Username, port)
					containername := generateContainername(form.Username)
					executeBash("./scripts/createCont.sh  -u "+strings.ToLower(form.Username)+" -cn "+containername+" -port "+strconv.Itoa(port), true)
					log.Println("./scripts/createCont.sh  -u " + strings.ToLower(form.Username) + " -cn " + containername + " -port " + strconv.Itoa(port))
					json.NewEncoder(w).Encode("CREATING Container.")
				case "DELETE":
					if form.Seq > 0 {
						manageContainer("DELETE", form.Username, form.Seq)
						executeBash("./scripts/removeCont.sh -u "+strings.ToLower(form.Username), true)
						json.NewEncoder(w).Encode("DELETE CONTAINER")
					}
				case "CHECK":
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(listContainers(form.Username))
				}
			} else {
				w.WriteHeader(401)
				json.NewEncoder(w).Encode("ERROR: Authentication failure.")
			}
		} else {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode("ERROR: Missing one (or more) required query argument(s).")
		}
	}
}

func executeBash(path string, special bool) string {
	var out []byte
	var err error
	if special {
		out, err = exec.Command("/bin/bash", "-c", path).Output()
		log.Println(string(out), err)
	} else {
		out, err = exec.Command("/bin/bash", path).Output()
		log.Println(string(out), err)
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
