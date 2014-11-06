package xingyun

import "net/http"

type ContextHandler interface {
	ServeContext(ctx *Context)
}

type ContextHandlerFunc func(ctx *Context)

func (h ContextHandlerFunc) ServeContext(ctx *Context) {
	h(ctx)
}

func GetContext(r *http.Request) *Context {
	return nil
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

	Logger         Logger
	SessionStorage SessionStorage

	UserID string
	Params map[string]string
	Data   map[string]interface{}

	flash      *Flash
	staticData map[string][]string
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
