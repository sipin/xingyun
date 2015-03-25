package xingyun

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) getURLVarLoaderPipeHandler() PipeHandler {
	return PipeHandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		s.logger.Tracef("enter urlVar loader")
		defer s.logger.Tracef("exit urlVal loader")

		ctx := GetContext(r)
		urlVars := mux.Vars(r)
		for k, v := range urlVars {
			_, ok := ctx.Params[k]
			if ok {
				s.logger.Warnf("param %s is overide by urlVar", k)
			}
			ctx.Params[k] = v
			s.logger.Debugf("load urlVal: %s = %s", k, v)
		}

		next(w, r)
	})
}
