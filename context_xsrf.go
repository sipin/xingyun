package xingyun

func (ctx *Context) XSRFToken() string {
	if ctx.xsrf == nil {
		return ""
	}
	return ctx.xsrf.GetToken()
}

func (ctx *Context) XSRFFormField() string {
	return "<input type=\"hidden\" name=\"" + formatName(ctx.xsrf.GetFormName()) + "\" value=\"" +
		ctx.XSRFToken() + "\"/>"
}
