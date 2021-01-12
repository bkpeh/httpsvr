package main

import (
	"net/http"

	logging "github.com/bkpeh/httpsvr/util"
	hsvr "github.com/bkpeh/httpsvr/web"
)

func main() {
	hsvr.SetLog("logs/http.log")
	http.HandleFunc("/", hsvr.Index)

	//defer logging.GetLogFile().Close()

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		logging.LogError("logs/http.log", "Main:"+err.Error())
	}
}
