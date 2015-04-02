package specer

import (
	"github.com/naveabe/pkgwrap/pkgwrap/core/request"
	"github.com/naveabe/pkgwrap/pkgwrap/templater"
	"testing"
)

var (
	testPkgName = "annolityx"
	testTmplMgr = templater.TemplatesManager{"/Users/abs/workbench/GoLang/src/github.com/naveabe/pkgwrap/data/templates"}
	testUrl     = "https://github.com/metrilyx/annolityx"
)

func Test_DEBSpec_WriteDirStructure(t *testing.T) {
	debSpec := DEBSpec{
		PackageMetadata: PackageMetadata{
			Name:    testPkgName,
			Version: testPkgVersion,
			Release: 1,
		},
		BuildDeps: "libzmq3-dev",
		Deps:      "libzmq3",
	}

	if err := debSpec.WriteDirStructure("/tmp"); err != nil {
		t.Fatalf("%s", err)
	}
}

func Test_BuildDEBSpec(t *testing.T) {
	distro, _ := request.NewDistribution("ubuntu", "14.04")
	pkg, _ := request.NewUserPackage(testPkgName, testPkgVersion, testPkgPath, testBsRunnable)
	pkg.URL = testUrl
	pkg.Packager = "metrilyx"
	if err := BuildDebStructure(&testTmplMgr, pkg, distro, "/tmp"); err != nil {
		t.Fatalf("%s", err)
	}
}
