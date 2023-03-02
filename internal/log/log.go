package log

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
)

var l logr.Logger = stdr.New(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile))
var lock = &sync.RWMutex{}

// SetLogger ...
func SetLogger(log logr.Logger) {
	lock.Lock()
	defer lock.Unlock()
	l = log
}

/*
// Info ...
func Info(msg string, keysAndValues ...interface{}) {
	lock.RLock()
	defer lock.RUnlock()
	l.Info(msg, keysAndValues...)
}

// Error ...
func Error(err error, msg string, keysAndValues ...interface{}) {
	lock.RLock()
	defer lock.RUnlock()
	l.Error(err, msg, keysAndValues...)
}
*/

func getMsg(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

// Infof ...
func Infof(format string, args ...interface{}) {
	lock.RLock()
	defer lock.RUnlock()
	l.Info(getMsg(format, args...))
}

// Errorf ...
func Errorf(format string, args ...interface{}) {
	lock.RLock()
	defer lock.RUnlock()
	l.Info(getMsg(format, args...))
}
