package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gorilla/mux"
)

func main() {
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
	WRS.HandleFunc("/reset/portcount", resetPort).Methods("PATCH")

	http.ListenAndServe((":1234"), WRS)
}

func rootEndpoint(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	json.NewEncoder(w).Encode("Root endpoint hit.")
}

func manipulateUser(command string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uQuery, ok1 := r.URL.Query()["user"]
		pQuery, ok2 := r.URL.Query()["pass"]
		if (ok1 || len(uQuery) > 0 && ok2 || len(pQuery) > 0) && command == "CREATE" {
			executeBash("/usr/local/nwrs/scripts/createUser.sh -u " + uQuery[0] + " -p " + pQuery[0])
			json.NewEncoder(w).Encode("CREATE USER")
		} else if ok1 || len(uQuery) > 0 && command == "REMOVE" {
			executeBash("/usr/local/nwrs/scripts/removeUser.sh -u " + uQuery[0])
			json.NewEncoder(w).Encode("REMOVE USER")
		} else {
			w.WriteHeader(400)
		}
	}
}

func manipulateContainer(command string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uQuery, ok := r.URL.Query()["user"]
		if ok || len(uQuery) > 0 {
			switch command {
			case "CREATE":
				executeBash("/usr/local/nwrs/scripts/createContainer.sh -u " + uQuery[0])
				json.NewEncoder(w).Encode("CREATE CONTAINER")
			case "DELETE":
				executeBash("/usr/local/nwrs/scripts/removeContainer.sh -u " + uQuery[0])
				json.NewEncoder(w).Encode("DELETE CONTAINER")
			}
		} else {
			w.WriteHeader(400)
		}
	}
}

func resetPort(w http.ResponseWriter, _ *http.Request) {
	output := executeBash("/usr/local/nwrs/scripts/resetPort.sh")
	json.NewEncoder(w).Encode(("PortReset finished, status: " + output))
	fmt.Println("Done")
}

func executeBash(path string) string {
	out, err := exec.Command("/bin/bash", "-c", path).Output()
	if err != nil {
		fmt.Println(err)
	}
	return (string(out))
}
