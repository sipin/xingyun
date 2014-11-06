package xingyun

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"strings"
)

func (ctx *Context) GetStaticUrl(url string) string {
	cfg := ctx.Config
	if url[0] != '/' {
		return url
	}
	if cfg.StaticHost == "" {
		return url
	}

	if ctx.isExcludeFile(url) {
		return url
	}

	if ctx.isExcludeType(url) {
		return url
	}

	hash := ctx.getStaticFileHash(url)
	if hash == "" {
		return cfg.StaticHost + url
	}

	if cfg.StaticHost == "/" {
		return url + "?hash=" + hash
	}
	return cfg.StaticHost + url + "?hash=" + hash
}

func isIE(version, user_agent string) bool {
	key := "MSIE " + version
	return strings.Contains(user_agent, key)
}

var browserMatchMap = map[string]func(user_agent string) bool{
	"ie8": func(user_agent string) bool { return isIE("8.0", user_agent) },
	"ie9": func(user_agent string) bool { return isIE("9.0", user_agent) },
}

func browserMatch(t, user_agent string) bool {
	if user_agent == "" {
		return false
	}
	if t == "" {
		return true
	}
	handler, ok := browserMatchMap[t]
	if !ok {
		return false
	}
	return handler(user_agent)
}

func splitExclude(s string) (t, v string) {
	if !strings.Contains(s, ":") {
		return "", s
	}
	ss := strings.Split(s, ":")
	return ss[0], ss[1]
}

func (ctx *Context) GetUserAgent() string {
	agents := ctx.Request.Header["User-Agent"]
	if len(agents) > 0 {
		return agents[0]
	}
	return ""
}

func (ctx *Context) isExcludeType(url string) bool {
	types := strings.Split(ctx.Config.StaticHostExcludeType, ",")
	for _, t := range types {
		if t == "" {
			continue
		}
		t1, v := splitExclude(t)
		agent := ctx.GetUserAgent()
		if browserMatch(t1, agent) && strings.HasSuffix(url, v) {
			return true
		}
	}
	return false
}

func (ctx *Context) isExcludeFile(url string) bool {
	files := strings.Split(ctx.Config.StaticHostExcludeFile, ",")
	for _, f := range files {
		if f == "" {
			continue
		}
		t1, v := splitExclude(f)
		if browserMatch(t1, ctx.GetUserAgent()) && url == v {
			return true
		}
	}
	return false
}

func (ctx *Context) getStaticFileHash(name string) string {
	dir := ctx.Server.StaticDir
	if dir == nil {
		return ""
	}
	f, err := dir.Open(name)
	if err != nil {
		return ""
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return ""
	}
	if stat.IsDir() {
		return ""
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", md5.Sum(data))
}

func (ctx *Context) addDataUnique(dataType string, values ...string) {
	var data []string
	var ok bool
	if ctx.staticData == nil {
		ctx.staticData = make(map[string][]string)
	}

	if data, ok = ctx.staticData[dataType]; !ok {
		data = []string{}
	}

	for _, val := range values {
		isNewVal := true
		for _, ele := range data {
			if ele == val {
				isNewVal = false
			}
		}

		if isNewVal {
			data = append(data, val)
		}
	}
	ctx.staticData[dataType] = data
}

func (ctx *Context) getData(dataType string) (values []string) {
	return ctx.staticData[dataType]
}

func (ctx *Context) AddJS(val ...string) {
	ctx.addDataUnique("js", val...)
}

func (ctx *Context) GetJS() (values []string) {
	return ctx.getData("js")
}

func (ctx *Context) AddCSS(val ...string) {
	ctx.addDataUnique("css", val...)
}

func (ctx *Context) GetCSS() (values []string) {
	return ctx.getData("css")
}
