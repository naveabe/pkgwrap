package notifications

import (
	"testing"
)

var (
	testEmail     = "user@domain.io"
	testBuildNoti BuildNotifications
)

func Test_BuildNotifications_GetIrcNotifiers(t *testing.T) {

	testBuildNoti = BuildNotifications{
		[]string{testIrcStr, testIrcStr2},
		[]string{testEmail}}
	/*
		nots, err := testBuildNoti.GetIrcNotifiers()
		if err != nil {
			t.Fatalf("%s", err)
		}
		if nots[0].Channel != "#metrilyx" {
			t.Fatalf("Failed to parse: %s", nots[0])
		}

		if nots[1].Channel != "##ipkgio" {
			t.Fatalf("Failed to parse: %s", nots[1])
		}

		t.Logf("%s", nots)
	*/
	t.Logf("%s", testBuildNoti)
}

/*
func Test_BuildNotifications_GetEmailNotifiers(t *testing.T) {
	list, err := testBuildNoti.GetEmailNotifiers()
	if err != nil {
		t.Fatalf("%s", err)
	}
	t.Logf("%s", list)
}
*/
