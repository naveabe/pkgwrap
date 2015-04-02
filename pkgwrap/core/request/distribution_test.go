package request

import (
	"testing"
)

var (
	testPkgr = "mock"
)

func Test_NewDistribution_Error(t *testing.T) {
	_, err := NewDistribution(DISTRO_CENTOS, "8")
	if err == nil {
		t.Fatalf("mismatch")
	}
}

func Test_NewDistribution(t *testing.T) {
	d, err := NewDistribution(DISTRO_CENTOS, "")
	if err != nil {
		t.Fatalf("%s", err)
	}

	bdir := d.BuildDir(testPkgr)
	t.Logf("%s", bdir)

	t.Logf(d.BuildCommand())

	if d.PackageType() != OS_PKG_TYPE_RPM {
		t.Fatalf("mismatch")
	}
}
