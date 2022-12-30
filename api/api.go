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

	"gopkg.in/vansante/go-ffprobe.v2"
)

const GIF_FRAME_RATE = "10"

func CutVideo(w http.ResponseWriter, r *http.Request) {
	methodStart := time.Now()

	// Cut the video into a smaller clip
	if r.Method == "POST" {
		output, err := handleCutVideo(w, r, model.CutVideo)

		writeResponse(w, r, output, err, methodStart, nil)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func GifFromVideo(w http.ResponseWriter, r *http.Request) {
	methodStart := time.Now()

	// Cut the video into a smaller clip
	if r.Method == "POST" {
		output, err := handleGifFromVideo(w, r, model.GetFrames)

		writeResponse(w, r, output, err, methodStart, nil)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func SoundFromVideo(w http.ResponseWriter, r *http.Request) {
	methodStart := time.Now()

	// Cut the video into a smaller clip
	if r.Method == "POST" {
		output, err := handleExtractAudio(w, r, model.GetAudio)

		writeResponse(w, r, output, err, methodStart, nil)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func VideoInfo(w http.ResponseWriter, r *http.Request) {
	methodStart := time.Now()

	// Get ffProbe info
	if r.Method == "GET" {
		request := model.Request{}
		_ = json.NewDecoder(r.Body).Decode(&request)
		probeData, err := service.ProbeVideo(fmt.Sprintf("%s%s", util.INPUT_LOCATION, request.VideoLocation))

		writeResponse(w, r, "", err, methodStart, probeData)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func PixelateVideo(w http.ResponseWriter, r *http.Request){
	methodStart := time.Now()

	// Pixelate video
	if r.Method == "POST" {
		request := model.Request{}
		_ = json.NewDecoder(r.Body).Decode(&request)
		output, err := service.PixelateVideo(request.VideoLocation)

		writeResponse(w, r, output, err, methodStart, nil)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleCutVideo(w http.ResponseWriter, r *http.Request, requestType string) (string, error) {
	// Translate request body into internal struct
	request := model.Request{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		// log.Printf("There was an error decoding the request body into the struct: " + r.Body)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not decode the request"))
		return "", err
	}

	// Validate request
	err = model.ValidateRequest(request, requestType)
	if err != nil {
		handleFailure(err, w, request)
	}

	// debug print request
	log.Printf("%#v\n", request)

	// Cut video down to clip
	return service.CutVideo(request)
}

func handleGifFromVideo(w http.ResponseWriter, r *http.Request, requestType string) (string, error) {
	request := model.Request{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		// log.Printf("There was an error decoding the request body into the struct: " + r.Body)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not decode the request"))
		return "", err
	}
	// Validate request
	err = model.ValidateRequest(request, requestType)
	if err != nil {
		handleFailure(err, w, request)
	}

	// debug print request
	log.Printf("%#v\n", request)

	// Create image frames from video using ffmpeg
	output, err := service.VideoToGifFrames(request, GIF_FRAME_RATE)
	if err != nil {
		log.Printf("An error occurred while writing gif frames: %s", err)
	}

	// Merge frames using imagemagick
	return service.FramesToGif(output, GIF_FRAME_RATE, util.OUTPUT_LOCATION, request)
}

func handleExtractAudio(w http.ResponseWriter, r *http.Request, requestType string) (string, error) {
	request := model.Request{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		// log.Printf("There was an error decoding the request body into the struct: " + r.Body)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not decode the request"))
		return "", err
	}
	// Validate request
	err = model.ValidateRequest(request, requestType)
	if err != nil {
		handleFailure(err, w, request)
	}

	// debug print request
	log.Printf("%#v\n", request)

	// Extract audio with ffmpeg
	return service.ExtractAudio(request)
}

func writeResponse(w http.ResponseWriter, r *http.Request, output string, err error, methodStart time.Time, probeData *ffprobe.ProbeData) {
	response := model.Response{}
	response.Duration = util.FormatDuration(time.Since(methodStart))

	// write response
	if err != nil {
		response.Success = false
		response.Error = err
		log.Println(err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		// Handle Success
		response.Success = true
		response.Location = output
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}

	if probeData != nil {
		probeResponse := model.ProbeResponse{
			Response: response,
			Data:     probeData,
		}
		json.NewEncoder(w).Encode(probeResponse)
	} else {
		json.NewEncoder(w).Encode(response)
	}
}

func handleFailure(err error, w http.ResponseWriter, request model.Request) {
	fmt.Printf("Request was not valid: %#v\n", request)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}
