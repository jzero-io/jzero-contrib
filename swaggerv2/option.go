package swaggerv2

func WithSwaggerHost(swaggerHost string) Opts {
	return func(config *swaggerConfig) {
		config.SwaggerHost = swaggerHost
	}
}

func WithSpecURL(specURL string) Opts {
	return func(config *swaggerConfig) {
		config.SpecURL = specURL
	}
}
