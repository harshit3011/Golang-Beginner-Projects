package main

import (
	"fmt"
	"myprogram/pkg/handlers"
	"net/http"
)

const port = ":8080"



func main() {
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about",handlers.About)

	fmt.Println(fmt.Sprintf("Server staring on port %v",port))
	http.ListenAndServe(port,nil)
}