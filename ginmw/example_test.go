package ginmw_test

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/peakxie/utils/ginmw"
)

// ExampleSetLogFunc shows how to set a global log function.
func ExampleSetLogFunc() {
	ginmw.SetLogFunc(func(c *gin.Context, format string, args ...any) {
		log.Printf(format, args...)
	})
	fmt.Println("log function set")
	ginmw.SetLogFunc(nil) // reset
	// Output: log function set
}

// ExampleSetErrorLogFunc shows how to set a separate error-level log function.
func ExampleSetErrorLogFunc() {
	ginmw.SetLogFunc(func(c *gin.Context, format string, args ...any) {
		log.Printf("[INFO] "+format, args...)
	})
	ginmw.SetErrorLogFunc(func(c *gin.Context, format string, args ...any) {
		log.Printf("[ERROR] "+format, args...)
	})
	fmt.Println("error log function set")
	ginmw.SetLogFunc(nil)      // reset
	ginmw.SetErrorLogFunc(nil) // reset
	// Output: error log function set
}
