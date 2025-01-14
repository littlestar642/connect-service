package logger

import (
	"log"
	"os"
)

type multiLogger struct {
	fileLogger *log.Logger
}

var loggerInstance *multiLogger

func Init() error {
	logFile, err := os.OpenFile("./logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	loggerInstance = &multiLogger{
		fileLogger: log.New(logFile, "", log.Ltime),
	}

	return nil
}

func PrintToFile(v ...interface{}) {
	loggerInstance.fileLogger.Println(v...)
}
