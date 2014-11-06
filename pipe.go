package xingyun

import "net/http"

type PipeHandler interface {
	ServePipe(w http.ResponseWriter, r *http.Request, next http.Handler)
}

type PipeHandlerFunc func(w http.ResponseWriter, r *http.Request, next http.Handler)

func (h PipeHandlerFunc) ServePipe(w http.ResponseWriter, r *http.Request, next http.Handler) {
	h(w, r, next)
}

func Wrap(h PipeHandler, f ContextHandlerFunc) ContextHandlerFunc {
	return func(ctx *Context) {
		h.ServePipe(ctx.ResponseWriter, ctx.Request, ToHTTPHandler(f))
	}
}

type Pipe struct {
	Server   *Server
	Handlers []PipeHandler
}

func NewPipe(server *Server, handlers ...PipeHandler) *Pipe {
	pipe := &Pipe{Server: server}
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
		sub := &Pipe{Server: p.Server, Handlers: p.Handlers[1:]}
		handler.ServePipe(w, r, sub.HTTPHandler(h))
	}
}

var voidHTTPHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

func (p *Pipe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.ServePipe(w, r, voidHTTPHandler)
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
		sub := &Pipe{Server: p.Server, Handlers: p.Handlers[1:]}
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
	if len(p.Server.DefaultPipeHandlers) != 0 {
		p.Handlers = p.Server.DefaultPipeHandlers
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
