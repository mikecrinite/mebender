package api

import (
	"encoding/json"
	"fmt"
	"log"
	"mebender/model"
	"mebender/service"
	"mebender/util"
	"net/http"
	"time"
)

const GIF_FRAME_RATE = "10"

func CutVideo(w http.ResponseWriter, r *http.Request) {
	methodStart := time.Now()

	// Cut the video into a smaller clip
	if r.Method == "POST" {
		response := model.Response{}
		output, stdout, stderr, err := handleCutVideo(w, r, model.CutVideo)
		response.Duration = util.FormatDuration(time.Since(methodStart))

		// write response
		if err != nil {
			response.Success = false
			response.Error = err
			log.Println(stderr)
			log.Println(err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
		} else {
			// Handle Success
			response.Success = true
			response.ClipLocation = output
			log.Println(stdout)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func GifFromVideo(w http.ResponseWriter, r *http.Request) {
	methodStart := time.Now()

	// Cut the video into a smaller clip
	if r.Method == "POST" {
		response := model.Response{}
		output, _, stderr, err := handleGifFromVideo(w, r, model.GetVideo)
		response.Duration = util.FormatDuration(time.Since(methodStart))

		// write response
		if err != nil {
			response.Success = false
			response.Error = err
			log.Println(err)
			log.Println(stderr)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
		} else {
			// Handle Success
			response.Success = true
			response.ClipLocation = output
			//log.Println(stdout)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func SoundFromVideo(w http.ResponseWriter, r *http.Request) {
		/*
		requestType := model.GetAudio
			// Translate request body into internal struct
			output, stdout, stderr, err := handleCutVideo(w, r, requestType)
	*/
}

func handleCutVideo(w http.ResponseWriter, r *http.Request, requestType string) (string, string, string, error) {
		// Translate request body into internal struct
	request := model.Request{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		// log.Printf("There was an error decoding the request body into the struct: " + r.Body)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not decode the request"))
		return "", "", "", err
	}

	// Validate request
	err = model.ValidateRequest(request, requestType)
	if err != nil {
		handleFailure(err, w, request)
	}

	// debug print request
	fmt.Printf("%#v\n", request)

	// Cut video down to clip
	return service.CutVideo(request)
}

func handleGifFromVideo(w http.ResponseWriter, r *http.Request, requestType string) (string, string, string, error) {
	request := model.Request{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		// log.Printf("There was an error decoding the request body into the struct: " + r.Body)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not decode the request"))
		return "", "", "", err
	}
	// Validate request
	err = model.ValidateRequest(request, requestType)
	if err != nil {
		handleFailure(err, w, request)
	}

	// debug print request
	fmt.Printf("%#v\n", request)

	// Create image frames from video using ffmpeg
	output, err := service.VideoToGifFrames(request, GIF_FRAME_RATE)
	if err != nil {
		log.Printf("An error occurred while writing gif frames: %s", err)
	}

	// Merge frames using imagemagick 
	return service.FramesToGif(output, GIF_FRAME_RATE, util.OUTPUT_LOCATION, request)
}

func handleFailure(err error, w http.ResponseWriter, request model.Request) {
	fmt.Printf("Request was not valid: %#v\n", request)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}
