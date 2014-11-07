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
	router *mux.Router
}

func newRouter() Router {
	gorillaRouter := mux.NewRouter()
	return &router{router: gorillaRouter}
}

func (r *router) Get(path string, h ContextHandlerFunc) {
	r.router.Handle(path, ToHTTPHandler(h)).Methods("GET")
}

func (r *router) Post(path string, h ContextHandlerFunc) {
	r.router.Handle(path, ToHTTPHandler(h)).Methods("POST")
}

func (r *router) Handle(path string, h ContextHandler) {
	r.router.Handle(path, ToHTTPHandler(h))
}

func (r *router) HandleFunc(path string, h ContextHandlerFunc) {
	r.router.HandleFunc(path, ToHTTPHandlerFunc(h))
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
