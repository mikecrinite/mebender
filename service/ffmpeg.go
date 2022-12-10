package service

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"time"

	"mebender/model"
	"mebender/util"
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
