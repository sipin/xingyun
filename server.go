package xingyun

import (
	"net/http"

	"code.1dmy.com/xyz/logex"
	"github.com/gorilla/securecookie"
)

type Config struct {
	CookieSecret          string
	StaticDir             string
	StaticHost            string
	StaticHostExcludeType string
	StaticHostExcludeFile string
}

type Server struct {
	Router
	Config    *Config
	StaticDir http.FileSystem

	Name         string
	Logger       Logger
	SecureCookie *securecookie.SecureCookie
}

func NewServer(config *Config) *Server {
	server := &Server{
		Router: NewRouter(),
		Logger: logex.NewLogger(1),
	}
	server.StaticDir = http.Dir(config.StaticDir)
	server.SecureCookie = securecookie.New([]byte(config.CookieSecret), []byte(config.CookieSecret))
	return server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

func (s *Server) name() string {
	if s.Name == "" {
		return "xingyun"
	}
	return s.Name
}

func (s *Server) ListenAndServe(addr string) error {
	s.Logger.Infof("%s start on %s", s.name(), addr)
	err := http.ListenAndServe(addr, s)
	s.Logger.Errorf("%s stop, err='%s'", err)
	return err
}
