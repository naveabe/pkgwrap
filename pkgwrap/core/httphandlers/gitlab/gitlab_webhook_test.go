package gitlab

import (
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
	"net/http"
	"os"
	"testing"
	"time"
)

var (
	testGlEventFile = "/Users/abs/workbench/GoLang/src/github.com/naveabe/pkgwrap/test/gitlab_tag_event.json"
	testGlEvent     = GitlabTagEvent{
		Ref: "ref/tags/v0.0.1-dev",
	}
)

func Test_Version(t *testing.T) {

	version, err := testGlEvent.CheckVersion()
	if err != nil {
		t.Fatalf("%s", err)
	}
	if version != "0.0.1" {
		t.Fatalf("Version mismatch: %s", version)
	}

	testGlEvent.Ref = "ref/tags/v0.0.1"
	version, err = testGlEvent.CheckVersion()
	if err != nil {
		t.Fatalf("%s", err)
	}
	if version != "0.0.1" {
		t.Fatalf("Version mismatch: %s", version)
	}
}

func Test_Version_Error(t *testing.T) {
	testGlEvent := GitlabTagEvent{
		Ref: "ref/tags/v-dev",
	}
	if _, err := testGlEvent.CheckVersion(); err == nil {
		t.Fatalf("mismatch: %s", err)
	}
}

func Test_GitlabWebHook(t *testing.T) {
	go func() {
		glHandler := GitlabWebHook{logging.NewStdLogger()}
		http.ListenAndServe(":7654", &glHandler)
	}()
	time.Sleep(2)

	fh, err := os.Open(testGlEventFile)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer fh.Close()

	_, err = http.Post("http://localhost:7654/", "application/json", fh)
	if err != nil {
		t.Fatalf("%s", err)
	}
	//t.Logf("%v", resp)
}
