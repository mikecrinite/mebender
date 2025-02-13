package util

import (
	"errors"
	"fmt"
	"log"
	"mebender/model"
	"os"
	"strings"
	"time"
)

var INPUT_LOCATION = ""
var HOST_OUTPUT_LOCATION = ""
var OUTPUT_LOCATION = "/root/resources/output/"

func init() {
	MOUNT_TARGET_DIR := os.Getenv("MOUNT_TARGET_DIR")
	if MOUNT_TARGET_DIR == "" {
		log.Fatal("No value found for MOUNT_TARGET_DIR")
	} 
	INPUT_LOCATION = MOUNT_TARGET_DIR

	MOUNT_HOST_DIR := os.Getenv("MOUNT_HOST_DIR")
	if MOUNT_HOST_DIR == "" {
		log.Println("No value found for MOUNT_HOST_DIR. Falling back to default, which won't reflect where the file actually is on the host system")
	} else {
		HOST_OUTPUT_LOCATION = MOUNT_HOST_DIR
	}
}

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

// If running in docker, the output filename will represent the filename within the docker container mount point
// Consequently, this is kind of useless to the end user. In that case, we also return the host file location
// for the purpose of displaying to the user
func GetOutputLocation(videoLocation string, isGif bool, requestType string, outputFilename string) (string, string) {
	preparedOutputFilename := getOutputFilename(videoLocation, isGif, requestType, outputFilename)
	
	return fmt.Sprintf("%s%s", OUTPUT_LOCATION, preparedOutputFilename), fmt.Sprintf("%s%s", HOST_OUTPUT_LOCATION, preparedOutputFilename)
}

func getOutputFilename(videoLocation string, isGif bool, requestType string, outputFilename string) string {
	parts := strings.Split(videoLocation, ".")
	now := time.Now().UnixNano()

	// shorten filename to 20 characters... filenames can get really long
	if outputFilename == "" {
		filename := strings.ReplaceAll(parts[0], " ", "")
	
		if len(filename) >= 20 {
			filename = filename[0:20]
		}
	
		parts[0] = filename
	} else {
		parts[0] = outputFilename
	}

	if isGif {
		//return fmt.Sprintf("%s_clip_%d.%s", parts[0], time.Now().UnixNano(), "gif")
		return fmt.Sprintf("%s_frames_%d", parts[0], now)
	} else {
		var file_identifier string
		switch requestType {
		case model.CutVideo:
			file_identifier = "clip"
		case model.PixelateVideo:
			file_identifier = "pixelated"
		}

		return fmt.Sprintf("%d_%s_%s.%s", now, parts[0], file_identifier, parts[1])
	}
}

func FormatDuration(duration time.Duration) string {
	return fmt.Sprintf("%.2f s", duration.Seconds())
}
