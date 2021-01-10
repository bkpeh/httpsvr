package main

import (
	"net/http"

	logging "github.com/bkpeh/httpsvr/util"
	hsvr "github.com/bkpeh/httpsvr/web"
)

func main() {

	http.HandleFunc("/", hsvr.Index)

	//defer logging.GetLogFile().Close()

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		logging.LogError("Main:" + err.Error())
	}
}
