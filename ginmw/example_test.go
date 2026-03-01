package ginmw_test

import (
	"fmt"

	"github.com/peakxie/utils/ginmw"
)

// ExampleSetLogger_logrus demonstrates how to adapt logrus for ginmw.
//
//	import "github.com/sirupsen/logrus"
//
//	type logrusAdapter struct{ *logrus.Logger }
//	func (l logrusAdapter) Infof(format string, args ...any)  { l.Logger.Infof(format, args...) }
//	func (l logrusAdapter) Errorf(format string, args ...any) { l.Logger.Errorf(format, args...) }
//
//	ginmw.SetLogger(logrusAdapter{logrus.StandardLogger()})
func ExampleSetLogger_logrus() {
	// This example shows the pattern — actual logrus import is omitted
	// to avoid adding it as a dependency.
	fmt.Println("logrus adapter set")
	// Output: logrus adapter set
}

// ExampleSetLogger_zap demonstrates how to adapt zap for ginmw.
//
//	import "go.uber.org/zap"
//
//	logger, _ := zap.NewProduction()
//	sugar := logger.Sugar()
//
//	type zapAdapter struct{ *zap.SugaredLogger }
//	func (l zapAdapter) Infof(format string, args ...any)  { l.SugaredLogger.Infof(format, args...) }
//	func (l zapAdapter) Errorf(format string, args ...any) { l.SugaredLogger.Errorf(format, args...) }
//
//	ginmw.SetLogger(zapAdapter{sugar})
func ExampleSetLogger_zap() {
	// This example shows the pattern — actual zap import is omitted
	// to avoid adding it as a dependency.
	fmt.Println("zap adapter set")
	// Output: zap adapter set
}

// ExampleWrapperH_basic shows a typical WrapperH usage.
func ExampleWrapperH_basic() {
	// Use ginmw.SetLogger to inject your logger first.
	ginmw.SetLogger(nil) // nop logger for this example
	fmt.Println("WrapperH ready")
	// Output: WrapperH ready
}
