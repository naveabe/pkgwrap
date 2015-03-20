package specer

import (
	//"encoding/json"
	"github.com/naveabe/pkgwrap/pkgwrap/initscript"
	//"io/ioutil"
	"testing"
)

var (
	testTgz     = "/Users/abs/workbench/GoLang/src/github.com/naveabe/pkgwrap/test/annolityx.tgz"
	testDst     = "/tmp"
	testRepoDir = "/Users/abs/workbench/GoLang/src/github.com/naveabe/pkgwrap/data/repository"
	testBldFile = "/Users/abs/workbench/GoLang/src/github.com/naveabe/pkgwrap/test/pkgbuild.json"
)

func Test_UserPackage_Error(t *testing.T) {
	var (
		err error
		up  = UserPackage{
			PackageMetadata: PackageMetadata{
				Name:    testPkgName,
				Version: "0.0.1",
			},
			Path: "file/not/found",
		}
	)

	if err = up.Uncompress(""); err == nil {
		t.Fatalf("mismatch")
	}

	up = UserPackage{
		PackageMetadata: PackageMetadata{
			Name:    testPkgName,
			Version: "0.0.1",
		},
		Path: "annolityx/0.0.1/annolityx.tgz",
	}
	if err = up.Uncompress("/not/found"); err == nil {
		t.Fatalf("mismatch")
	}
}

func Test_UserPackage(t *testing.T) {
	runnable := initscript.BasicRunnable{"/usr/local/bin/anno", "-l info", map[string]string{}}
	up, err := NewUserPackage(testPkgName, testPkgVersion, "annolityx/0.0.1/annolityx.tgz", runnable)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if err = up.Uncompress(testRepoDir); err != nil {
		t.Fatalf("%s", err)
	}
}
