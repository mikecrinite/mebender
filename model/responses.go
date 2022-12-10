package model

import "gopkg.in/vansante/go-ffprobe.v2"

type Response struct {
	Location string `json:"location"`
	Success  bool   `json:"success"`
	Error    error  `json:"error"`
	Duration string `json:"duration"`
}

type ProbeResponse struct {
	Response
	Data *ffprobe.ProbeData `json:"data"`
}
