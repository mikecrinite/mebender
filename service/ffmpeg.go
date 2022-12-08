package service

import (
    "bytes"
    "fmt"
    "log"
    "os/exec"
	"time"
	"strings"

	"mebender/model"
)

//const ffmpegFormatString = "ffmpeg -ss 00:01:00 -to 00:02:00 -i input.mp4 -c copy output.mp4"
//const ShellToUse = "bash"

func CutVideo(cutVideoRequest model.CutVideoRequest) (string, string, error) {
	start, end, err := getTimes(cutVideoRequest)
	if err != nil {
		//todo
		log.Println(err)
	}
	output := getOutputLocation(cutVideoRequest)

    var stdout bytes.Buffer
    var stderr bytes.Buffer
    //cmd := exec.Command(ShellToUse, "-c", command)

	cmd := exec.Command("ffmpeg", "-ss", formatTime(start), "-to", formatTime(end), "-i", cutVideoRequest.VideoLocation, /*"-c copy",*/ output)

    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    err = cmd.Run()
    return stdout.String(), stderr.String(), err
}

func getTimes(cutVideoRequest model.CutVideoRequest) (time.Duration, time.Duration, error){
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

func getOutputLocation(cutVideoRequest model.CutVideoRequest) string {
	parts := strings.Split(cutVideoRequest.VideoLocation, ".")

	return fmt.Sprintf("%s_clip_%d.%s", parts[0], time.Now().UnixNano(), parts[1])
}

func formatTime(duration time.Duration) string {
	//return fmt.Sprintf("%f:%f:%f", duration.Hours(), duration.Minutes(), duration.Seconds())
	return fmt.Sprintf("%f", duration.Seconds())
}