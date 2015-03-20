package templater

import (
	"os"
	"testing"
)

var (
	err          error
	testTmplName = "startup.sh"
	testTmplFile = "/Users/abs/workbench/GoLang/src/github.com/naveabe/pkgwrap/data/templates/startup/startup.sh"
)

func Test_NewTemplateBuilder_Error(t *testing.T) {
	if _, err = NewTemplateBuilder(testTmplName, "mab"); err == nil {
		t.Fatalf("mismatch")
	}
}

func Test_NewTemplateBuilder(t *testing.T) {

	tb, err := NewTemplateBuilder(testTmplName, testTmplFile)
	if err != nil {
		t.Fatalf("%s", err)
	}

	if err := tb.Build(struct {
		Name         string
		RunnablePath string
		RunnableArgs string
		Logfile      string
	}{"testname", "testRunPath", "testRunArgs", "testLogfile"}, os.Stdout); err != nil {
		t.Fatalf("%s", err)
	}
}
