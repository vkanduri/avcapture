package recorder

import (
	"sync"
)

type Runnable interface {
	Start() error
	Stop() error
}

type Recorder struct {
	Running bool
	mtx     sync.Mutex

	FFmpegCmd Runnable
	ChromeCmd Runnable
}

type FFmpeg struct {
	Params  [][]string `json:"params"`
	Options [][]string `json:"options"`
}

type Chrome struct {
	URL     string   `json:"url"`
	Options []string `json:"options"`
}
