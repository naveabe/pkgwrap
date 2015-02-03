package specer

import (
	"github.com/naveabe/pkgwrap/pkgwrap/initscript"
	"testing"
)

var (
	testPkgVersion = "0.0.1"
	testTmpRepoDir = "/tmp"
	testPkgPath    = testPkgName + "/" + testPkgVersion + "/annolityx"
)

func Test_NewRPMSpec(t *testing.T) {
	spec, err := NewRPMSpec(testPkgName, testPkgVersion)
	if err != nil {
		t.Fatalf("%s", err)
	}
	//spec.Spec(os.Stdout)
}

func Test_BuildRPMSpec(t *testing.T) {
	runnable := initscript.BasicRunnable{"/path/to/bin", "-l info", map[string]string{}}
	pkg, _ := NewUserPackage(testPkgName, testPkgVersion, testPkgPath, runnable)

	tDistro, _ := NewDistribution("centos", "6")
	tDistro.Deps = []string{"zeromq3"}
	_, err := BuildRPMSpec(pkg, tDistro, testTmpRepoDir)
	if err != nil {
		t.Fatalf("%s", err)
	}
}
