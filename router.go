package xingyun

import "net/http"

type Router interface {
	Route(r *http.Request) http.Handler
}

func RouterHandler(router Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h := router.Route(r)
		h.ServeHTTP(w, r)
	}
}
