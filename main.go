package main

import (
	"log"
	//"fmt"
	"net/http"

	"mebender/web"
)

func main(){
	log.Print("Starting mebender on port 8080")

	http.HandleFunc("/cut", web.CutVideo)

	log.Fatal(http.ListenAndServe(":8080", nil))
}