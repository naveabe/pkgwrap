package config

import (
	"encoding/json"
	//"io/ioutil"
	"testing"
)

var (
	testCfgFile      = "/Users/abs/workbench/GoLang/src/github.com/naveabe/pkgwrap/etc/pkgwrap/pkgwrapd.conf.json.sample"
	testCfgFileError = "/Users/abs/workbench/GoLang/src/github.com/naveabe/pkgwrap/test/pkgwrapd.conf.sample.error"
)

func Test_LoadConfigFromFile(t *testing.T) {
	cfg, err := LoadConfigFromFile(testCfgFile)
	if err != nil {
		t.Fatalf("%s", err)
	}
	b, _ := json.MarshalIndent(cfg, "", "  ")
	t.Logf("%s\n", b)
}

func Test_LoadConfigFromFile_Error(t *testing.T) {
	_, err := LoadConfigFromFile("/not/found")
	if err == nil {
		t.Fatalf("mismatch")
	}
	_, err = LoadConfigFromFile(testCfgFileError)
	if err == nil {
		t.Fatalf("Bad json check failed")
	}
}
