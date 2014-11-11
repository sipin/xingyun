package xingyun

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode"

	"code.google.com/p/xsrftoken"
)

type xsrf struct {
	// Header name value for setting and getting xsrf token.
	Header string
	// Form name value for setting and getting xsrf token.
	Form string
	// Cookie name value for setting and getting xsrf token.
	Cookie string
	// Token generated to pass via header, cookie, or hidden form value.
	Token string
	// This value must be unique per user.
	ID string
	// Secret used along with the unique id above to generate the Token.
	Secret string
	// ErrorFunc is the custom function that replies to the request when ValidToken fails.
	ErrorFunc func(w http.ResponseWriter)
}

// Returns the name of the HTTP header for xsrf token.
func (c *xsrf) GetHeaderName() string {
	return c.Header
}

// Returns the name of the form value for xsrf token.
func (c *xsrf) GetFormName() string {
	return c.Form
}

// Returns the name of the cookie for xsrf token.
func (c *xsrf) GetCookieName() string {
	return c.Cookie
}

// Returns the current token. This is typically used
// to populate a hidden form in an HTML template.
func (c *xsrf) GetToken() string {
	return c.Token
}

// Validates the passed token against the existing Secret and ID.
func (c *xsrf) ValidToken(t string) bool {
	return xsrftoken.Valid(t, c.Secret, c.ID, "POST")
}

// Error replies to the request when ValidToken fails.
func (c *xsrf) Error(w http.ResponseWriter) {
	c.ErrorFunc(w)
}

// xsrfOptions maintains options to manage behavior of Generate.
type xsrfOptions struct {
	// The global secret value used to generate Tokens.
	Secret string
	// HTTP header used to set and get token.
	Header string
	// Form value used to set and get token.
	Form string
	// Cookie value used to set and get token.
	Cookie string
	// If true, send token via X-XSRFToken header.
	SetHeader bool
	// If true, send token via _xsrf cookie.
	SetCookie bool
	// Set the Secure flag to true on the cookie.
	Secure bool
	// The function called when Validate fails.
	ErrorFunc func(w http.ResponseWriter)
	// Array of allowed origins. Will be checked during generation from a cross site request.
	// Must be the complete origin. Example: 'https://golang.org'. You will only need to set this
	// if you are supporting CORS.
	AllowedOrigins []string
}

const domainReg = `/^\.?[a-z\d]+(?:(?:[a-z\d]*)|(?:[a-z\d\-]*[a-z\d]))(?:\.[a-z\d]+(?:(?:[a-z\d]*)|(?:[a-z\d\-]*[a-z\d])))*$/`

func formatName(name string) string {
	r := []rune{}
	for _, c := range name {
		if !unicode.IsSpace(c) {
			r = append(r, unicode.ToLower(c))
		}
	}
	return string(r)
}

func getXSRFxsrfOptions(name string, cfg *Config) *xsrfOptions {
	opts := &xsrfOptions{
		Header:         "X-XSRFToken",
		Form:           formatName("_" + name + "_xsrf"),
		Cookie:         formatName("_" + name + "_xsrf"),
		Secret:         cfg.XSRFSecret,
		SetCookie:      true,
		AllowedOrigins: cfg.XSRFAllowedOrigins,
	}
	if opts.ErrorFunc == nil {
		opts.ErrorFunc = func(w http.ResponseWriter) {
			http.Error(w, "Invalid xsrf token.", http.StatusBadRequest)
		}
	}
	return opts
}

func isAllowedOrigin(opts *xsrfOptions, r *http.Request) bool {
	if r.Header.Get("Origin") == "" {
		return true
	}
	originUrl, err := url.Parse(r.Header.Get("Origin"))
	if err != nil {
		return false
	}
	if originUrl.Host == r.Host {
		return true
	}
	isAllowed := false
	for _, origin := range opts.AllowedOrigins {
		if originUrl.String() == origin {
			isAllowed = true
			break
		}
	}
	return isAllowed
}

func removeCookie(opts *xsrfOptions, w http.ResponseWriter, r *http.Request) {
	expire := time.Now().AddDate(0, 0, -1)
	domain := strings.Split(r.Host, ":")[0]
	if ok, err := regexp.Match(domainReg, []byte(domain)); !ok || err != nil {
		domain = ""
	}

	cookie := &http.Cookie{
		Name:       opts.Cookie,
		Value:      "",
		Path:       "/",
		Domain:     domain,
		Expires:    expire,
		RawExpires: expire.Format(time.UnixDate),
		MaxAge:     0,
		Secure:     opts.Secure,
		HttpOnly:   false,
		Raw:        fmt.Sprintf("%s=%s", opts.Cookie, ""),
		Unparsed:   []string{fmt.Sprintf("token=%s", "")},
	}
	http.SetCookie(w, cookie)
}

func getXSRFId(ctx *Context, name string) string {
	cookie_name := formatName("_" + name + "_xsrf_id")
	id, err := ctx.GetStringCookie(cookie_name)
	if err != nil {
		id = GenRandString(16)
		ctx.SetCookie(cookie_name, id)
		ctx.Logger.Debugf("new xsrf id %s", id)
	} else {
		ctx.Logger.Debugf("load xsrf id %s", id)
	}
	return id
}

func (s *Server) GetXSRFGeneratePipeHandler() PipeHandler {
	return PipeHandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		s.Logger.Tracef("enter xsrf generater")
		defer s.Logger.Tracef("exit xsrf generater")

		ctx := GetContext(r)
		opts := getXSRFxsrfOptions(s.name(), s.Config)
		x := &xsrf{
			Secret:    opts.Secret,
			Header:    opts.Header,
			Form:      opts.Form,
			Cookie:    opts.Cookie,
			ErrorFunc: opts.ErrorFunc,
		}
		x.ID = getXSRFId(ctx, s.name())
		ctx.xsrf = x

		if !isAllowedOrigin(opts, r) {
			ctx.Forbidden()
			return
		}

		// If cookie present, map existing token, else generate a new one.
		if val, err := ctx.GetStringCookie(opts.Cookie); err == nil && val != "" {
			x.Token = val
			s.Logger.Debugf("get xsrf token %s", x.Token)
		} else {
			x.Token = xsrftoken.Generate(x.Secret, x.ID, "POST")
			if opts.SetCookie {
				ctx.SetCookie(opts.Cookie, x.Token)
			}
			s.Logger.Debugf("generate xsrf token %s", x.Token)
		}

		if opts.SetHeader {
			w.Header().Add(opts.Header, x.Token)
		}

		next(w, r)
	})
}

func (s *Server) GetXSRFValidatePipeHandler() PipeHandler {
	return PipeHandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		s.Logger.Tracef("enter xsrf validater")
		defer s.Logger.Tracef("exit xsrf validater")

		if r.Method == "GET" || r.Method == "HEAD" {
			next(w, r)
			return
		}

		ctx := GetContext(r)
		x := ctx.xsrf
		if token := r.Header.Get(x.GetHeaderName()); token != "" {
			if !x.ValidToken(token) {
				s.Logger.Debugf("invalid headker token %s", token)
				opts := getXSRFxsrfOptions(s.name(), s.Config)
				removeCookie(opts, w, r)
				x.Error(w)
				return
			}
			next(w, r)
			return
		}

		if token := r.FormValue(x.GetFormName()); token != "" {
			if !x.ValidToken(token) {
				s.Logger.Debugf("invalid cookie token %s", token)
				opts := getXSRFxsrfOptions(s.name(), s.Config)
				removeCookie(opts, w, r)
				x.Error(w)
				return
			}

			next(w, r)
			return
		}

		s.Logger.Debugf("can't get token from header or form")
		opts := getXSRFxsrfOptions(s.name(), s.Config)
		removeCookie(opts, w, r)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	})
}
