## 行云 ##

### 对象 ###

- Server: 包含处理请求的所有相关组件。继承Router，用于设置路由。Pipe容器，管理Pipe
- Pipe: PipeHandler队列
- Context: 请求相关的数据

### 接口 ###

- ContextHandler:

		type ContextHandler interface {
			ServeContext(ctx *Context)
		}


- PipeHandler (类似negroni.Handler)
- Config
- Logger: 统一LOG格式
- Router: 设置路由

### 目标 ###

- 使用接口类似web.go
- 提供类似negroni的web中间件机制

### Context ###

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
