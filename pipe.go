package xingyun

import (
	"net/http"

	"github.com/codegangsta/negroni"
)

var (
	DefaultPipeHandler = []PipeHandler{}
)

type PipeHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type Pipe struct {
	negroni *negroni.Negroni
}

func (p *Pipe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.negroni.ServeHTTP(w, r)
}

func (p *Pipe) Use(handlers ...PipeHandler) {
	if len(DefaultPipeHandler) != 0 && len(p.negroni.Handlers()) == 0 {
		for _, h := range DefaultPipeHandler {
			p.negroni.Use(h)
		}
	}
	for _, h := range handlers {
		p.negroni.Use(h)
	}
}
