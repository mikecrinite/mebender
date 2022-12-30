package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"mebender/api"

	"github.com/99designs/basicauth-go"
	"github.com/go-chi/chi"
)

func main() {
	AUTH_USER := os.Getenv("AUTH_USER")
	if AUTH_USER == "" {
		log.Fatal("No value found for AUTH_USER")
	}

	AUTH_PASS := os.Getenv("AUTH_PASS")
	if AUTH_PASS == "" {
		log.Fatal("No value found for AUTH_PASS")
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Println("No value found for PORT in env. Setting to 8080")
		PORT = "8080"
	}

	log.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lshortfile)
	log.Printf("Starting mebender on port %s", PORT)

	// Add basic auth middleware. the map is a collection of user credentials
	router := chi.NewRouter()
	authMap := map[string][]string{
		AUTH_USER: {
			AUTH_PASS,
		},
	}
	router.Use(basicauth.New("mebender", authMap))

	router.HandleFunc("/cut", api.CutVideo)
	router.HandleFunc("/gif", api.GifFromVideo)
	router.HandleFunc("/sound", api.SoundFromVideo)
	router.HandleFunc("/info", api.VideoInfo)
	// router.HandleFunc("/pixelate", api.PixelateVideo)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PORT), router))
}
