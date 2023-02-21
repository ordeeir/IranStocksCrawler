package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func SendError(w http.ResponseWriter, status int, message string) {
	error := make(map[string]string)

	error["Message"] = message
	error["Status"] = strconv.Itoa(status)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)

}

func CheckError(err error, comment string) {
	if err != nil {
		fmt.Println("Error: ", err)
		if comment != "" {
			fmt.Println("## ", comment)
		}
	}
}
