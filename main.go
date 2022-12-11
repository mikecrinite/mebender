package main

import (
	"fmt"
	"log"
	"net/http"

	"mebender/api"
)

const PORT = 8080

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lshortfile)
	log.Printf("Starting mebender on port %d", PORT)

	http.HandleFunc("/cut", api.CutVideo)
	http.HandleFunc("/gif", api.GifFromVideo)
	http.HandleFunc("/sound", api.SoundFromVideo)
	http.HandleFunc("/info", api.VideoInfo)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil))
}
