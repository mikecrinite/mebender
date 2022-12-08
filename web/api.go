package web

import (
	"net/http"
	"encoding/json"
	"log"
	"fmt"
	"mebender/model"
	"mebender/service"
)

func CutVideo (w http.ResponseWriter, r *http.Request) {
	// Cut the video into a smaller clip

	if r.Method == "POST" { 
		// Translate request body into internal struct
		request := model.CutVideoRequest{}
		err := json.NewDecoder(r.Body).Decode(&request) 
		if err != nil { 
			// log.Printf("There was an error decoding the request body into the struct: " + r.Body) 
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Could not decode the request"))
			return
		} 
		// debug print request
		fmt.Printf("%#v\n", request)

		// Validate request
		err = model.ValidateCutVideoRequest(request)
		if err != nil {
			handleFailure(err, w, request)
		}

		stdout, stderr, err := service.CutVideo(request)
		log.Println(stdout)
		log.Println(stderr)
		log.Println(err)

		// Handle Success
		////w.Header().Set("Content-Type", "application/json") 
		w.WriteHeader(http.StatusOK) 
	}
}

func handleFailure(err error, w http.ResponseWriter, request model.CutVideoRequest) {
	fmt.Printf("Request was not valid: %#v\n", request)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}