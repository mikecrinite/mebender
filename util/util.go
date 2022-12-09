package util

import (
	"errors"
	"os"
)

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