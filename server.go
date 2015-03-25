package xingyun

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/securecookie"
)

type Server struct {
	Router
	Config    *Config
	StaticDir http.FileSystem

	Name                string
	SecureCookie        *securecookie.SecureCookie
	DefaultPipeHandlers []PipeHandler
	PanicHandler        ContextHandlerFunc
	ErrorPageHandler    ContextHandlerFunc
	SessionStorage      SessionStorage

	pipes       map[string]*Pipe
	defaultPipe *Pipe
	l           net.Listener
	logger      Logger
}

func NewServer(config *Config) *Server {
	if config == nil {
		config = &Config{}
	}
	setDefaultConfig(config)
	server := &Server{
		logger: &debugLogger{Logger: NewSimpleLevelLogger(os.Stdout), enableDebug: config.EnableDebug},
		Config: config,
	}
	server.PanicHandler = DefaultPanicHandler
	server.pipes = map[string]*Pipe{}
	server.DefaultPipeHandlers = []PipeHandler{
		server.GetErrorPagePipeHandler(),
		server.GetLogPipeHandler(),
		server.GetRecoverPipeHandler(),
		server.GetStaticPipeHandler(),
	}

	server.Router = newRouter(server.getURLVarLoaderPipeHandler())
	server.StaticDir = http.Dir(config.StaticDir)
	server.SessionStorage = NewMemoryStore()

	server.SecureCookie = securecookie.New(
		[]byte(config.CookieSecret),
		[]byte(config.CookieSecret),
	)

	return server
}

func (s *Server) SetLogger(l Logger) {
	s.logger = &debugLogger{Logger: l, enableDebug: s.Config.EnableDebug}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pipeHandlers := []PipeHandler{s.GetContextPipeHandler()}
	pipeHandlers = append(pipeHandlers, s.DefaultPipeHandlers...)
	s.defaultPipe = newPipe(s, pipeHandlers...)
	h := s.defaultPipe.HTTPHandler(s.Router)
	h.ServeHTTP(wrapResponseWriter(w), r)
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
	l, err := net.Listen("tcp", addr)
	if err != nil {
		s.logger.Errorf(err.Error())
	}
	s.l = l
	s.logger.Infof("%s start in: %s", s.name(), addr)
	err = http.Serve(s.l, s)
	// todo: must handle error when serve failed
	// s.logger.Errorf("%s stop, err='%s'", err)
	return err
}

func (s *Server) Stop() error {
	if s.l != nil {
		return s.l.Close()
	}
	return nil
}
