package service

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"mebender/model"
	"mebender/util"
)

const INPUT_LOCATION = "/root/resources/input/"
const OUTPUT_LOCATION = "/root/resources/output/"

func CutVideo(cutVideoRequest model.Request) (string, string, string, error) {
	start, end, err := getTimes(cutVideoRequest)
	if err != nil {
		//todo
		log.Println(err)
	}
	output := getOutputLocation(cutVideoRequest, false)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("ffmpeg", "-ss", formatTime(start), "-to", formatTime(end), "-i", fmt.Sprintf("%s/%s", INPUT_LOCATION, cutVideoRequest.VideoLocation) /*"-c copy",*/, output)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	return output, stdout.String(), stderr.String(), err
}

func VideoToGif(gifRequest model.Request) (string, string, string, error) {
	const frameRate = "24"
	output := getOutputLocation(gifRequest, true)
	// TODO: cut video if necessary

	// ffmpeg -ss 00:00:00.000 -i yesbuddy.mov -pix_fmt rgb24 -r 10 -s 320x240 -t 00:00:10.000 output.gif
	// ffmpeg -i <input.mp4>  -r 10 -f image2pipe -vcodec ppm - | convert -delay 10 -loop 0 -layers Optimize - <output.gif>

	/*
	ffmpegCommand := exec.Command(fmt.Sprintf("ffmpeg -i %s  -r %d -f image2pipe -vcodec ppm - ", gifRequest.VideoLocation, frameRate))
	imageMagick := exec.Command(fmt.Sprintf("convert -delay %d -loop 0 -layers Optimize - %s", frameRate, output))

	r, w := io.Pipe()
	ffmpegCommand.Stdout = w
	imageMagick.Stdin = r

	imageMagick.Stdin = r

	var b2 bytes.Buffer
	imageMagick.Stdout = &b2

	ffmpegCommand.Start()
	imageMagick.Start()
	imageMagick.Wait()
	written, err := io.Copy(os.Stdout, &b2)
	if err != nil {
		log.Println(err)
	}else{
		log.Printf("Written: %d", written)
	}

	return "", "", "", nil
	*/

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	//stringCommand := fmt.Sprintf("ffmpeg -i %s  -r %d -f image2pipe -vcodec ppm - | convert -delay %d -loop 0 -layers Optimize - %s", gifRequest.VideoLocation, frameRate, frameRate, output)
	//stringCommand := fmt.Sprintf("ffmpeg -i %s  -r %d %s/frame-%03d.png", gifRequest.VideoLocation, frameRate)
	// ffmpeg -i <input.mp4> -r 10 <path-to-output-folder>/frame-%03d.png

	/*
	log.Printf("Running:\nsh -c %s", stringCommand)
	cmd := exec.Command("sh", "-c", stringCommand)
	*/

	err := util.MkdirIfNotExists(output)
	if err != nil {
		log.Println(err)
		return "", "", "", err
	}
	cmd := exec.Command("ffmpeg", "-i", fmt.Sprintf("%s%s", INPUT_LOCATION, gifRequest.VideoLocation), "-r", frameRate, "-vcodec", "png", fmt.Sprintf("%s/%s", output, "frame-%03d.png"))

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	return output, stdout.String(), stderr.String(), err
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

func getOutputLocation(cutVideoRequest model.Request, isGif bool) string {
	parts := strings.Split(cutVideoRequest.VideoLocation, ".")
	now := time.Now().UnixNano()

	if isGif {
		//return fmt.Sprintf("%s_clip_%d.%s", parts[0], time.Now().UnixNano(), "gif")
		return fmt.Sprintf("%s%s_frames_%d", OUTPUT_LOCATION, parts[0], now)
	} else {
		return fmt.Sprintf("%s/%s_clip_%d.%s", OUTPUT_LOCATION, parts[0], now, parts[1])
	}
}

func formatTime(duration time.Duration) string {
	return fmt.Sprintf("%f", duration.Seconds())
}
