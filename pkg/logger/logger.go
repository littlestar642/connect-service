package logger

import (
	"log"
	"os"
)

type multiLogger struct {
	fileLogger *log.Logger
	stdOutLogger *log.Logger
}

var loggerInstance *multiLogger

func Init() error {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	loggerInstance = &multiLogger{
		fileLogger: log.New(logFile, "", log.Ltime),
		stdOutLogger: log.New(os.Stdout, "counter-service: ", log.LstdFlags),
	}
	
	return nil
}

func StdOut() *log.Logger {
	return loggerInstance.stdOutLogger
}

func PrintToFile(v ...interface{}) {
	loggerInstance.fileLogger.Println(v...)
}

