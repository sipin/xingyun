package xingyun

import "net/http"

var (
	DefaultPipeHandlers = []PipeHandler{}
)

type PipeHandler interface {
	ServePipe(w http.ResponseWriter, r *http.Request, next http.Handler)
}

type Pipe struct {
	Handlers []PipeHandler
}

func NewPipe(handlers ...PipeHandler) *Pipe {
	pipe := &Pipe{}
	pipe.Use(handlers...)
	return pipe
}

func (p *Pipe) ServePipe(w http.ResponseWriter, r *http.Request, h http.Handler) {
	switch len(p.Handlers) {
	case 0:
		h.ServeHTTP(w, r)
	case 1:
		handler := p.Handlers[0]
		handler.ServePipe(w, r, h)
	default:
		handler := p.Handlers[0]
		sub := &Pipe{p.Handlers[1:]}
		handler.ServePipe(w, r, sub.HTTPHandler(h))
	}
}

var VoidHTTPHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

func (p *Pipe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.ServePipe(w, r, VoidHTTPHandler)
}

func (p *Pipe) ServeContext(ctx *Context) {
	p.ServeHTTP(ctx.ResponseWriter, ctx.Request)
}

func (p *Pipe) HTTPHandler(h http.Handler) http.Handler {
	switch len(p.Handlers) {
	case 0:
		return h
	case 1:
		handler := p.Handlers[0]
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.ServePipe(w, r, h)
		})
	default:
		handler := p.Handlers[0]
		sub := &Pipe{p.Handlers[1:]}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.ServePipe(w, r, sub.HTTPHandler(h))
		})
	}
}

func (p *Pipe) ContextHandler(h ContextHandler) ContextHandler {
	httpHandler := ToHTTPHandler(h)
	return FromHTTPHandler(p.HTTPHandler(httpHandler))
}

func (p *Pipe) Use(handlers ...PipeHandler) {
	if len(DefaultPipeHandlers) != 0 {
		p.Handlers = DefaultPipeHandlers
	}
	for _, h := range handlers {
		p.Handlers = append(p.Handlers, h)
	}
}

func (p *Pipe) Wrap(h ContextHandlerFunc) ContextHandlerFunc {
	return ContextHandlerFunc(func(ctx *Context) {
		wrapHandler := p.ContextHandler(h)
		wrapHandler.ServeContext(ctx)
	})
}
