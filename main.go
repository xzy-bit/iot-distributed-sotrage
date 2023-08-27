package main

import (
	"IOT_Storage/src/Web"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", Web.HelloWorld)
	http.HandleFunc("/login", Web.Login)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
