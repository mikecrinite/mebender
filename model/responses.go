package model

import "time"

type Response struct {
	ClipLocation string        `json:"clipLocation"`
	Success      bool          `json:"success"`
	Error        error         `json:"error"`
	Duration     time.Duration `json:"duration"`
}
