package initscript

import (
	"testing"
)

var (
	testBin     = "/usr/local/bin/annolityx"
	testBinArgs = "-l info -c annolityx.toml"
	testEnvVars = map[string]string{"APPHOME": "/opt/app", "TMP": "/tmp"}
)

func Test_BasicRunnable(t *testing.T) {
	r := BasicRunnable{testBin, testBinArgs, testEnvVars}
	t.Logf("%s", r.Command())
}

func Test_BasicRunnable_NoEnv(t *testing.T) {
	r := BasicRunnable{testBin, testBinArgs, nil}
	t.Logf("%s", r.Command())
}
