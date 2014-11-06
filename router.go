package xingyun

import "net/http"

type Router interface {
	Handle(pattern string, h Context)
	HandleFunc(pattern string, h ContextHandlerFunc)
	Get(pattern string, h ContextHandlerFunc)
	Post(pattern string, h ContextHandlerFunc)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
