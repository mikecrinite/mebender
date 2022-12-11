package service

import (
	"bytes"
	"log"
	"mebender/util"
	"os/exec"
	"time"
)

const PRINT_STDOUT = false

func RunCommand(cmd *exec.Cmd, commandName string, taskName string) error {
	methodStart := time.Now()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.Printf("Running %s command: %s\n", commandName, cmd.String())
	err:= cmd.Run()
	if err != nil {
		log.Println(stderr.String())
		return err
	} else {
		log.Printf("%s %s task took %s\n", commandName, taskName, util.FormatDuration(time.Since(methodStart)))
		if(PRINT_STDOUT){
			log.Println(stdout.String())
		}
	}
	return err
}