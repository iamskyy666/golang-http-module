package main

import (
	"fmt"
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

	fmt.Println("status_code:",resp.StatusCode)
	fmt.Println("status:",resp.Status)
}