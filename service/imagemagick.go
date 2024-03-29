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
	var outputFilename string
	if request.OutputFilename != nil {
		outputFilename = *request.OutputFilename
	} else {
		outputFilename = "animation"
	}
	
	output := fmt.Sprintf("%s%d_%s.gif", util.OUTPUT_LOCATION, time.Now().UnixNano(), outputFilename)

	cmd := exec.Command("convert", "-delay", frameRate, "-loop", "0", "-layers", "optimize", fmt.Sprintf("%s/%s", framesDirectory, "*.png"), output)
	err := RunCommand(cmd, "imagemagick", "FramesToGif")

	if err == nil {
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
