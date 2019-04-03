package belt

type Logger interface {
	Debug(...interface{})
	Debugf(string, ...interface{})

	Info(...interface{})
	Infof(string, ...interface{})

	Warn(...interface{})
	Warnf(string, ...interface{})

	Error(...interface{})
	Errorf(string, ...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})

	Panic(...interface{})
	Panicf(string, ...interface{})

	WithFields(map[string]string) Logger
	WithField(string, string) Logger
}
