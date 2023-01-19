package model

// request types
const (
	CutVideo = "CUT_VIDEO" // short video clip
	GetAudio = "GET_AUDIO" // audio track only
	GetFrames = "GET_FRAMES" // video frames only
	ProbeVideo = "PROBE_VIDEO" // ffProbe video info
	PixelateVideo = "PIXELATE_VIDEO" // pixellate video
)

type Request struct {
	StartTime     *string
	EndTime       *string
	VideoLocation string
	OutputFilename *string
}
