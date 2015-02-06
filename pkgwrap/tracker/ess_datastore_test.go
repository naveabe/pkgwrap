package tracker

import (
	"github.com/naveabe/pkgwrap/pkgwrap/config"
	"testing"
)

var (
	testCfgFile     = "/Users/abs/workbench/GoLang/src/github.com/naveabe/pkgwrap/etc/pkgwrap/pkgwrapd.conf.json"
	testConfig, _   = config.LoadConfigFromFile(testCfgFile)
	testSearchTerms = map[string]string{"username": "metrilyx"}
	testEds         *EssDatastore
)

func Test_EssDatastore_GetBuildsForUser(t *testing.T) {
	var err error
	testEds, err = NewEssDatastore(&testConfig.JobTracker.Datastore, nil)
	if err != nil {
		t.Fatalf("%s", err)
	}

	rslt, err := testEds.GetBuildsForUser("metrilyx")
	if err != nil {
		t.Fatalf("%s", err)
	}
	if len(rslt) == 0 {
		t.Fatalf("No results")
	}
	t.Logf("%v", rslt)
}

func Test_EssDatastore_GetBuildsForPackage(t *testing.T) {
	rslt, err := testEds.GetBuildsForPackage("metrilyx", "annolityx")
	if err != nil {
		t.Fatalf("%s", err)
	}
	if len(rslt) == 0 {
		t.Fatalf("No results")
	}
	t.Logf("%v", rslt)
}

func Test_EssDatastore_GetBuildsForPackageVersion(t *testing.T) {
	rslt, err := testEds.GetBuildsForPackageVersion("metrilyx", "annolityx", "0.0.1")
	if err != nil {
		t.Fatalf("%s", err)
	}
	if len(rslt) == 0 {
		t.Fatalf("No results")
	}
	t.Logf("%v", rslt)
}
