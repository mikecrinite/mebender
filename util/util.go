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

	// shorten filename to 20 characters... filenames can get really long
	filename := strings.ReplaceAll(parts[0], " ", "")
	
	if len(filename) >= 20 {
		filename = filename[0:20]
	}

	parts[0] = filename

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

		return fmt.Sprintf("%s%d_%s_%s.%s", OUTPUT_LOCATION, now, parts[0], file_identifier, parts[1])
	}
}

func FormatDuration(duration time.Duration) string {
	return fmt.Sprintf("%.2f s", duration.Seconds())
}
