package repository

import (
	"testing"
)

var (
	testPkgr        = "metrilyx"
	testRepoDir     = "/Users/abs/workbench/GoLang/src/github.com/naveabe/pkgwrap/data/repository"
	testPkgName     = "annolityx"
	testPkgVersion  = "0.0.1"
	testRepo        = BuildRepository{testRepoDir}
	testDistroLabel = "centos-6"
)

func Test_LastRelease(t *testing.T) {
	t.Logf("%d", testRepo.LastRelease(testPkgr, testPkgName, testPkgVersion, testDistroLabel))
}
func Test_NextRelease(t *testing.T) {
	t.Logf("%d", testRepo.NextRelease(testPkgr, testPkgName, testPkgVersion, testDistroLabel))
}

func Test_ListPackages(t *testing.T) {
	//t.Logf("%d", testRepo.NextRelease(testPkgr, testPkgName, testPkgVersion, "centos-6"))
	list, err := testRepo.ListPackages(testPkgr, testPkgName, testPkgVersion, testDistroLabel)
	if err != nil {
		t.Fatalf("%s", err)
	}
	t.Logf("%v", list)
}
