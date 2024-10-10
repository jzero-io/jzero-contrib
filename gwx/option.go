package gwx

func WithDir(dir string) Opts {
	return func(config *gwxConfig) {
		config.Dir = dir
	}
}
