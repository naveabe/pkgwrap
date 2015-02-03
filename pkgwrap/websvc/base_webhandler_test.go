package websvc

import (
	"net/http"
	"testing"
	"time"
)

var (
	testPrefix        = "/api/buildertest"
	testAddr          = ":3434"
	testMethodHandler = DefaultMethodHandler{}
)

func Test_PkgBuilderHandler(t *testing.T) {
	go func() {
		NewPkgBuilderHandler(testPrefix, &testMethodHandler, nil)
		t.Logf("Starting web server... %s", testAddr)
		http.ListenAndServe(testAddr, nil)
	}()
	time.Sleep(2)
	// name test
	//httpResp, err := httpCall.Get("http://localhost:3434" + testPrefix + "/annolityx")
	resp, err := http.Get("http://localhost:3434" + testPrefix + "/annolityx")
	if err != nil {
		t.Fatalf("%s", err)
	}
	if resp.StatusCode != 405 {
		t.Fatalf("Status code: %d", resp.StatusCode)
	}

	resp, err = http.Post("http://localhost:3434"+testPrefix+"/annolityx", "text/plain", nil)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if resp.StatusCode != 405 {
		t.Fatalf("Status code: %d", resp.StatusCode)
	}

	// version test
	resp, err = http.Get("http://localhost:3434" + testPrefix + "/annolityx/0.0.1")
	if err != nil {
		t.Fatalf("%s", err)
	}
	if resp.StatusCode != 405 {
		t.Fatalf("Status code: %d", resp.StatusCode)
	}
	// error test
	resp, err = http.Get("http://localhost:3434" + testPrefix)
	if resp.StatusCode != 404 {
		t.Fatalf("Status code mismatch: %#v", resp)
	}
}

func Test_PkgBuilderHandler_TrailingSlash(t *testing.T) {
	go func() {
		NewPkgBuilderHandler("/with/prefix-slash/", &testMethodHandler, nil)
		http.ListenAndServe(testAddr, nil)
	}()
	// name test
	_, err := http.Get("http://localhost:3434/with/prefix-slash/annolityx")
	if err != nil {
		t.Fatalf("%s", err)
	}
}
