package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type TestReq struct {
	Name string `json:"name"`
}

func writeJSON(w http.ResponseWriter, status int, data any){
	w.Header().Set("Content-Type","application.json")
	w.WriteHeader(status)
	_=json.NewEncoder(w).Encode(data)
}


func decodeHandler(w http.ResponseWriter, r *http.Request){
	if r.Method !=http.MethodPost{
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"ok":false,
			"error":"⚠️ Only POST method is allowed!",
		})
		return
	}
	defer r.Body.Close()

	var req TestReq

	dec:=json.NewDecoder(r.Body)

	if err:=dec.Decode(&req); err!=nil{
		writeJSON(w,http.StatusBadRequest, map[string]any{
			"ok":false,
			"error":"Invalid JSON format! ❌",
		})
		return
	}

	req.Name = strings.TrimSpace(req.Name)

	if req.Name == ""{
		writeJSON(w,http.StatusBadRequest, map[string]any{
			"ok":false,
			"error":"⚠️ Name msutn't be empty!",
		})
		return
	
	}

	writeJSON(w,http.StatusOK, map[string]any{
		"ok":true,
		"data":req,
		"timeStamp":time.Now().UTC(),
	})
}

func main() {
	http.HandleFunc("/ok",decodeHandler)
	err := http.ListenAndServe(":5000",nil)
	fmt.Println("ERR:",err.Error())
}