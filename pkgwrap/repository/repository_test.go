package repository

import (
	"testing"
)

var (
	testRepoDir    = "/Users/abs/workbench/GoLang/src/github.com/naveabe/pkgwrap/data/repository"
	testPkgName    = "annolityx"
	testPkgVersion = "0.0.1"
	testRepo       = BuildRepository{testRepoDir}
)

func Test_LastRelease(t *testing.T) {
	t.Logf("%d", testRepo.LastRelease(testPkgName, testPkgVersion, "rpm"))
}
func Test_NextRelease(t *testing.T) {
	t.Logf("%d", testRepo.NextRelease(testPkgName, testPkgVersion, "rpm"))
}
