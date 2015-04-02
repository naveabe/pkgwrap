package initscript

import (
	//"os"
	"testing"
)

/*
var (
	testRepodir = "/tmp"
	testOutfile = "/tmp/pkgbuilder.test.out"
)
*/
func Test_NewBasicInitScript(t *testing.T) {
	bis, err := NewBasicInitScript("annolityx")
	if err != nil {
		t.Fatalf("%s", err)
	}
	bis.Runnable = BasicRunnable{"/bin/bash", "-c 'echo hi'", nil}
	/*
		if err = bis.StartupScript(os.Stdout); err != nil {
			t.Fatalf("%s", err)
		}

		_, err = bis.Write(testRepodir, "0.0.2")
		if err == nil {
			t.Fatalf("mismatch")
		}
		_, err = bis.Write(testRepodir, "0.0.1")
		if err != nil {
			t.Fatalf("%s", err)
		}
	*/
}
