package logging

import (
	"log"
	"os"
)

//For logging
var logfile *os.File
var infolog *log.Logger
var errlog *log.Logger

func init() {
	logfile, err := os.OpenFile("logs/http.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Println(err)
	}

	infolog = log.New(logfile, "INFO: ", log.LstdFlags)
	errlog = log.New(logfile, "ERROR: ", log.LstdFlags)

	infolog.Println("init:Start Log")
}

func LogInfo(i interface{}) {
	infolog.Println(i)
}

func LogError(i interface{}) {
	errlog.Println(i)
}

func GetLogFile() *os.File {
	return logfile
}
