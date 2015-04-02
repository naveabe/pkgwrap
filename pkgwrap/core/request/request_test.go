package request

import (
	"testing"
)

var (
	testPkgName = "annolityx"
	testYmlConf = "/Users/abs/workbench/GoLang/src/github.com/naveabe/pkgwrap/test/pkgwrap.yml"
)

func Test_NewPackageRequest(t *testing.T) {
	pkgreq := NewPackageRequest(testPkgName)
	t.Logf("%#v", pkgreq)
	//if err := pkgreq.Validate(); err != nil {
	//	t.Fatalf("%s", err)
	//}
}

func Test_NewPackageRequestFromYamlConfig(t *testing.T) {

	pkgreq, err := NewPackageRequestFromYamlConfig(testYmlConf)
	if err != nil {
		t.Fatalf("%s", err)
	}
	t.Logf("%v", pkgreq)
}

func Test_NewPackageRequestFromYamlConfig_Error(t *testing.T) {

	_, err := NewPackageRequestFromYamlConfig("")
	if err == nil {
		t.Fatalf("file not found - failed")
	}
}
