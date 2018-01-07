package margo

import (
	"log"
	"fmt"
)

func logInfo(format string, args ... interface{}) {
	log.Printf("[margo] %s\n", fmt.Sprintf(format, args...))
}
