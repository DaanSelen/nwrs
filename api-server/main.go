package main

import (
	"fmt"
	"os/exec"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	initHttp()
}

func initHttp() {
	fmt.Println("Starting WRS API-Server")
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
		switch command {
		case "CREATE":
			json.NewEncoder(w).Encode("CREATE USER")
		case "DELETE":
			json.NewEncoder(w).Encode("DELETE USER")
		}
	}
}

func manipulateContainer(command string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch command {
		case "CREATE":
			json.NewEncoder(w).Encode("CREATE CONTAINER")
		case "DELETE":
			json.NewEncoder(w).Encode("DELETE CONTAINER")
		}
	}
}

func resetPort(w http.ResponseWriter, _ *http.Request) {
	output := executeBash("/home/celdserv/apps/NWRS/scripts/resetPort.sh")
	json.NewEncoder(w).Encode(("PortReset finished, status: " + output))
}

func executeBash(path string) string {
	out, err := exec.Command("/bin/bash", path).Output()
	if err != nil {
		fmt.Println(err)
	}
	return(string(out))
}