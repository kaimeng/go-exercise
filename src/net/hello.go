package main

import (
	"net/http"
	"io"
	"log"
)

func helloHandler(w http.ResponseWriter, r *http.Request)  {
	io.WriteString(w, "Hello, Meng!")
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err.Error())
	}
}
