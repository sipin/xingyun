package xingyun

import "net/http"

type Server struct {
	Config Config

	pipeRouter Router
}

func NewServer(router Router) *Server {
	return &Server{pipeRouter: router}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := RouterHandler(s.pipeRouter)
	h(w, r)
}

func (s *Server) ListenAndServe(addr string) {
	http.ListenAndServe(addr, s)
}
