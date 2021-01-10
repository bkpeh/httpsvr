package logging

import (
	"log"
	"os"
)

//For logging
var logfile *os.File
var infolog *log.Logger
var errlog *log.Logger

func openFile() {
	logfile, err := os.OpenFile("../logs/http.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Println(err)
	}

	infolog = log.New(logfile, "INFO: ", log.LstdFlags)
	errlog = log.New(logfile, "ERROR: ", log.LstdFlags)

	infolog.Println("init:Start Log")
}

func LogInfo(i interface{}) {
	openFile()
	infolog.Println(i)
	logfile.Sync()
	logfile.Close()
}

func LogError(i interface{}) {
	openFile()
	errlog.Println(i)
	logfile.Sync()
	logfile.Close()
}

func GetLogFile() *os.File {
	return logfile
}
