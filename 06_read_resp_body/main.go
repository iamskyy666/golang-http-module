package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	url := "https://jsonplaceholder.typicode.com/todos"
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
		fmt.Println("ERR:",err)
		return
	}

	bodyTxt:=string(bodyBytes)
	max := min(len(bodyTxt), 250)

	fmt.Println(bodyTxt[:max])
}