package service

import (
	"bytes"
	"fmt"
	"log"
	"mebender/model"
	"mebender/util"
	"os/exec"
	"time"
)

func FramesToGif(framesDirectory string, frameRate string, outputDirectory string, request model.Request) (string, string, string, error) {
	methodStart := time.Now()
	
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	output := fmt.Sprintf("%s%d.gif", util.OUTPUT_LOCATION, methodStart.UnixNano())

	cmd := exec.Command("convert", "-delay", frameRate, "-loop", "0", "-layers", "optimize", fmt.Sprintf("%s/%s", framesDirectory, "*.png"), output)
	
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	log.Printf("Running imagemagick command: %s\n", cmd.String())
	err := cmd.Run()

	if err != nil {
		// TODO: If successful, remove the frames directory in separate goroutine
		log.Printf("imagemagick task took %s\n", util.FormatDuration(time.Since(methodStart)))
		go func(framesDirectory string) {
			log.Printf("Removing directory: %s", framesDirectory)
			log.Printf("Successfully removed directory: %s", framesDirectory)
		}(framesDirectory)
	}

	return output, stdout.String(), stderr.String(), err
}