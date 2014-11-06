package xingyun

import (
	"net/http"

	"github.com/gorilla/context"
)

func (s *Server) GetContextPipeHandler() PipeHandler {
	return PipeHandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.Handler) {
		s.Logger.Tracef("enter")
		defer s.Logger.Tracef("exit")

		initContext(r, w, s)
		defer context.Clear(r)
		next.ServeHTTP(w, r)
	})
}
