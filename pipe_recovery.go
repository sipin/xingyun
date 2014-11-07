package xingyun

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"code.1dmy.com/xyz/logex"
)

func DefaultPanicHandler(ctx *Context) {
	w := ctx.ResponseWriter
	r := ctx.Request
	w.WriteHeader(http.StatusInternalServerError)

	logex.Errorf("%s %s: %s\n%s", r.Method, r.RequestURI, ctx.PanicError, ctx.StackMessage)
}

func (s *Server) GetRecoverPipeHandler() PipeHandler {
	return PipeHandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.Handler) {
		s.Logger.Tracef("enter recover handler")

		defer func() {
			if err := recover(); err != nil {
				ctx := GetContext(r)
				ctx.IsPanic = true
				ctx.PanicError = err
				var stacks []string
				for i := 1; ; i += 1 {
					_, file, line, ok := runtime.Caller(i)
					if !ok {
						break
					}
					stacks = append(stacks, fmt.Sprintf("\t%s:%d", file, line))
				}
				ctx.StackMessage = strings.Join(stacks, "\n")
				s.PanicHandler(ctx)
			}
			s.Logger.Tracef("exit recover handler")
		}()

		next.ServeHTTP(w, r)
	})
}
