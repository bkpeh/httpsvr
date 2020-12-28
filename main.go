package main

import "fmt"
import "net/http"
import "encoding/json"
import "io/ioutil"
import "os"
import "net/url"

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
		fmt.Println("readjson:Error opening file.", err)
	}

	err = json.Unmarshal(jfile, &rawlist)

	if err != nil {
		fmt.Println("readjson:Error to Unmarshal.", err)
	}

	return rawlist
}

//Create list to return for Respond
func readinfo(q url.Values) map[string]emp {
	rawlist := readjson()
	newlist := map[string]emp{}

	if _, ok := q["Id"]; ok && len(q) > 0 {
		for _, v := range q["Id"] {
			newlist[v] = rawlist[v]
		}
	} else {
		newlist = rawlist
	}

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
			fmt.Println("deleteinfo:Error in Marshal.", err)
			return
		}

		for i, _ := range plist {
			delete(rawlist, i)
		}
	case url.Values:
		param := a.(url.Values)

		if len(param) == 0 {
			fmt.Println("deleteinfo:Empty url.Values.")
			return
		}

		if _, ok := param["Id"]; ok && len(param) > 0 {
			for _, v := range param["Id"] {
				delete(rawlist, v)
			}
		}
	}

	wfile, err := os.OpenFile("json/list.json", os.O_WRONLY|os.O_TRUNC, 0777)

	if err != nil {
		fmt.Println("deleteinfo:Error in OpenFile.", err)
	}

	defer wfile.Close()

	jfile, err := json.MarshalIndent(rawlist, "", "	")

	if err != nil {
		fmt.Println("deleteinfo:Error in Marshal.", err)
	}

	if _, err = wfile.Write(jfile); err != nil {
		fmt.Println("deleteinfo:Error in writing to JSON file.", err)
	}
}

//Insert data into JSON file
func updateinfo(by []byte) {
	plist := map[string]emp{}

	err := json.Unmarshal(by, &plist)

	//Invalid JSON data
	if err != nil {
		fmt.Println("updateinfo:Error in Marshal.", err)
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
		fmt.Println("updateinfo:Error in Marshal.", err)
	}

	if _, err = wfile.Write(jfile); err != nil {
		fmt.Println("updateinfo:Error in writing to JSON file.", err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":
		err := json.NewEncoder(w).Encode(readinfo(r.URL.Query()))

		if err != nil {
			fmt.Println("index:GET:Error encoding", err)
		}

	case "POST":

	case "PUT":
		by, err := ioutil.ReadAll(r.Body)

		if err != nil {
			fmt.Println("index:PUT:Error reading body", err)
		}

		updateinfo(by)

	case "DELETE":
		if len(r.URL.Query()) > 0 {
			deleteinfo(r.URL.Query())
		}

		if r.Header.Get("Content-type") == "application/json" {
			by, err := ioutil.ReadAll(r.Body)

			if err != nil {
				fmt.Println("index:PUT:Error reading body", err)
			}

			deleteinfo(by)
		}

	default:
		fmt.Println("Default")

	}
}

func init() {

}

func main() {

	http.HandleFunc("/", index)

	err := http.ListenAndServe(":8080", nil)

	fmt.Println(err)
}
