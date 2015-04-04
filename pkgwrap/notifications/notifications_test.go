package notifications

import (
	"testing"
)

var (
	testEmail = "user@domain.io"
)

func Test_GetNotificationMessage(t *testing.T) {
	msg := GetNotificationMessage("pkgname", "0.0.1", "Succeeded", "ubuntu", "14.04")
	t.Logf("%s", msg)
}

func Test_EmailNotifier(t *testing.T) {
	en := EmailNotifier{
		"euforia@gmail.com",
		"packager@ipkg.io",
		"Subject",
		"Test Message",
	}
	if err := en.Notify(); err == nil {
		t.Fatalf("Connection should be refused: %s", err)
	}
}
