package main

import "fmt"
import "net/http"
import "encoding/json"
import "io/ioutil"
import "os"

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
		fmt.Println("Error opening file.")
	}

	err = json.Unmarshal(jfile, &rawlist)

	if err != nil {
		fmt.Println("Error to Unmarshal.")
	}

	return rawlist
}

//Create list to return
func readinfo(m map[string][]string) []emp {
	rawlist := readjson()
	plist := []emp{}

	if _, ok := m["Id"]; ok && len(m["Id"]) > 0 {
		for _, v := range m["Id"] {
			plist = append(plist, rawlist[v])
		}
	} else {
		for _, v := range rawlist {
			plist = append(plist, v)
		}
	}

	return plist
}

func createinfo() {

}

func deleteinfo() {

}

//Insert data into JSON file
func updateinfo(by []byte) {
	plist := map[string]emp{}
	wfile, err := os.OpenFile("json/list.json", os.O_WRONLY, 0777)
	defer wfile.Close()

	json.Unmarshal(by, &plist)
	rawlist := readjson()

	for i, v := range plist {
		rawlist[i] = v
	}

	jfile, err := json.MarshalIndent(rawlist, "", "	")

	if err != nil {
		fmt.Println("Error in Marshal.")
	}

	if _, err = wfile.Write(jfile); err != nil {
		fmt.Println("Error in writing to JSON file.")
	}
}

func index(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":
		err := json.NewEncoder(w).Encode(readinfo(r.URL.Query()))

		if err != nil {
			fmt.Fprintf(w, "Error encoding")
		}

	case "POST":

	case "PUT":
		by, err := ioutil.ReadAll(r.Body)

		if err != nil {
			fmt.Println("Error reading body")
		}

		updateinfo(by)

	case "DELETE":

	default:
		fmt.Fprintf(w, "Default")

	}
}

func init() {

}

func main() {

	http.HandleFunc("/", index)

	err := http.ListenAndServe(":8080", nil)

	fmt.Println(err)
}
