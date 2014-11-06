package xingyun

type logger interface {
	Infof(s string, o ...interface{})
	Errorf(s string, o ...interface{})
	Debugf(s string, o ...interface{})
	Tracef(s string, o ...interface{})
}
