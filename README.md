## 行云 ##

### 对象 ###

- Server: pipe容器
- Pipe: 对请求做前、后置处理
- Context: 请求相关的数据

### 接口 ###

- http.Handler
- pipe.Handler (negroni.Handler)
- Config
- Logger
- Router: 用于Pipe路由Handler路由

### 目标 ###

- 提供机制、而不**只是**实现
- 方便扩展
- 默认配置功能强大、容易使用
- 使用代码生成作为元编程
- Everything is http.Handler

### Core SubContext ###

- Cookie
- Session
- XSRF
- Flash
- Render

### Default PipeHandler ###

- Logger
- Recover
- Static