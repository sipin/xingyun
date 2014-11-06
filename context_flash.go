package xingyun

type Flash struct {
	Alert  string
	Notice string
}

var (
	FlashExpire    int64  = 60
	FlashAlertKey  string = "ZQFA"
	FlashNoticeKey string = "ZQFN"
)

func (ctx *Context) SetFlashAlert(msg string) {
	ctx.SetExpireCookie(FlashAlertKey, msg, FlashExpire)
}

func (ctx *Context) SetFlashNotice(msg string) {
	ctx.SetExpireCookie(FlashNoticeKey, msg, FlashExpire)
}

func (ctx *Context) getFlash() *Flash {
	flash := &Flash{}
	var err error
	flash.Alert, err = ctx.GetStringCookie(FlashAlertKey)
	if err == nil {
		ctx.RemoveCookie(FlashAlertKey)
	}

	flash.Notice, err = ctx.GetStringCookie(FlashNoticeKey)
	if err == nil {
		ctx.RemoveCookie(FlashNoticeKey)
	}
	return flash
}

func (ctx *Context) GetFlash() *Flash {
	if ctx.flash != nil {
		return ctx.flash
	}
	ctx.flash = ctx.getFlash()
	return ctx.flash
}
