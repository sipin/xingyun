package xingyun

import (
	"mime"
	"strings"
)

func (ctx *Context) SetContentType(val string) string {
	var ctype string
	if strings.ContainsRune(val, '/') {
		ctype = val
	} else {
		if !strings.HasPrefix(val, ".") {
			val = "." + val
		}
		ctype = mime.TypeByExtension(val)
	}
	if ctype != "" {
		ctx.ResponseWriter.Header().Set("Content-Type", ctype)
	}
	return ctype
}
func (ctx *Context) WriteString(s string) {
	_, err := ctx.Write([]byte(s))
	if err != nil {
		ctx.Logger.Errorf(err.Error())
		return
	}
}

func (ctx *Context) NotModified() {
	ctx.WriteHeader(304)
}

func (ctx *Context) NotFound(message string) {
	ctx.WriteHeader(404)
	_, err := ctx.Write([]byte(message))
	if err != nil {
		ctx.Logger.Errorf(err.Error())
		return
	}
}

func (ctx *Context) Unauthorized() {
	ctx.WriteHeader(401)
}

func (ctx *Context) Forbidden() {
	ctx.WriteHeader(403)
}
