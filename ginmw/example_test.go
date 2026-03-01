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
