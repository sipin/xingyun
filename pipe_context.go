package xingyun

import (
	"net/http"

	"github.com/gorilla/context"
)

func (s *Server) GetContextPipeHandler() PipeHandler {
	return PipeHandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.Handler) {
		s.Logger.Tracef("enter context handler")
		defer s.Logger.Tracef("exit context handler")

		ctx := initContext(r, w, s)
		next.ServeHTTP(w, r)
		ctx.Logger.Debugf("clear context, &r=%p", r)
		context.Clear(r)
	})
}
