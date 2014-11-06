package xingyun

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"code.1dmy.com/xyz/logex"
)

func (s *Server) GetRecoverPipeHandler() PipeHandler {
	return PipeHandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.Handler) {
		s.Logger.Tracef("enter")
		defer s.Logger.Tracef("exit")

		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				var stacks []string
				for i := 1; ; i += 1 {
					_, file, line, ok := runtime.Caller(i)
					if !ok {
						break
					}
					stacks = append(stacks, fmt.Sprintf("\t%s:%d", file, line))
				}
				stackMsg := strings.Join(stacks, "\n")
				logex.Errorf("%s %s: %s\n%s", r.Method, r.RequestURI, err, stackMsg)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
