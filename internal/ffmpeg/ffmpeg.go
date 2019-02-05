package ffmpeg

// Builder builds an FFmpeg command with safe default options
type Builder struct {
	execPath string
	options  [][]string
	args     [][]string
}

func NewBuilder() Builder {
	return Builder{
		execPath: "/usr/bin/ffmpeg",
		options: [][]string{
			{"-y", ""},
			{"-v", "info"},
			{"-f", "x11grab"},
			{"-draw_mouse", "0"},
			{"-r", "24"},
			{"-s", "1280x720"},
			{"-thread_queue_size", "4096"},
			{"-i", ":99.0+0,0"},
			{"-f", "pulse"},
			{"-thread_queue_size", "4096"},
			{"-i", "default"},
			{"-acodec", "aac"},
			{"-strict", "-2"},
			{"-ar", "44100"},
			{"-c:v", "libx264"},
			{"-x264opts", "no-scenecut"},
			{"-preset", "veryfast"},
			{"-profile:v", "main"},
			{"-level", "3.1"},
			{"-pix_fmt", "yuv420p"},
			{"-r", "24"},
			{"-crf", "25"},
			{"-g", "48"},
			{"-keyint_min", "48"},
			{"-force_key_frames", "\"expr:gte(t,n_forced*2)\""},
			{"-tune", "zerolatency"},
			{"-b:v", "2800k"},
			{"-maxrate", "2996k"},
			{"-bufsize", "4200k"},
		},
	}
}

func (b Builder) WithOptions(options ...[]string) Builder {
	if len(options) != 0 {
		b.options = options
	}
	return b
}

func (b Builder) WithArguments(args ...[]string) Builder {
	b.args = args
	return b
}

func (b Builder) Build() ([]string, error) {
	cmdArgs := []string{b.execPath}
	for _, opt := range b.options {
		cmdArgs = append(cmdArgs, opt...)
	}
	for _, opt := range b.args {
		cmdArgs = append(cmdArgs, opt...)
	}
	return cmdArgs, nil
}
