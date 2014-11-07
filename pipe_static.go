package xingyun

import (
	"net/http"
	"path"
	"strings"
)

func (s *Server) GetStaticPipeHandler() PipeHandler {
	return PipeHandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.Handler) {
		cfg := s.Config
		if r.Method != "GET" && r.Method != "HEAD" {
			next.ServeHTTP(rw, r)
			return
		}
		file := r.URL.Path
		if cfg.StaticPrefix != "" {
			if !strings.HasPrefix(file, cfg.StaticPrefix) {
				next.ServeHTTP(rw, r)
				return
			}
			file = file[len(cfg.StaticPrefix):]
			if file != "" && file[0] != '/' {
				next.ServeHTTP(rw, r)
				return
			}
		}
		f, err := s.StaticDir.Open(file)
		if err != nil {
			next.ServeHTTP(rw, r)
			return
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			next.ServeHTTP(rw, r)
			return
		}

		if fi.IsDir() {
			if !strings.HasSuffix(r.URL.Path, "/") {
				http.Redirect(rw, r, r.URL.Path+"/", http.StatusFound)
				return
			}

			file = path.Join(file, cfg.StaticIndexFile)
			f, err = s.StaticDir.Open(file)
			if err != nil {
				next.ServeHTTP(rw, r)
				return
			}
			defer f.Close()

			fi, err = f.Stat()
			if err != nil || fi.IsDir() {
				next.ServeHTTP(rw, r)
				return
			}
		}

		http.ServeContent(rw, r, file, fi.ModTime(), f)
	})
}
