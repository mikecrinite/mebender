package model

var REQUEST_TYPES = [...]string{CutVideo, GetAudio, GetVideo}

func ValidateRequest(req Request, requestType string) error {
	// Start Time only is valid

	// validate start time is before end time

	// validate video file location exists and is formatted correctly (windows/unix/linux??)

	return nil
}
