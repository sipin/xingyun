package xingyun

import (
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
		h.ServeContext(getUnInitedContext(r, w))
	}
}

func FromHTTPHandlerFunc(h http.HandlerFunc) ContextHandlerFunc {
	return func(ctx *Context) {
		h.ServeHTTP(ctx.ResponseWriter, ctx.Request)
	}
}

func ToHTTPHandler(h ContextHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeContext(getUnInitedContext(r, w))
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

	isInited   bool
	flash      *Flash
	staticData map[string][]string
	opts       *Options
	xsrf       XSRF
}

func GetContext(r *http.Request) *Context {
	obj, ok := context.GetOk(r, CONTEXT_KEY)
	if !ok {
		panic("can't get context")
	}
	ctx := obj.(*Context)
	if !ctx.isInited {
		panic("get uninited context")
	}
	return ctx
}

func initContext(r *http.Request, w http.ResponseWriter, s *Server) *Context {
	ctx := getUnInitedContext(r, w)
	if ctx.isInited {
		return ctx
	}
	*ctx = Context{
		ResponseWriter: w,
		Request:        r,
		Server:         s,
		Config:         s.Config,
		Logger:         s.Logger,
		Params:         map[string]string{},
		Data:           map[string]interface{}{},
		staticData:     map[string][]string{},
	}
	ctx.parseParams()
	ctx.isInited = true
	context.Set(r, CONTEXT_KEY, ctx)
	return ctx
}

func getUnInitedContext(r *http.Request, w http.ResponseWriter) *Context {
	ctx, ok := context.GetOk(r, CONTEXT_KEY)
	if !ok {
		newctx := &Context{Request: r, ResponseWriter: w}
		context.Set(r, CONTEXT_KEY, newctx)
		return newctx
	}
	return ctx.(*Context)
}

func (ctx *Context) parseParams() {
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
