package xingyun

type Logger interface {
	Infof(s string, o ...interface{})
	Errorf(s string, o ...interface{})
	Warnf(s string, o ...interface{})
	Debugf(s string, o ...interface{})
	Tracef(s string, o ...interface{})
}

type debugLogger struct {
	Logger
	enableDebug bool
}

func (l *debugLogger) Debugf(s string, o ...interface{}) {
	if l.enableDebug {
		l.Logger.Debugf(s, o...)
	}
}

func (l *debugLogger) Tracef(s string, o ...interface{}) {
	if l.enableDebug {
		l.Logger.Tracef(s, o...)
	}
}
