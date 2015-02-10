package specer

import (
	"github.com/naveabe/pkgwrap/pkgwrap/templater"
	"testing"
)

var (
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
		BuildDeps: []string{"libzmq3-dev"},
		Deps:      []string{"libzmq3"},
	}

	if err := debSpec.WriteDirStructure("/tmp"); err != nil {
		t.Fatalf("%s", err)
	}
}

func Test_BuildDEBSpec(t *testing.T) {
	distro, _ := NewDistribution("ubuntu", "14.04")
	pkg, _ := NewUserPackage(testPkgName, testPkgVersion, testPkgPath, testBsRunnable)
	pkg.URL = testUrl
	pkg.Packager = "metrilyx"
	if err := BuildDebStructure(&testTmplMgr, pkg, distro, "/tmp"); err != nil {
		t.Fatalf("%s", err)
	}
}
