package github

import (
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

var (
	testGhEventFile = "/Users/abs/workbench/GoLang/src/github.com/naveabe/pkgwrap/test/github_create_event.json"
)

func Test_GithubWebHook(t *testing.T) {
	go func() {
		glHandler := GithubWebHook{logging.NewStdLogger()}
		http.ListenAndServe(":7655", &glHandler)
	}()
	time.Sleep(2)

	fh, err := os.Open(testGhEventFile)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer fh.Close()

	resp, err := http.Post("http://localhost:7655/", "application/json", fh)
	if err != nil {
		t.Fatalf("%s", err)
	}
	b, _ := ioutil.ReadAll(resp.Body)
	t.Logf("%s", b)
}
