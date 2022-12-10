package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"time"

	"mebender/model"
	"mebender/util"

	"gopkg.in/vansante/go-ffprobe.v2"
)

func CutVideo(cutVideoRequest model.Request) (string, string, string, error) {
	methodStart := time.Now()
	start, end, err := getTimes(cutVideoRequest)
	if err != nil {
		//todo
		log.Println(err)
	}
	output := util.GetOutputLocation(cutVideoRequest, false)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("ffmpeg", "-ss", formatTime(start), "-to", formatTime(end), "-i", fmt.Sprintf("%s/%s", util.INPUT_LOCATION, cutVideoRequest.VideoLocation) /*"-c copy",*/, output)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	log.Printf("Running ffmpeg command: %s\n", cmd.String())
	err = cmd.Run()
	if err != nil {
		log.Printf("ffmpeg CutVideo task took %s\n", util.FormatDuration(time.Since(methodStart)))
	}
	return output, stdout.String(), stderr.String(), err
}

func VideoToGifFrames(gifRequest model.Request, frameRate string) (string, error) {
	methodStart := time.Now()
	output := util.GetOutputLocation(gifRequest, true)
	// TODO: cut video if necessary

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	// Have to create the directory beforehand or else ffmpeg will fail
	err := util.MkdirIfNotExists(output)
	if err != nil {
		return "", err
	}
	cmd := exec.Command("ffmpeg", "-i", fmt.Sprintf("%s%s", util.INPUT_LOCATION, gifRequest.VideoLocation), "-r", frameRate, "-vcodec", "png", fmt.Sprintf("%s/%s", output, "frame-%03d.png"))

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	log.Printf("Creating gif frames in output folder: %s\n", output)
	log.Printf("Running ffmpeg command: %s\n", cmd.String())
	err = cmd.Run()
	if err != nil {
		log.Println(stderr.String())
		return "", err
	} else {
		log.Println("Successfully created gif frames in output folder")
		log.Printf("ffmpeg VideoToGifFrames task took %s\n", util.FormatDuration(time.Since(methodStart)))
		// log.Println(stdout.String())
		return output, nil
	}
}

func ExtractAudio(request model.Request) (string, string, string, error) {
	methodStart := time.Now()
	fullVideoLocation := fmt.Sprintf("%s%s", util.INPUT_LOCATION, request.VideoLocation)

	// ffmpeg -i <input> -map 0:a:0 output0.wav -map 0:a:1 output1.wav -map 0:a:2 output2.wav -map 0:a:3 output3.wav
	// ffmpeg -i <input> -vn -acodec copy output-audio.aac
	// ffmpeg -i <input>.mov <input>.mp3
	probeData, err := ProbeVideo(fullVideoLocation)
	if err != nil {
		return "", "", "", err
	}

	var audioStream *ffprobe.Stream
	streams := probeData.Streams
	for i, s := range streams {
		if s.CodecType == "audio" {
			tagList := s.TagList
			if tagList != nil {
				language := tagList["language"]
				if language != nil {
					if language == "eng" {
						log.Printf("Stream %d had an 'eng' language tag. This stream will be accepted\n", i)
						audioStream = s
					}
				} else {
					log.Printf("Stream %d had no 'language' tag. It will be skipped\n", i)
				}
			} else {
				log.Printf("Stream %d had no tags. It will be skipped", i)
			}
		} 
	}

	if audioStream == nil {
		log.Printf("Could not find an audio stream for video %s\n", request.VideoLocation)
		return "", "", "", errors.New("no audio stream with english language tag")
	}

	streamIndex := audioStream.Index
	audioOutputLocation := fmt.Sprintf("%s%d.wav", util.OUTPUT_LOCATION, time.Now().UnixNano())

	// ffmpeg -i <input> -map 0:a:0 output0.wav
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("ffmpeg", "-i", fullVideoLocation, "-map", fmt.Sprintf("0:%d", streamIndex), audioOutputLocation)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	log.Printf("Extracting audio from video: %s\n", fullVideoLocation)
	log.Printf("Running ffmpeg command: %s\n", cmd.String())
	err = cmd.Run()
	if err != nil {
		log.Println(stderr.String())
		return "", stdout.String(), stderr.String(), err
	} else {
		log.Println("Successfully created audio file in output directory")
		log.Printf("ffmpeg ExtractAudio task took %s\n", util.FormatDuration(time.Since(methodStart)))
		// log.Println(stdout.String())
		return audioOutputLocation, stdout.String(), stderr.String(), nil
	}
}

func ProbeVideo(videoLocation string) (*ffprobe.ProbeData, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()
	return ffprobe.ProbeURL(ctx, videoLocation)
}

func getTimes(cutVideoRequest model.Request) (time.Duration, time.Duration, error) {
	// TODO: find a better way to do this?
	emptyDuration := time.Since(time.Now())

	start, err := time.ParseDuration(*cutVideoRequest.StartTime)
	if err != nil {
		return emptyDuration, emptyDuration, err
	}

	end, err := time.ParseDuration(*cutVideoRequest.EndTime)
	if err != nil {
		return emptyDuration, emptyDuration, err
	}

	return start, end, nil
}

func formatTime(duration time.Duration) string {
	return fmt.Sprintf("%f", duration.Seconds())
}
