package xingyun

import (
	"net/http"

	"github.com/gorilla/context"
)

func (s *Server) GetContextPipeHandler() PipeHandler {
	return PipeHandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.Handler) {
		NewContext(r, w, s)
		defer context.Clear(r)
		next.ServeHTTP(w, r)
	})
}
