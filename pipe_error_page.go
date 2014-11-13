package xingyun

import "net/http"

func (s *Server) GetErrorPagePipeHandler() PipeHandler {
	return PipeHandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		s.Logger.Tracef("enter error page handler")
		defer s.Logger.Tracef("exit error page handler")

		next(w, r)
		ctx := GetContext(r)
		status := ctx.ResponseWriter.Status()
		if status >= 400 && status <= 599 {
			if s.ErrorPageHandler != nil && ctx.ResponseWriter.Size() == 0 {
				s.ErrorPageHandler.ServeContext(ctx)
			}
		}
	})
}
