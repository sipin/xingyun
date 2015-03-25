package xingyun

import (
	"io"
	"log"
)

type Logger interface {
	Infof(s string, o ...interface{})
	Errorf(s string, o ...interface{})
	Debugf(s string, o ...interface{})
	Warnf(s string, o ...interface{})
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

type simpleLevelLogger struct {
	l *log.Logger
}

func NewSimpleLevelLogger(w io.Writer) *simpleLevelLogger {
	return &simpleLevelLogger{log.New(w, "", log.LstdFlags)}
}

func (l *simpleLevelLogger) output(level, s string, o ...interface{}) {
	log.Printf(level+" "+s, o...)
}

func (l *simpleLevelLogger) Infof(s string, o ...interface{}) {
	l.output("INFO", s, o...)
}

func (l *simpleLevelLogger) Errorf(s string, o ...interface{}) {
	l.output("ERROR", s, o...)
}

func (l *simpleLevelLogger) Warnf(s string, o ...interface{}) {
	l.output("WARN", s, o...)
}

func (l *simpleLevelLogger) Debugf(s string, o ...interface{}) {
	l.output("DEBUG", s, o...)
}

func (l *simpleLevelLogger) Tracef(s string, o ...interface{}) {
	l.output("TRACE", s, o...)
}
