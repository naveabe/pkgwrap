package builder

import (
	"encoding/json"
	"github.com/naveabe/pkgwrap/pkgwrap/config"
	"github.com/naveabe/pkgwrap/pkgwrap/initscript"
	"github.com/naveabe/pkgwrap/pkgwrap/specer"
	"testing"
)

var (
	testCfgFile    = "/Users/abs/workbench/GoLang/src/github.com/naveabe/pkgwrap/etc/pkgwrap/pkgwrapd.conf.json.sample"
	testDockerUri  = "tcp://localhost:5555"
	testPkgName    = "annolityx"
	testPkgVersion = "0.0.1"
	testPkgPath    = "annolityx/0.0.1/annolityx.tgz"
	testRepoDir    = "/Users/abs/workbench/GoLang/src/github.com/naveabe/pkgwrap/data/repository"
	testContRun    ContainerRunner
	testRunnable   = initscript.BasicRunnable{
		Path: "/usr/local/bin/" + testPkgName,
		Args: "-l info", Env: map[string]string{},
	}
	testDistro, _  = specer.NewDistribution(specer.DISTRO_CENTOS, "")
	testUserPkg, _ = specer.NewUserPackage(testPkgName, testPkgVersion, testPkgPath, testRunnable)
	testConfig, _  = config.LoadConfigFromFile(testCfgFile)
)

func Test_ContainerRunner(t *testing.T) {

	testContRun = ContainerRunner{
		Distro:  testDistro,
		Package: testUserPkg,
	}
	testContRun.ContainerConfig = testContRun.initContainerConfig()

	jStr, _ := json.MarshalIndent(testContRun.ContainerConfig, "", "  ")
	t.Logf("%s", jStr)
}

func Test_NewContainerRunner(t *testing.T) {
	testUserPkg.BuildEnv = "go"
	testUserPkg.Packager = "metrilyx"

	testDistro2, _ := specer.NewDistribution(specer.DISTRO_CENTOS, "6")
	testDistro2.UserBuildCmd = []string{"make install"}
	testDistro2.BuildDeps = []string{"gcc-c++", "gcc"}

	cntr, err := NewContainerRunner(testConfig.Builder, testDistro2, testUserPkg)
	if err != nil {
		t.Fatalf("%s", err)
	}

	if _, err = cntr.Start(testDockerUri); err != nil {
		t.Fatalf("%s", err)
	}
}
