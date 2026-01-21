package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func successHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json") // header
	w.WriteHeader(http.StatusOK) // status

	res:=map[string]any{
		"ok":true,
		"message":"JSON encode successfull! ✅",
		"datetime":time.Now().UTC(),
	}

	_= json.NewEncoder(w).Encode(res)
}

func main() {

	http.HandleFunc("/ok",successHandler)

	err := http.ListenAndServe(":5001",nil)
	fmt.Println("ERR:",err.Error())
}
