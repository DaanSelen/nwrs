package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
)

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

func checkAuth(username, tryPasswd string) bool {
	var passwd string
	var trialID int
	nerthusDB.QueryRow("SELECT id, passwd FROM user WHERE username = '"+username+"'").Scan(&trialID, &passwd)
	hashedTryPasswd := hashWithSalt(tryPasswd, trialID)
	if hashedTryPasswd == passwd && len(passwd) != 0 {
		return true
	} else {
		return false
	}
}

func checkDuplicate(user string) bool {
	var duplicateAmount int
	nerthusDB.QueryRow("SELECT COUNT(*) FROM user WHERE EXISTS (SELECT username FROM user WHERE username == '" + user + "');").Scan(&duplicateAmount)
	if duplicateAmount == 0 {
		return false
	} else {
		return true
	}
}

func getMaxID() int {
	var maxID int
	err := nerthusDB.QueryRow("SELECT MAX(id) FROM user;").Scan(&maxID)
	switch err {
	case nil:
		return maxID
	default:
		return 0
	}
}

func hashWithSalt(passwd string, id int) string {
	log.Println(id)
	hash := sha256.New()
	hash.Write([]byte((passwd + strconv.Itoa(id))))
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}

func resetPort(w http.ResponseWriter, _ *http.Request) {
	output := executeBash("/usr/local/nwrs/scripts/resetPort.sh", false)
	json.NewEncoder(w).Encode(("PortReset finished, status: " + output))
	fmt.Println("Done")
}
