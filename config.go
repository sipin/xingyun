package xingyun

type Config struct {
	CookieSecret string

	StaticDir       string
	StaticPrefix    string
	StaticIndexFile string

	StaticHost            string
	StaticHostExcludeType string
	StaticHostExcludeFile string

	XSRFSecret         string
	XSRFAllowedOrigins []string
}

func setDefaultConfig(config *Config) {
	if config.StaticIndexFile == "" {
		config.StaticIndexFile = "index.html"
	}
}
