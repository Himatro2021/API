package util

import (
	"fmt"
	"os"
	"time"
)

var errLogName = os.Getenv("ERR_LOG_FILE_NAME")
var errLogFile, _ = os.OpenFile(errLogName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

func LogErr(level, msg string) {
	errLogFile.WriteString(fmt.Sprintf("[%s]-[%s]: %s\n", level, time.Now().String(), msg))
}
