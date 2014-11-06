package xingyun

import (
	"net/http"

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
	SecureCookie securecookie.SecureCookie
}

func NewServer(router Router, logger Logger) *Server {
	return &Server{Router: router, Logger: logger}
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
