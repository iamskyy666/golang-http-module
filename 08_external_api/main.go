package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CatFactResp struct {
	Fact string `json:"fact"`
	Length int `json:"length"`
}

func writeJson(w http.ResponseWriter, status int, data any){
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}

func fetchcatFact()(CatFactResp, error){
	url:="https://catfact.ninja/fact"

	res,err:=http.Get(url)
	if err!=nil{
		return CatFactResp{},err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK{
		return CatFactResp{}, fmt.Errorf("External api failed: %s",res.Status)
	}

	bodyBytes,err:=io.ReadAll(res.Body)

	if err!=nil{
		return CatFactResp{},err
	}

	var data CatFactResp

	if err:=json.Unmarshal(bodyBytes,&data);err!=nil{
		return CatFactResp{},err
	}

	return data,nil
	
}

func externalApiHandler(w http.ResponseWriter, r *http.Request){
	if r.Method!=http.MethodGet{
		writeJson(w,http.StatusMethodNotAllowed,map[string]any{
			"ok":false,
			"error":"⚠️ Only GET method is allowed!",
		})
		return
	}

	data,err:=fetchcatFact()
	if err!=nil{
		writeJson(w,http.StatusBadRequest,map[string]any{
			"ok":false,
			"error":"⚠️ Failed to fetch data!",
		})
		return
	}

	writeJson(w,http.StatusOK, map[string]any{
		"ok":true,
		"timeStamp":time.Now().UTC(),
		"external":map[string]any{
			"src":"catfact.ninja",
			"fact":data.Fact,
			"length":data.Length,
		},
	})
}

func main() {
	http.HandleFunc("/external",externalApiHandler)
	err := http.ListenAndServe(":5000",nil)
	fmt.Println("ERR:",err.Error())
}