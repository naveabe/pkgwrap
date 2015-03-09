package websvc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/naveabe/pkgwrap/pkgwrap/config"
	"github.com/naveabe/pkgwrap/pkgwrap/initscript"
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
	"github.com/naveabe/pkgwrap/pkgwrap/repository"
	"github.com/naveabe/pkgwrap/pkgwrap/specer"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var (
	testPkgName             = "annolityx"
	testConfigFile          = "/Users/abs/workbench/GoLang/src/github.com/naveabe/pkgwrap/pkgwrapd.conf.json"
	testConfig, _           = config.LoadConfigFromFile(testConfigFile)
	testPkgBldMethodHandler = PkgBuilderMethodHandler{
		Config:     testConfig,
		Repository: repository.BuildRepository{testConfig.Repository},
		Logger:     logging.NewStdLogger(),
	}
	testPkgReq       = specer.NewPackageRequest(testPkgName)
	testRunnable     = initscript.BasicRunnable{}
	testDistro, _    = specer.NewDistribution("centos", "")
	testUploadParams = url.Values{
		"distros":        []string{"centos,centos-6"},
		"requires":       []string{"zeromq3"},
		"build_requires": []string{"zeromq3-devel"},
	}
	testBinPkgPath = "/Users/abs/workbench/GoLang/src/github.com/naveabe/pkgwrap/test/annolityx.tgz"
)

func newfileUploadRequest(uri string, params map[string]string, fileName, filePath string) error {
	var (
		body   bytes.Buffer
		writer       = multipart.NewWriter(&body)
		err    error = nil
	)

	if fileName != "" && filePath != "" {

		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		part, err := writer.CreateFormFile(fileName, filepath.Base(filePath))
		if err != nil {
			return err
		}
		if _, err = io.Copy(part, file); err != nil {
			return err
		}
	}

	for key, val := range params {

		if err = writer.WriteField(key, val); err != nil {
			return err
		}
	}
	// Write trailing boundary
	if err = writer.Close(); err != nil {
		return err
	}

	//fmt.Printf("%s\n", body)

	request, err := http.NewRequest("POST", uri, &body)
	if err != nil {
		return err
	}

	//request.Header.Set("Content-Type", "multipart/form")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	} else {
		body := &bytes.Buffer{}
		if _, err = body.ReadFrom(resp.Body); err != nil {
			return err
		}
		resp.Body.Close()

		fmt.Printf("%d %v\n", resp.StatusCode, resp.Header)
		fmt.Println(body)

		return nil
	}
}

func Test_PkgBuilderMethodHandler_JSONBody(t *testing.T) {
	go func() {
		NewPkgBuilderHandler(testConfig.APIPrefix, &testPkgBldMethodHandler, nil)
		t.Logf("Starting web server... :%d", testConfig.Port)
		http.ListenAndServe(fmt.Sprintf(":%d", testConfig.Port), nil)
	}()
	time.Sleep(2)

	//testPkgReq.Package, _ = specer.NewUserPackage(testPkgName, "0.0.1", "a/b/c", testRunnable)

	httpUrl := fmt.Sprintf("http://localhost:%d%s/%s/0.0.1",
		testConfig.Port, testConfig.APIPrefix, testPkgName)
	t.Logf("%s", httpUrl)

	t.Logf("%#v", testPkgReq)

	// Distribution error check
	body, _ := json.Marshal(testPkgReq)
	resp, err := http.Post(httpUrl, "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("%s", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode <= 304 {
		t.Fatalf("Mismatch: %d", resp.StatusCode)
	}

	// URL error check
	testPkgReq.Distributions = []specer.Distribution{testDistro}
	//testPkgReq.Package.BuildEnv = "go"
	body, _ = json.Marshal(testPkgReq)

	resp, err = http.Post(httpUrl, "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("%s", err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 304 {
		t.Fatalf("Mismatch: %d", resp.StatusCode)
	}

	// Packager error check
	testPkgReq.Package, _ = specer.NewUserPackage(testPkgName, "0.0.1",
		"annolityx/0.0.1/annolityx", testRunnable)
	testPkgReq.Package.URL = "https://github.com/metrilyx/annolityx"
	testPkgReq.Package.Packager = ""

	body, _ = json.Marshal(testPkgReq)
	resp, err = http.Post(httpUrl, "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("%s", err)
	}
	if resp.StatusCode != 400 {
		t.Fatalf("Mismatch: %d", resp.StatusCode)
	}

	testPkgReq.Package.Packager = "metrilyx"
	testPkgReq.Package.BuildEnv = "go"

	body, _ = json.Marshal(testPkgReq)
	resp, err = http.Post(httpUrl, "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("%s", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 304 {
		b, _ := ioutil.ReadAll(resp.Body)
		t.Errorf("%s", b)
		t.Fatalf("Status code: %d", resp.StatusCode)
	}

}

/*
func Test_PkgBuilderMethodHandler_Form(t *testing.T) {
	go func() {
		NewPkgBuilderHandler(testConfig.APIPrefix, &testPkgBldMethodHandler, nil)
		t.Logf("Starting web server... :%d", testConfig.Port)
		http.ListenAndServe(fmt.Sprintf(":%d", testConfig.Port), nil)
	}()
	time.Sleep(2)

	httpUrl := fmt.Sprintf("http://localhost:%d%s/%s/0.0.1",
		testConfig.Port, testConfig.APIPrefix, testPkgName)
	t.Logf("%s", httpUrl)

	//if err := newfileUploadRequest(httpUrl, testUploadParams, "package", testBinPkgPath); err != nil {
	//	t.Fatalf("%s", err)
	//}
	rsp, err := http.PostForm(httpUrl, testUploadParams)
	if err != nil {
		t.Fatalf("%s", err)
	}
	t.Logf("%v", rsp)
}
*/
