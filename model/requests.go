package model

// request types
const (
	CutVideo = "CUT_VIDEO" // short video clip
	GetAudio = "GET_AUDIO" // audio track only
	GetVideo = "GET_VIDEO" // video frames only
)

type Request struct {
	StartTime     *string
	EndTime       *string
	VideoLocation string
}
