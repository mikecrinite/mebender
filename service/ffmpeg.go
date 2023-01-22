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
	output := util.GetOutputLocation(cutVideoRequest.VideoLocation, false, model.CutVideo, *cutVideoRequest.OutputFilename)
	cmd := exec.Command("ffmpeg", "-ss", formatTime(start), "-to", formatTime(end), "-i", fmt.Sprintf("%s/%s", util.INPUT_LOCATION, cutVideoRequest.VideoLocation) /*"-c copy",*/, output)
	err = RunCommand(cmd, "ffmpeg", "CutVideo")

	return output, err
}

func VideoToGifFrames(gifRequest model.Request, frameRate string) (string, error) {
	output := util.GetOutputLocation(gifRequest.VideoLocation, true, model.GetFrames, *gifRequest.OutputFilename)

	// Have to create the directory beforehand or else ffmpeg will fail
	err := util.MkdirIfNotExists(output)
	if err != nil {
		return "", err
	}

	var cmd *exec.Cmd
	if hasStartAndEnd(gifRequest) {
		start, end, err := getTimes(gifRequest)
		if err != nil {
			//todo
			log.Println(err)
		}
		cmd = exec.Command("ffmpeg", "-ss", formatTime(start), "-to", formatTime(end), "-i", fmt.Sprintf("%s%s", util.INPUT_LOCATION, gifRequest.VideoLocation), "-r", frameRate, "-vcodec", "png", fmt.Sprintf("%s/%s", output, "frame-%03d.png"))
	} else {
		cmd = exec.Command("ffmpeg", "-i", fmt.Sprintf("%s%s", util.INPUT_LOCATION, gifRequest.VideoLocation), "-r", frameRate, "-vcodec", "png", fmt.Sprintf("%s/%s", output, "frame-%03d.png"))
	}
	err = RunCommand(cmd, "ffmpeg", "VideoToGifFrames")

	return output, err
}

func ExtractAudio(request model.Request) (string, error) {
	fullVideoLocation := fmt.Sprintf("%s%s", util.INPUT_LOCATION, request.VideoLocation)

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
	var audioOutputLocation string
	if request.OutputFilename != nil {
		audioOutputLocation = fmt.Sprintf("%s%d_%s.wav", util.OUTPUT_LOCATION, time.Now().UnixNano(), *request.OutputFilename)
	}else{
		audioOutputLocation = fmt.Sprintf("%s%d_sound.wav", util.OUTPUT_LOCATION, time.Now().UnixNano())
	}

	var cmd *exec.Cmd
	if hasStartAndEnd(request) {
		start, end, err := getTimes(request)
		if err != nil {
			//todo
			log.Println(err)
		}

		cmd = exec.Command("ffmpeg", "-ss", formatTime(start), "-to", formatTime(end), "-i", fullVideoLocation, "-map", fmt.Sprintf("0:%d", streamIndex), audioOutputLocation)
	} else {
		cmd = exec.Command("ffmpeg", "-i", fullVideoLocation, "-map", fmt.Sprintf("0:%d", streamIndex), audioOutputLocation)
	}
	err = RunCommand(cmd, "ffmpeg", "ExtractAudio")

	return audioOutputLocation, err
}

func ProbeVideo(videoLocation string) (*ffprobe.ProbeData, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()
	return ffprobe.ProbeURL(ctx, videoLocation)
}

func PixelateVideo(videoLocation string) (string, error) {
	// TODO: This doesn't currently work because to enable frei0r, we have to update the ffmpeg configuration, and to do that
	// we have to recompile ffmpeg, and right now I don't feel like doing that
	
	// ffmpeg -i input -vf "frei0r=filter_name=pixeliz0r:filter_params=0.02|0.02" output
	pixel_dimensions := "0.02"

	output := util.GetOutputLocation(videoLocation, false, model.PixelateVideo, "") //, cutVideoRequest.OutputFilename)
	cmd := exec.Command("ffmpeg", "-i", fmt.Sprintf("%s%s", util.INPUT_LOCATION, videoLocation), "-vf", fmt.Sprintf("\"frei0r=filter_name=pixeliz0r:filter_params=%s|%s\"", pixel_dimensions, pixel_dimensions), output)
	err := RunCommand(cmd, "ffmpeg", "PixelateVideo")

	return output, err
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

func hasStartAndEnd(request model.Request) bool {
	return request.StartTime != nil && request.EndTime != nil
}
