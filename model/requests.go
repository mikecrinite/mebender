package model

type CutVideoRequest struct {
	StartMinutes *int `json:"startMinutes"`
	StartSeconds *int `json:"startSeconds"`
	EndMinutes *int `json:"endMinutes"`
	EndSeconds *int `json:"endSeconds"`
	VideoLocation *string `json:"videoLocation"`
}