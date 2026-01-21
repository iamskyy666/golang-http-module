package main

import (
	"fmt"
	"net/http"
)

func HelloHandler(w http.ResponseWriter,r *http.Request){
	if r.Method != http.MethodGet{
		http.Error(w,"Only GET method is allowed!",http.StatusMethodNotAllowed)
		return
	}

	_,_=w.Write([]byte("Hello from GOLANG http-server!"))

	// if err!=nil{
	// 	fmt.Println("ERROR:",err.Error())
	// 	return
	// }
}

func main() {
	http.HandleFunc("/hello", HelloHandler)
	fmt.Println("Try going to 8080 port..")

	err:= http.ListenAndServe(":8080",nil)
	fmt.Println("ERROR:",err.Error())
}

//03:03:00