package xingyun

import (
	"net/http"
	"time"
)

func (s *Server) GetLogPipeHandler() PipeHandler {
	return PipeHandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.Handler) {
		s.Logger.Tracef("enter logger")
		defer s.Logger.Tracef("exit logger")

		start := time.Now()
		next.ServeHTTP(w, r)

		res := w.(ResponseWriter)
		log := s.Logger.Infof
		status := res.Status()
		if status >= 500 && status <= 599 {
			log = s.Logger.Errorf
		}
		if status >= 400 && status <= 499 {
			log = s.Logger.Warnf
		}
		log("%v %s %s %s in %v", res.Status(), r.Method, r.Host, r.URL.Path, time.Since(start))
	})
}
