package service

import (
	"fmt"
	"log"
	"mebender/model"
	"mebender/util"
	"os"
	"os/exec"
	"time"
)

func FramesToGif(framesDirectory string, frameRate string, outputDirectory string, request model.Request) (string, error) {
	output := fmt.Sprintf("%s%d.gif", util.OUTPUT_LOCATION, time.Now().UnixNano())

	cmd := exec.Command("convert", "-delay", frameRate, "-loop", "0", "-layers", "optimize", fmt.Sprintf("%s/%s", framesDirectory, "*.png"), output)
	err := RunCommand(cmd, "imagemagick", "FramesToGif")

	if err == nil {
		// TODO: Make this a goroutine. The reason it isn't already is because if you just call `go removeImageFramesDirectory`
		// the parent method will end its execution before this method even gets a chance to run, and consequently it will almost
		// never execute
		go removeDir(framesDirectory)
	}

	return output, err
}

func removeDir(framesDirectory string) {
	log.Printf("Removing directory: %s", framesDirectory)
	err := os.RemoveAll(framesDirectory)
	if err != nil {
		log.Printf("Error while trying to remove framesDirectory: %s", err.Error())
	} else {
		log.Printf("Successfully removed directory: %s", framesDirectory)
	}
}