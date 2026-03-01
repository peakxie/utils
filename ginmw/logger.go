package ginmw

import "sync"

// Logger defines a minimal logging interface.
// Users must inject their own implementation via SetLogger.
//
// Adapting logrus:
//
//	type logrusAdapter struct{ *logrus.Logger }
//	func (l logrusAdapter) Infof(format string, args ...any)  { l.Logger.Infof(format, args...) }
//	func (l logrusAdapter) Errorf(format string, args ...any) { l.Logger.Errorf(format, args...) }
//	ginmw.SetLogger(logrusAdapter{logrus.StandardLogger()})
//
// Adapting zap:
//
//	type zapAdapter struct{ *zap.SugaredLogger }
//	func (l zapAdapter) Infof(format string, args ...any)  { l.SugaredLogger.Infof(format, args...) }
//	func (l zapAdapter) Errorf(format string, args ...any) { l.SugaredLogger.Errorf(format, args...) }
//	ginmw.SetLogger(zapAdapter{sugar})
type Logger interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
}

var (
	logger Logger = nopLogger{}
	logMu  sync.RWMutex
)

// SetLogger sets the global logger used by ginmw.
// If l is nil, logging is disabled (nop logger).
func SetLogger(l Logger) {
	logMu.Lock()
	defer logMu.Unlock()
	if l == nil {
		logger = nopLogger{}
		return
	}
	logger = l
}

func getLogger() Logger {
	logMu.RLock()
	defer logMu.RUnlock()
	return logger
}

type nopLogger struct{}

func (nopLogger) Infof(string, ...any)  {}
func (nopLogger) Errorf(string, ...any) {}
