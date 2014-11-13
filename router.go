package xingyun

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Router interface {
	Handle(pattern string, h ContextHandler)
	HandleFunc(pattern string, h ContextHandlerFunc)
	Get(pattern string, h ContextHandlerFunc)
	Post(pattern string, h ContextHandlerFunc)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type router struct {
	router     *mux.Router
	afterRoute PipeHandler
}

func newRouter(afterRoute PipeHandler) *router {
	gorillaRouter := mux.NewRouter()
	gorillaRouter.KeepContext = true
	gorillaRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})
	return &router{router: gorillaRouter, afterRoute: afterRoute}
}

func (r *router) getWrapHandler(h ContextHandlerFunc) http.Handler {
	return ToHTTPHandlerFunc(Wrap(r.afterRoute, h))
}

func (r *router) Get(path string, h ContextHandlerFunc) {
	r.router.Handle(path, r.getWrapHandler(h)).Methods("GET")
}

func (r *router) Post(path string, h ContextHandlerFunc) {
	r.router.Handle(path, r.getWrapHandler(h)).Methods("POST")
}

func (r *router) Handle(path string, h ContextHandler) {
	r.router.Handle(path, r.getWrapHandler(h.ServeContext))
}

func (r *router) HandleFunc(path string, h ContextHandlerFunc) {
	r.router.HandleFunc(path, r.getWrapHandler(h).ServeHTTP)
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
