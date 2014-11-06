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

type Getter interface {
	Get(name string) (string, error)
}

type Setter interface {
	Set(name, value string) error
}

type GetSetter interface {
	Getter
	Setter
}

type Cookier interface {
	GetSetter
	SetWithExpire(name, value string, expire int)
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

	Logger         Logger
	SessionStorage SessionStorage

	UserID string
	Params map[string]string
	Data   map[string]interface{}
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

func (ctx *Context) WriteString(s string) {
	_, err := ctx.Write([]byte(s))
	if err != nil {
		ctx.Logger.Errorf(err.Error())
		return
	}
}

func (ctx *Context) NotModified() {
	ctx.WriteHeader(304)
}

func (ctx *Context) NotFound(message string) {
	ctx.WriteHeader(404)
	_, err := ctx.Write([]byte(message))
	if err != nil {
		ctx.Logger.Errorf(err.Error())
		return
	}
}

func (ctx *Context) Unauthorized() {
	ctx.WriteHeader(401)
}

func (ctx *Context) Forbidden() {
	ctx.WriteHeader(403)
}
