package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type emp struct {
	Fname string `json:"Fname"`
	Lname string `json:"Lname"`
	Id    int    `json:"Id"`
	Dept  string `json:"Dept"`
}

//For logging
var logfile *os.File
var infolog *log.Logger
var errlog *log.Logger

//Read from JSON file
func readjson() map[string]emp {
	rawlist := map[string]emp{}

	jfile, err := ioutil.ReadFile("json/list.json")

	if err != nil {
		errlog.Println("readjson:Error opening file.", err.Error())
	}

	err = json.Unmarshal(jfile, &rawlist)

	if err != nil {
		errlog.Println("deleteinfo:Error in Unmarshal.", err.Error())
	}

	return rawlist
}

//Create list to return for Respond
func readinfo(q url.Values) map[string]emp {
	rawlist := readjson()
	newlist := map[string]emp{}

	if _, ok := q["Id"]; ok && len(q) > 0 {
		for _, v := range q["Id"] {
			if _, ok := rawlist[v]; ok {
				newlist[v] = rawlist[v]
			}
		}
	} else {
		newlist = rawlist
	}

	infolog.Println("readinfo:Read", newlist)

	return newlist
}

func createinfo() {

}

//Delete data from JSON file
func deleteinfo(a interface{}) {
	plist := map[string]emp{}
	rawlist := readjson()

	switch a.(type) {
	case []byte:
		err := json.Unmarshal(a.([]byte), &plist)

		//Invalid JSON data
		if err != nil {
			errlog.Println("deleteinfo:Error in Marshal.", err.Error())
			return
		}

		for i, _ := range plist {
			delete(rawlist, i)
		}
	case url.Values:
		param := a.(url.Values)

		if len(param) == 0 {
			infolog.Println("deleteinfo:Empty url.Values.")
			return
		}

		if _, ok := param["Id"]; ok && len(param) > 0 {
			for _, v := range param["Id"] {
				infolog.Println("deleteinfo:Delete", rawlist[v])
				delete(rawlist, v)
			}
		}
	}

	wfile, err := os.OpenFile("json/list.json", os.O_WRONLY|os.O_TRUNC, 0777)

	if err != nil {
		errlog.Println("deleteinfo:Error in OpenFile.", err.Error())
	}

	defer wfile.Close()

	jfile, err := json.MarshalIndent(rawlist, "", "	")

	if err != nil {
		errlog.Println("deleteinfo:Error in Marshal.", err.Error())
	}

	if _, err = wfile.Write(jfile); err != nil {
		errlog.Println("deleteinfo:Error in writing to JSON file.", err.Error())
	}
}

//Insert data into JSON file
func updateinfo(by []byte) {
	plist := map[string]emp{}

	err := json.Unmarshal(by, &plist)

	//Invalid JSON data
	if err != nil {
		errlog.Println("updateinfo:Error in Unmarshal.", err.Error())
		return
	}

	rawlist := readjson()
	wfile, err := os.OpenFile("json/list.json", os.O_WRONLY|os.O_TRUNC, 0777)
	defer wfile.Close()

	for i, v := range plist {
		rawlist[i] = v
	}

	jfile, err := json.MarshalIndent(rawlist, "", "	")

	if err != nil {
		errlog.Println("updateinfo:Error in Marshal.", err.Error())
	}

	if _, err = wfile.Write(jfile); err != nil {
		errlog.Println("updateinfo:Error in writing to JSON file.", err.Error())
	}
}

func index(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":
		err := json.NewEncoder(w).Encode(readinfo(r.URL.Query()))

		if err != nil {
			errlog.Println("index:GET:Error encoding.", err.Error())
		}

	case "POST":
		if len(r.URL.Query()) > 0 || r.Header.Get("Content-type") == "application/x-www-form-urlencoded" {
			r.ParseForm()
			err := json.NewEncoder(w).Encode(readinfo(r.Form))

			if err != nil {
				errlog.Println("index:POST:Error encoding.", err.Error())
			}
		}

	case "PUT":
		by, err := ioutil.ReadAll(r.Body)

		if err != nil {
			errlog.Println("index:PUT:Error reading body.", err.Error())
		}

		updateinfo(by)

	case "DELETE":
		if len(r.URL.Query()) > 0 {
			deleteinfo(r.URL.Query())
		}

		if r.Header.Get("Content-type") == "application/json" {
			by, err := ioutil.ReadAll(r.Body)

			if err != nil {
				errlog.Println("index:PUT:Error reading body.", err.Error())
			}

			deleteinfo(by)
		}

	default:
		infolog.Println("index:Default")

	}
}

func init() {
	logfile, err := os.OpenFile("logs/http.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Println(err)
	}

	infolog = log.New(logfile, "INFO: ", log.LstdFlags)
	errlog = log.New(logfile, "ERROR: ", log.LstdFlags)

	infolog.Println("init:Start Log")
}

func main() {
	http.HandleFunc("/", index)

	//Defer closing of logfile initalise in init()
	defer logfile.Close()

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		errlog.Println("Main:", err.Error())
	}
}
