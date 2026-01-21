package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CatFact struct{
	Fact string `json:"fact"`
	Length int `json:"length"`
}

func main() {
	url := "https://catfact.ninja/fact"
	resp, err := http.Get(url)
	if err!=nil{
		fmt.Println("ERR:",err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode!=http.StatusOK{
		fmt.Println("resp_status:",resp.Status)
		return
	}

	bodyBytes,err:=io.ReadAll(resp.Body)
	if err!=nil{
		fmt.Println("ERR. reading body:",err.Error())
		return
	}

	var data CatFact

	if err:=json.Unmarshal(bodyBytes,&data);err!=nil{
		fmt.Println("ERR:",err.Error())
		return
	}

	fmt.Println("data:",data.Fact, data.Length)
}

// O/P:
// data: Cats take between 20-40 breaths per minute. 43