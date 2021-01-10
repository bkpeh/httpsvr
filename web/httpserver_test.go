package httpserver_test

import (
	"bytes"
	"flag"
	"io/ioutil"
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
	body     string
}

func TestIndex(t *testing.T) {
	testtable := []test{
		{testname: "TEST1", method: "GET", body: ""},
		{testname: "TEST2", method: "GET", body: "1020"},
		//{testname: "POST"},
		//{testname: "PUT"},
		//{testname: "DELETE"},
	}

	for _, v := range testtable {
		t.Run(v.testname, func(t *testing.T) {
			gp := filepath.Join("../test", v.testname+".golden")

			req := httptest.NewRequest(v.method, "http://localhost:8080", nil)

			if v.body != "" {
				req.URL.RawQuery = url.Values{"Id": {"1020"}}.Encode()
			}

			rec := httptest.NewRecorder()

			h.Index(rec, req)

			result, _ := ioutil.ReadAll(rec.Body)

			if *update {
				ioutil.WriteFile(gp, result, 0644)
			}

			expect, _ := ioutil.ReadFile(gp)

			if !bytes.Equal(result, expect) {
				logging.LogError("TestIndex:")
			}
		})
	}
}
