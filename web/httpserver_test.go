package httpserver_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	logging "github.com/bkpeh/httpsvr/util"
	h "github.com/bkpeh/httpsvr/web"
)

var update = flag.Bool("update", false, "update .golden files")

type test struct {
	testname string
	method   string
	header   string
	body     string
	httpcode int
}

func TestIndex(t *testing.T) {
	testtable := []test{
		{testname: "test1", method: "GET", header: "", body: "", httpcode: http.StatusOK},
		{testname: "test2", method: "GET", header: "", body: "../test/t2bodyinput.json", httpcode: http.StatusOK},
		{testname: "test3", method: "POST", header: "", body: "../test/t3bodyinput.json", httpcode: http.StatusOK},
		{testname: "test4", method: "POST", header: "../test/t4bodyinput.json", body: "../test/t4bodyinput.json", httpcode: http.StatusOK},
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
					id    url.Values
					auth  string
					ctype string
				}{}
				expect, _ := ioutil.ReadFile(v.header)

				json.Unmarshal(expect, &hinput)
				req = httptest.NewRequest("POST", "http://localhost:8080/", strings.NewReader(hinput.id.Encode()))
				req.Header.Add("Authorization", hinput.auth)
				req.Header.Add("Content-Type", hinput.ctype)
				req.Header.Add("Content-Length", strconv.Itoa(len(hinput.id)))
				//urlstr := url.Values{}

				//for i, v := range hinput {
				// if v.(Type) == []string {

				//}
				//urlstr.Add(i, v)
				//}

				//req.URL.RawQuery = urlstr.Encode()
			}

			//For request body
			if v.body != "" {
				param := url.Values{}
				expect, _ := ioutil.ReadFile(v.body)

				json.Unmarshal(expect, &param)
				req.URL.RawQuery = param.Encode()
			}

			h.SetLog("logs/test.log")
			h.Index(rec, req)
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
