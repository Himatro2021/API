package util

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

var errLogName = os.Getenv("ERR_LOG_FILE_NAME")
var errLogFile, err = os.OpenFile(errLogName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

func LogErr(level, clue, msg string) {
	if err != nil {
		log.Print(err.Error())
	}
	errLogFile.WriteString(fmt.Sprintf("[%s]-[%s]-[%s]: %s\n", level, clue, time.Now().String(), msg))
}
