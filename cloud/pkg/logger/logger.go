package logger

import (
	"log"
	"os"
)

func Log() *log.Logger {
	return log.New(os.Stdout, "[LB] ", log.LstdFlags)
}
