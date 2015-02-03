package specer

import (
	"testing"
)

func Test_WriteDirStructure(t *testing.T) {
	debSpec := DEBSpec{
		PackageMetadata: PackageMetadata{
			Name:    testPkgName,
			Version: testPkgVersion,
			Release: 1,
		},
		BuildDeps: "zmq3-dev",
		Deps:      "zmq3",
	}

	if err := debSpec.WriteDirStructure("/tmp"); err != nil {
		t.Fatalf("%s", err)
	}
}
