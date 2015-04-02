package specer

import (
	"github.com/naveabe/pkgwrap/pkgwrap/core/initscript"
	"github.com/naveabe/pkgwrap/pkgwrap/core/request"
	"testing"
)

var (
	testPkgVersion = "0.0.1"
	testTmpRepoDir = "/tmp"
	testPkgPath    = testPkgName + "/" + testPkgVersion + "/annolityx"
	testBsRunnable = initscript.BasicRunnable{"/path/to/bin", "-l info", map[string]string{}}
)

func Test_NewRPMSpec(t *testing.T) {
	_, err := NewRPMSpec(testPkgName, testPkgVersion)
	if err != nil {
		t.Fatalf("%s", err)
	}
	//spec.Spec(os.Stdout)
}

func Test_BuildRPMSpec(t *testing.T) {

	pkg, _ := request.NewUserPackage(testPkgName, testPkgVersion, testPkgPath, testBsRunnable)

	tDistro, _ := request.NewDistribution("centos", "6")
	tDistro.Deps = []string{"zeromq3"}
	_, err := BuildRPMSpec(&testTmplMgr, pkg, tDistro, testTmpRepoDir)
	if err != nil {
		t.Fatalf("%s", err)
	}
}
