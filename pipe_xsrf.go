package xingyun

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"code.google.com/p/xsrftoken"
)

// XSRF is used to get the current token and validate a suspect token.
type XSRF interface {
	// Return HTTP header to search for token.
	GetHeaderName() string
	// Return form value to search for token.
	GetFormName() string
	// Return cookie name to search for token.
	GetCookieName() string
	// Return the token.
	GetToken() string
	// Validate by token.
	ValidToken(t string) bool
	// Error replies to the request with a custom function when ValidToken fails.
	Error(w http.ResponseWriter)
}

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

// Options maintains options to manage behavior of Generate.
type Options struct {
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

func getXSRFOptions(cfg *Config) *Options {
	opts := &Options{
		Header:         "X-XSRFToken",
		Form:           "_xsrf",
		Cookie:         "_xsrf",
		Secret:         cfg.XSRFSecret,
		AllowedOrigins: cfg.XSRFAllowedOrigins,
	}
	if opts.ErrorFunc == nil {
		opts.ErrorFunc = func(w http.ResponseWriter) {
			http.Error(w, "Invalid xsrf token.", http.StatusBadRequest)
		}
	}
	return opts
}

func isAllowedOrigin(opts *Options, r *http.Request) bool {
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

func setCookie(opts *Options, x *xsrf, w http.ResponseWriter, r *http.Request) {
	expire := time.Now().AddDate(0, 0, 1)
	// Verify the domain is valid. If it is not, set as empty.
	domain := strings.Split(r.Host, ":")[0]
	if ok, err := regexp.Match(domainReg, []byte(domain)); !ok || err != nil {
		domain = ""
	}

	cookie := &http.Cookie{
		Name:       opts.Cookie,
		Value:      x.Token,
		Path:       "/",
		Domain:     domain,
		Expires:    expire,
		RawExpires: expire.Format(time.UnixDate),
		MaxAge:     0,
		Secure:     opts.Secure,
		HttpOnly:   false,
		Raw:        fmt.Sprintf("%s=%s", opts.Cookie, x.Token),
		Unparsed:   []string{fmt.Sprintf("token=%s", x.Token)},
	}
	http.SetCookie(w, cookie)
}

func (s *Server) GetXSRFGeneratePipeHandler() PipeHandler {
	return PipeHandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.Handler) {
		s.Logger.Tracef("enter")
		defer s.Logger.Tracef("exit")

		ctx := GetContext(r)
		opts := getXSRFOptions(s.Config)
		x := &xsrf{
			Secret:    opts.Secret,
			Header:    opts.Header,
			Form:      opts.Form,
			Cookie:    opts.Cookie,
			ErrorFunc: opts.ErrorFunc,
		}
		x.ID = ctx.UserID
		ctx.xsrf = x

		if !isAllowedOrigin(opts, r) {
			ctx.Forbidden()
			return
		}

		// If cookie present, map existing token, else generate a new one.
		if ex, err := r.Cookie(opts.Cookie); err == nil && ex.Value != "" {
			x.Token = ex.Value
		} else {
			x.Token = xsrftoken.Generate(x.Secret, x.ID, "POST")
			if opts.SetCookie {
				setCookie(opts, x, w, r)
			}
		}

		if opts.SetHeader {
			w.Header().Add(opts.Header, x.Token)
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) GetXSRFValidatePipeHandler() PipeHandler {
	return PipeHandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.Handler) {
		s.Logger.Tracef("enter")
		defer s.Logger.Tracef("exit")

		if r.Method == "GET" || r.Method == "HEAD" {
			next.ServeHTTP(w, r)
			return
		}

		ctx := GetContext(r)
		x := ctx.xsrf
		if token := r.Header.Get(x.GetHeaderName()); token != "" {
			if !x.ValidToken(token) {
				x.Error(w)
			}
			next.ServeHTTP(w, r)
			return
		}
		if token := r.FormValue(x.GetFormName()); token != "" {
			if !x.ValidToken(token) {
				x.Error(w)
			}

			next.ServeHTTP(w, r)
			return
		}

		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	})
}
