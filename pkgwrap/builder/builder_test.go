package builder

import (
	"github.com/naveabe/pkgwrap/pkgwrap/repository"
	"github.com/naveabe/pkgwrap/pkgwrap/specer"
	"github.com/naveabe/pkgwrap/pkgwrap/templater"
	"testing"
)

var (
	testDistro2, _ = specer.NewDistribution(specer.DISTRO_CENTOS, "6")
	testPkgReq     = specer.PackageRequest{
		Package:       testUserPkg,
		Distributions: []specer.Distribution{testDistro, testDistro2},
	}
	testBuildRepo = repository.BuildRepository{testRepoDir}
	testTmplMgr   = templater.TemplatesManager{testConfig.TemplatesDir()}
)

func Test_NewTargetedPackageBuild_Source(t *testing.T) {
	testPkgReq.Package.BuildType = specer.BUILDTYPE_SOURCE
	testPkgReq.Package.Packager = "metrilyx"

	b, err := NewTargetedPackageBuild(testConfig.Builder, testBuildRepo, &testPkgReq)
	if err != nil {
		t.Fatalf("%s", err)
	}
	t.Logf("%v", b.ListContainers())

	if err = b.SetupEnv(testTmplMgr); err != nil {
		t.Fatalf("Setup: %s", err)
	}

	cIds := b.StartBuilds(testDockerUri)
	if len(cIds) != 2 {
		t.Fatalf("Container id mistmatch: %d", len(cIds))
	}
}

func Test_NewTargetedPackageBuild_Bin_Error(t *testing.T) {
	testPkgReq.Package.BuildType = specer.BUILDTYPE_BIN

	testPkgReq.Package.Path += ".zip"
	b, _ := NewTargetedPackageBuild(testConfig.Builder, testBuildRepo, &testPkgReq)
	if err := b.SetupEnv(testTmplMgr); err == nil {
		t.Fatalf("Mismatch .zip")
	}

	testPkgReq.Package.Path += ".tbz2"
	b, _ = NewTargetedPackageBuild(testConfig.Builder, testBuildRepo, &testPkgReq)
	if err := b.SetupEnv(); err == nil {
		t.Fatalf("Mismatch bzip")
	}
}
