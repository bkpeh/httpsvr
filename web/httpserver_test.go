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

			}

			//For request body
			if v.body != "" {
				param := map[string][]string{}
				expect, _ := ioutil.ReadFile(v.body)

				json.Unmarshal(expect, &param)
				urlstr := url.Values{}

				for i, v := range param {
					for _, vv := range v {
						urlstr.Add(i, vv)
					}
				}

				req.URL.RawQuery = urlstr.Encode()
			}

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
				logging.LogError("../logs/http.log", "TestIndex:Result not expected")
				t.Fail()
			}

			if rec.Code != v.httpcode {
				logging.LogError("../logs/http.log", "TestIndex:Http code mismatch")
				t.Fail()
			}
		})
	}
}
