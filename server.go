package xingyun

import (
	"fmt"
	"net/http"

	"code.1dmy.com/xyz/logex"
	"github.com/gorilla/securecookie"
)

type Server struct {
	Router
	Config    *Config
	StaticDir http.FileSystem

	Name                string
	Logger              Logger
	SecureCookie        *securecookie.SecureCookie
	DefaultPipeHandlers []PipeHandler
	StaticHandler       http.Handler

	pipes map[string]*Pipe
}

func NewServer(config *Config) *Server {
	if config == nil {
		config = &Config{}
	}
	setDefaultConfig(config)
	server := &Server{
		Logger: logex.NewLogger(1),
		Config: config,
	}

	server.pipes = map[string]*Pipe{}
	server.DefaultPipeHandlers = []PipeHandler{
		server.GetLogPipeHandler(),
		server.GetRecoverPipeHandler(),
		server.GetContextPipeHandler(),
	}

	server.StaticHandler = server.GetStaticHandler()
	server.Router = NewRouter(server.StaticHandler)
	server.StaticDir = http.Dir(config.StaticDir)

	server.SecureCookie = securecookie.New(
		[]byte(config.CookieSecret),
		[]byte(config.CookieSecret),
	)

	return server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(NewResponseWriter(w), r)
}

func (s *Server) NewPipe(name string, handlers ...PipeHandler) *Pipe {
	p := newPipe(s, handlers...)
	_, ok := s.pipes[name]
	if ok {
		panic(fmt.Errorf("pipe %s is exist", name))
	}
	s.pipes[name] = p
	return p
}

func (s *Server) Pipe(name string) *Pipe {
	p, ok := s.pipes[name]
	if !ok {
		panic(fmt.Errorf("pipe %s is not exist", name))
	}
	return p
}

func (s *Server) name() string {
	if s.Name == "" {
		return "xingyun"
	}
	return s.Name
}

func (s *Server) ListenAndServe(addr string) error {
	s.Logger.Infof("%s start in: %s", s.name(), addr)
	err := http.ListenAndServe(addr, s)
	s.Logger.Errorf("%s stop, err='%s'", err)
	return err
}
