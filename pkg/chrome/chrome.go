package chrome

// Builder builds an Chrome command for execution
type Builder struct {
	execPath string
	options  []string
	URL      string
}

func NewBuilder() Builder {
	return Builder{
		execPath: "/usr/bin/google-chrome-stable",
		options: []string{
			"--enable-logging=stderr",
			"--autoplay-policy=no-user-gesture-required",
			"--no-sandbox",
			"--disable-infobars",
			"--kiosk",
			"--start-maximized --window-position=0,0",
			"--window-size=1280,720",
		},
	}
}

func (b Builder) WithOptions(options ...string) Builder {
	if len(options) > 0 {

		b.options = options
	}
	return b
}

func (b Builder) WithURL(URL string) Builder {
	b.URL = URL
	return b
}

func (b Builder) Build() ([]string, error) {
	cmdArgs := []string{b.execPath}
	cmdArgs = append(cmdArgs, b.options...)
	if b.URL != "" {
		cmdArgs = append(cmdArgs, "--app=\""+b.URL+"\"")
	}
	return cmdArgs, nil
}
