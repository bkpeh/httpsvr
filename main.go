package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	logging "github.com/httpsvr/util"
)

type emp struct {
	Fname string `json:"Fname"`
	Lname string `json:"Lname"`
	Id    int    `json:"Id"`
	Dept  string `json:"Dept"`
}

//Read from JSON file
func readjson() map[string]emp {
	rawlist := map[string]emp{}

	jfile, err := ioutil.ReadFile("json/list.json")

	if err != nil {
		logging.LogError("readjson:Error opening file:" + err.Error())
	}

	err = json.Unmarshal(jfile, &rawlist)

	if err != nil {
		logging.LogError("deleteinfo:Error in Unmarshal:" + err.Error())
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

	logging.LogInfo(newlist)

	return newlist
}

func createinfo() {

}

//Delete data from JSON file
func deleteinfo(a interface{}) int {
	var code int
	plist := map[string]emp{}
	rawlist := readjson()

	switch a.(type) {
	case []byte:
		err := json.Unmarshal(a.([]byte), &plist)

		//Invalid JSON data
		if err != nil {
			logging.LogError("deleteinfo:Error in Marshal:" + err.Error())
			code = http.StatusBadRequest
		}

		code = http.StatusNotFound
		for i, _ := range plist {
			_, ok := rawlist[i]

			if ok {
				delete(rawlist, i)
				code = http.StatusNoContent
			}
		}
	case url.Values:
		param := a.(url.Values)

		if len(param) == 0 {
			logging.LogInfo("deleteinfo:Empty url.Values.")
			code = http.StatusBadRequest
		}

		if _, ok := param["Id"]; ok && len(param) > 0 {
			code = http.StatusNotFound

			for _, v := range param["Id"] {
				_, ok := rawlist[v]

				if ok {
					logging.LogInfo("deleteinfo:Delete:" + v)
					delete(rawlist, v)
					code = http.StatusNoContent
				}
			}
		} else {
			code = http.StatusBadRequest
		}
	default:
		code = http.StatusBadRequest
	}

	if code == http.StatusBadRequest || code == http.StatusNotFound {
		return code
	}

	wfile, err := os.OpenFile("json/list.json", os.O_WRONLY|os.O_TRUNC, 0777)

	if err != nil {
		logging.LogError("deleteinfo:Error in OpenFile:" + err.Error())
		code = http.StatusInternalServerError
	}

	defer wfile.Close()

	jfile, err := json.MarshalIndent(rawlist, "", "	")

	if err != nil {
		logging.LogError("deleteinfo:Error in Marshal:" + err.Error())
		code = http.StatusInternalServerError
	}

	if _, err = wfile.Write(jfile); err != nil {
		logging.LogError("deleteinfo:Error in writing to JSON file:" + err.Error())
		code = http.StatusInternalServerError
	}

	return code
}

//Insert data into JSON file
func updateinfo(by []byte) int {
	var code int = http.StatusCreated
	plist := map[string]emp{}
	err := json.Unmarshal(by, &plist)

	//Invalid JSON data
	if err != nil {
		logging.LogError("updateinfo:Error in Unmarshal:" + err.Error())
		return http.StatusBadRequest
	}

	rawlist := readjson()
	wfile, err := os.OpenFile("json/list.json", os.O_WRONLY|os.O_TRUNC, 0777)
	defer wfile.Close()

	for i, v := range plist {
		_, ok := rawlist[i]

		if ok {
			//Update. Set to StatusOK
			code = http.StatusOK
		}
		rawlist[i] = v
	}

	jfile, err := json.MarshalIndent(rawlist, "", "	")

	if err != nil {
		logging.LogError("updateinfo:Error in Marshal:" + err.Error())
		code = http.StatusBadRequest
	}

	if _, err = wfile.Write(jfile); err != nil {
		logging.LogError("updateinfo:Error in writing to JSON file:" + err.Error())
		code = http.StatusInternalServerError
	}

	return code
}

func index(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":
		resp := readinfo(r.URL.Query())

		if len(resp) <= 0 {
			logging.LogInfo("index:GET:Not Found")
			w.WriteHeader(http.StatusNotFound)
		} else {
			err := json.NewEncoder(w).Encode(resp)

			if err != nil {
				logging.LogError("index:GET:Error encoding:" + err.Error())
				w.WriteHeader(http.StatusBadRequest)
			}
		}

	case "POST":
		if len(r.URL.Query()) > 0 || r.Header.Get("Content-type") == "application/x-www-form-urlencoded" {
			r.ParseForm()
			resp := readinfo(r.Form)

			if len(resp) <= 0 {
				logging.LogInfo("index:POST:Not Found")
				w.WriteHeader(http.StatusNotFound)
			} else {
				err := json.NewEncoder(w).Encode(resp)

				if err != nil {
					logging.LogError("index:POST:Error encoding:" + err.Error())
					w.WriteHeader(http.StatusBadRequest)
				}
			}
		} else if r.Header.Get("Content-type") == "application/json" {
			by, err := ioutil.ReadAll(r.Body)

			if err != nil {
				logging.LogError("index:POST:Error reading body:" + err.Error())
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(updateinfo(by))
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		logging.LogInfo("POST SUCCESS")
	case "PUT":
		by, err := ioutil.ReadAll(r.Body)

		if err != nil {
			logging.LogError("index:PUT:Error reading body:" + err.Error())
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(updateinfo(by))
		}

	case "DELETE":
		if len(r.URL.Query()) > 0 {
			w.WriteHeader(deleteinfo(r.URL.Query()))
		} else if r.Header.Get("Content-type") == "application/json" {
			by, err := ioutil.ReadAll(r.Body)

			if err != nil {
				logging.LogError("index:PUT:Error reading body:" + err.Error())
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(deleteinfo(by))
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

	default:
		logging.LogInfo("index:Default")

	}
}

func main() {
	http.HandleFunc("/", index)

	//defer logging.GetLogFile().Close()

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		logging.LogError("Main:" + err.Error())
	}
}
