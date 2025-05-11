package logger

import (
	"log"
	"os"
)

var logInstance = log.New(os.Stdout, "[LB] ", log.LstdFlags)

func Log() *log.Logger {
	return logInstance
}
