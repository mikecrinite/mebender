package main

import (
	"log"
	//"fmt"
	"net/http"

	"mebender/api"
)

func main() {
	log.Print("Starting mebender on port 8080")

	http.HandleFunc("/cut", api.CutVideo)
	http.HandleFunc("/gif", api.GifFromVideo)
	http.HandleFunc("/sound", api.SoundFromVideo)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
