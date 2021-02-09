package httpserver

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	logging "github.com/bkpeh/httpsvr/util"
)

var update = flag.Bool("update", false, "update .golden files")

type test struct {
	testname string
	method   string
	header   string
	body     string
	httpcode int
}

func TestReadjson(t *testing.T) {
	expect := map[string]emp{
		"1005": {"Jean", "Toh", 1005, "D5"},
		"1010": {"Peter", "Tan", 1010, "D1"},
		"1015": {"Jess", "Lim", 1015, "D4"},
		"1020": {"Mary", "Ong", 1020, "D2"},
		"1030": {"Jack", "Koh", 1030, "D3"},
		"1070": {"xxx", "yyy", 1070, "D1"},
		"1080": {"aaa", "bbb", 1080, "D1"},
	}

	result := readjson()

	equal := reflect.DeepEqual(expect, result)

	if !equal {
		t.Error("Readjson data mismatch")
	}
}

func TestReadinfo(t *testing.T) {
	//T1: To test with 1 item
	expect := map[string]emp{
		"1010": {"Peter", "Tan", 1010, "D1"},
	}

	input := url.Values{
		"Id": {"1010"},
	}
	SetLog("logs/test.log")
	result := readinfo(input)
	equal := reflect.DeepEqual(expect, result)

	if !equal {
		t.Error("T1:Readinfo data mismatch")
	}

	//T2: To test with full set of data
	expect = map[string]emp{
		"1005": {"Jean", "Toh", 1005, "D5"},
		"1010": {"Peter", "Tan", 1010, "D1"},
		"1015": {"Jess", "Lim", 1015, "D4"},
		"1020": {"Mary", "Ong", 1020, "D2"},
		"1030": {"Jack", "Koh", 1030, "D3"},
		"1070": {"xxx", "yyy", 1070, "D1"},
		"1080": {"aaa", "bbb", 1080, "D1"},
	}

	result = readinfo(url.Values{})

	equal = reflect.DeepEqual(expect, result)

	if !equal {
		t.Error("T2:Readinfo data mismatch")
	}
}

func TestDeleteinfo(t *testing.T) {
	//T1: To test delete 1 item (byte input case)
	expect := http.StatusNoContent
	data := map[string]emp{
		"1070": {"xxx", "yyy", 1070, "D1"},
	}

	input, _ := json.MarshalIndent(data, "", " ")
	result := deleteinfo(input)

	if expect != result {
		t.Errorf("T1:Deleteinfo unexpected result. Expect:%d, Result:%d", expect, result)
	}

	//T2: To test delete non existing item (byte input case)
	expect = http.StatusNotFound
	data = map[string]emp{
		"1075": {"xxx", "yyy", 1070, "D1"},
	}

	input, _ = json.MarshalIndent(data, "", " ")
	result = deleteinfo(input)

	if expect != result {
		t.Errorf("T2:Deleteinfo unexpected result. Expect:%d, Result:%d", expect, result)
	}

	//T3: To test delete 1 item (url.Value input case)
	expect = http.StatusNoContent
	dataT3 := url.Values{
		"Id": {"1080"},
	}

	result = deleteinfo(dataT3)

	if expect != result {
		t.Errorf("T3:Deleteinfo unexpected result. Expect:%d, Result:%d", expect, result)
	}

	//T4: To test delete non existing item (url.Value input case)
	expect = http.StatusNotFound
	dataT4 := url.Values{
		"Id": {"1085"},
	}

	result = deleteinfo(dataT4)

	if expect != result {
		t.Errorf("T4:Deleteinfo unexpected result. Expect:%d, Result:%d", expect, result)
	}
}

func TestUpdateinfo(t *testing.T) {
	//T1: To test create item
	expect := http.StatusCreated
	data := map[string]emp{
		"1070": {"xxx", "yyy", 1070, "D2"},
		"1080": {"aaa", "bbb", 1080, "D2"},
	}

	input, _ := json.MarshalIndent(data, "", " ")

	result := updateinfo(input)

	if expect != result {
		t.Errorf("T1:Updateinfo unexpected result. Expect:%d, Result:%d", expect, result)
	}

	//T2: To test update item
	expect = http.StatusOK
	data = map[string]emp{
		"1070": {"xxx", "yyy", 1070, "D1"},
		"1080": {"aaa", "bbb", 1080, "D1"},
	}

	input, _ = json.MarshalIndent(data, "", " ")

	result = updateinfo(input)

	if expect != result {
		t.Errorf("T2:Updateinfo unexpected result. Expect:%d, Result:%d", expect, result)
	}
}

func TestIndex(t *testing.T) {
	testtable := []test{
		{testname: "test1", method: "GET", header: "", body: "", httpcode: http.StatusOK},
		{testname: "test2", method: "GET", header: "", body: "../test/t2bodyinput.json", httpcode: http.StatusOK},
		{testname: "test3", method: "POST", header: "", body: "../test/t3bodyinput.json", httpcode: http.StatusOK},
		{testname: "test4", method: "POST", header: "../test/t4headinput.json", body: "", httpcode: http.StatusOK},
		//{testname: "PUT"},
		//{testname: "DELETE"},
	}

	for _, v := range testtable {
		t.Run(v.testname, func(t *testing.T) {
			gp := filepath.Join("../test", v.testname+".golden")

			req := httptest.NewRequest(v.method, "http://localhost:8080", nil)
			rec := httptest.NewRecorder()

			//For request header
			if v.header != "" {
				hinput := struct {
					Ids   url.Values
					Auth  string
					Ctype string
				}{}

				expect, _ := ioutil.ReadFile(v.header)

				err := json.Unmarshal(expect, &hinput)

				if err != nil {
					fmt.Println("Marshal Error", err)
				}

				req = httptest.NewRequest(v.method, "http://localhost:8080/", strings.NewReader(hinput.Ids.Encode()))
				req.Header.Add("Authorization", hinput.Auth)
				req.Header.Add("Content-Type", hinput.Ctype)
				req.Header.Add("Content-Length", strconv.Itoa(len(hinput.Ids)))
			}

			//For request body
			if v.body != "" {
				param := url.Values{}
				expect, _ := ioutil.ReadFile(v.body)

				json.Unmarshal(expect, &param)
				req.URL.RawQuery = param.Encode()
			}

			logging.LogInfo("logs/test.log", req)

			SetLog("logs/test.log")
			Index(rec, req)
			result, _ := ioutil.ReadAll(rec.Body)

			//Update Golden File
			if *update {
				ioutil.WriteFile(gp, result, 0644)
			}

			//Read JSON and replace /r (Carriage Return)
			expect, _ := ioutil.ReadFile(gp)
			expect = bytes.ReplaceAll(expect, []byte("\r"), []byte(""))

			if !bytes.Equal(result, expect) {
				logging.LogError("logs/test.log", "TestIndex:"+v.testname+":Result not expected")
				t.Fail()
			}

			if rec.Code != v.httpcode {
				logging.LogError("logs/test.log", "TestIndex:"+v.testname+":Http code mismatch")
				logging.LogError("logs/test.log", "Expect:"+strconv.Itoa(v.httpcode)+", Result:"+strconv.Itoa(rec.Code))
				t.Fail()
			}
		})
	}
}
