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

func cleanup(rec *Recorder) {
	rec.Running = false
	rec.FFmpegCmd = nil
	rec.ChromeCmd = nil
}

func setRunInfo(rec *Recorder, ffmpeg, chrome Runnable) {
	rec.ChromeCmd = chrome
	rec.FFmpegCmd = ffmpeg
	rec.Running = true
}
