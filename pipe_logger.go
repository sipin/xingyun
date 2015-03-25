package xingyun

import (
	"net/http"
	"time"
)

func (s *Server) GetLogPipeHandler() PipeHandler {
	return PipeHandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		s.logger.Tracef("enter logger")
		defer s.logger.Tracef("exit logger")

		start := time.Now()
		next(w, r)

		res := w.(ResponseWriter)
		log := s.logger.Infof
		status := res.Status()
		if status >= 500 && status <= 599 {
			log = s.logger.Errorf
		}
		if status >= 400 && status <= 499 {
			log = s.logger.Warnf
		}
		log("%v %s %s %s in %v", res.Status(), r.Method, r.Host, r.URL.Path, time.Since(start))
	})
}
