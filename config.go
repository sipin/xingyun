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
	if config.StaticDir == "" {
		config.StaticDir = "static"
	}
	// TODO generate random secret
	if config.CookieSecret == "" {
		config.CookieSecret = "D893DCDBCB524C6X"
	}
	if config.XSRFSecret == "" {
		config.XSRFSecret = "D893DCDBCB524C6X"
	}
}
