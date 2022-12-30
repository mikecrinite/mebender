package util

import (
	"errors"
	"fmt"
	"mebender/model"
	"os"
	"strings"
	"time"
)

const INPUT_LOCATION = "/root/resources/input/"
const OUTPUT_LOCATION = "/root/resources/output/"

func MkdirIfNotExists(dirName string) error {
	err := os.Mkdir(dirName, 0755)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		// check that the existing path is a directory
		info, err := os.Stat(dirName)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return errors.New("path exists but is not a directory")
		}
		return nil
	}
	return err
}

func GetOutputLocation(videoLocation string, isGif bool, requestType string) string {
	parts := strings.Split(videoLocation, ".")
	now := time.Now().UnixNano()

	if isGif {
		//return fmt.Sprintf("%s_clip_%d.%s", parts[0], time.Now().UnixNano(), "gif")
		return fmt.Sprintf("%s%s_frames_%d", OUTPUT_LOCATION, parts[0], now)
	} else {
		var file_identifier string
		switch requestType {
		case model.CutVideo:
			file_identifier = "clip"
		case model.PixelateVideo:
			file_identifier = "pixelated"
		}

		return fmt.Sprintf("%s%s_%s_%d.%s", OUTPUT_LOCATION, parts[0], file_identifier, now, parts[1])
	}
}

func FormatDuration(duration time.Duration) string {
	return fmt.Sprintf("%.2f s", duration.Seconds())
}
