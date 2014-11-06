package xingyun

import (
	"fmt"
	"net/http"

	"github.com/gorilla/context"
)

const (
	CONTEXT_KEY = "_XINGYUN_CONTEXT_"
)

type ContextHandler interface {
	ServeContext(ctx *Context)
}

type ContextHandlerFunc func(ctx *Context)

func (h ContextHandlerFunc) ServeContext(ctx *Context) {
	h(ctx)
}

func ToHTTPHandlerFunc(h ContextHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeContext(GetContext(r))
	}
}

func ToHTTPHandler(h ContextHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeContext(GetContext(r))
	})
}

func FromHTTPHandler(h http.Handler) ContextHandler {
	return ContextHandlerFunc(func(ctx *Context) {
		h.ServeHTTP(ctx.ResponseWriter, ctx.Request)
	})
}

type SessionStorage interface {
	SetSession(sessionID string, key string, data []byte)
	GetSession(sessionID string, key string) []byte
	ClearSession(sessionID string, key string)
}

type Context struct {
	http.ResponseWriter
	Request *http.Request
	Server  *Server
	Config  *Config

	Logger Logger

	UserID string
	Params map[string]string

	// use for user ContextHandler
	Data map[string]interface{}
	// use for user PipeHandler. avoid name conflict
	PipeHandlerData map[string]interface{}

	flash      *Flash
	staticData map[string][]string
	opts       *Options
	xsrf       XSRF
}

func NewContext(r *http.Request, w http.ResponseWriter, s *Server) *Context {
	ctx := &Context{
		ResponseWriter: w,
		Request:        r,
		Server:         s,
		Config:         s.Config,
		Logger:         s.Logger,
		Params:         map[string]string{},
		Data:           map[string]interface{}{},
		staticData:     map[string][]string{},
	}
	context.Set(r, CONTEXT_KEY, ctx)
	return ctx
}

func GetContext(r *http.Request) *Context {
	ctx, ok := context.GetOk(r, CONTEXT_KEY)
	if !ok {
		panic(fmt.Errorf("can't get context"))
	}
	return ctx.(*Context)
}

func (ctx *Context) parseParams(s string) {
	var err error
	err = ctx.Request.ParseForm()
	if err != nil {
		ctx.Logger.Errorf(err.Error())
		return
	}
	for k, v := range ctx.Request.Form {
		ctx.Params[k] = v[0]
	}
}
