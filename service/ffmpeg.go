package service

import (
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

func CutVideo(cutVideoRequest model.Request) (string, error) {
	start, end, err := getTimes(cutVideoRequest)
	if err != nil {
		//todo
		log.Println(err)
	}
	output := util.GetOutputLocation(cutVideoRequest, false)
	cmd := exec.Command("ffmpeg", "-ss", formatTime(start), "-to", formatTime(end), "-i", fmt.Sprintf("%s/%s", util.INPUT_LOCATION, cutVideoRequest.VideoLocation) /*"-c copy",*/, output)
	err = RunCommand(cmd, "ffmpeg", "CutVideo")

	return output, err
}

func VideoToGifFrames(gifRequest model.Request, frameRate string) (string, error) {
	output := util.GetOutputLocation(gifRequest, true)
	// TODO: cut video if necessary

	// Have to create the directory beforehand or else ffmpeg will fail
	err := util.MkdirIfNotExists(output)
	if err != nil {
		return "", err
	}
	cmd := exec.Command("ffmpeg", "-i", fmt.Sprintf("%s%s", util.INPUT_LOCATION, gifRequest.VideoLocation), "-r", frameRate, "-vcodec", "png", fmt.Sprintf("%s/%s", output, "frame-%03d.png"))
	err = RunCommand(cmd, "ffmpeg", "VideoToGifFrames")

	return output, err
}

func ExtractAudio(request model.Request) (string, error) {
	fullVideoLocation := fmt.Sprintf("%s%s", util.INPUT_LOCATION, request.VideoLocation)

	// ffmpeg -i <input> -map 0:a:0 output0.wav -map 0:a:1 output1.wav -map 0:a:2 output2.wav -map 0:a:3 output3.wav
	// ffmpeg -i <input> -vn -acodec copy output-audio.aac
	// ffmpeg -i <input>.mov <input>.mp3
	probeData, err := ProbeVideo(fullVideoLocation)
	if err != nil {
		return "", err
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
		return "", errors.New("no audio stream with english language tag")
	}

	streamIndex := audioStream.Index
	audioOutputLocation := fmt.Sprintf("%s%d.wav", util.OUTPUT_LOCATION, time.Now().UnixNano())

	cmd := exec.Command("ffmpeg", "-i", fullVideoLocation, "-map", fmt.Sprintf("0:%d", streamIndex), audioOutputLocation)
	err = RunCommand(cmd, "ffmpeg", "ExtractAudio")

	return audioOutputLocation, err
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
