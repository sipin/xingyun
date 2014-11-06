package xingyun

import (
	"net/http"
	"time"

	"code.1dmy.com/xyz/logex"

	"github.com/codegangsta/negroni"
)

func (s *Server) GetLogPipeHandler() PipeHandler {
	return PipeHandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.Handler) {
		s.Logger.Tracef("enter")
		defer s.Logger.Tracef("exit")

		start := time.Now()
		next.ServeHTTP(rw, r)

		res := rw.(negroni.ResponseWriter)
		log := logex.Infof
		status := res.Status()
		if status >= 500 && status <= 599 {
			log = logex.Errorf
		}
		log("%v %s %s %s in %v", res.Status(), r.Method, r.Host, r.URL.Path, time.Since(start))
	})
}
