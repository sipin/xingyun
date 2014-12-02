## XingYun ##

XingYun is a web framework with negroni like middleware and web.go like API.

Difference compare to negroni and web.go


- negroni: XingYun have Context that wrap Request and ResponseWriter. The Context have a easy to use API. XingYun have pipe to manage middleware.

- web.go: All feature in XingYun is composed by middleware and every middleware can replaced by user. The middleware also can do other pre-request or post-request features.    


### Object ###

Important object

- Server: Server contains everything and implement router. Server also is a pipe container to manage pipe.
- Pipe: PipeHandler queue. Pipe is also a PipeHandler and http.Handler
- Context: contain request related data and have easy to use API.

### Interface ###

Interface in XingYun

- ContextHandler:
  process request with context

		type ContextHandler interface {
			ServeContext(ctx *Context)
		}

- PipeHandler (like negroni.Handler)
- Logger: used to print log
- Router: set route

### Context ###

context support features below

- Cookie
- Session
- XSRF
- Flash
- Render

### Default PipeHandler ###

- Logger
- Recover
- Static
- XSRF
- ErrorPage
- Context
- URLVarLoader
