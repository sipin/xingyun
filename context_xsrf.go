package xingyun

import "fmt"

func (ctx *Context) checkXSRF() {
	if ctx.xsrf == nil {
		panic(fmt.Errorf("xsrf not inited"))
	}
}

func (ctx *Context) XSRFToken() string {
	ctx.checkXSRF()
	if ctx.xsrf == nil {
		return ""
	}
	return ctx.xsrf.GetToken()
}

func (ctx *Context) XSRFFormField() string {
	ctx.checkXSRF()
	return "<input type=\"hidden\" name=\"" + formatName(ctx.xsrf.GetFormName()) + "\" value=\"" +
		ctx.XSRFToken() + "\"/>"
}
