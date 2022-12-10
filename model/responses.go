package model

type Response struct {
	ClipLocation string `json:"clipLocation"`
	Success      bool   `json:"success"`
	Error        error  `json:"error"`
	Duration     string `json:"duration"`
}
