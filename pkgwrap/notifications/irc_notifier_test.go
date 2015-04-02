package notifications

import (
	"testing"
)

var (
	testIrcStr  = "chat.freenode.net#metrilyx"
	testIrcStr2 = "chat.freenode.net##ipkgio"
	testMsg     = `Multi
line
Test`
)

func Test_IRCNotifier_Notify(t *testing.T) {
	nc, err := NewIRCNotifierFromString(testIrcStr)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if err = nc.Notify(testMsg); err != nil {
		t.Fatalf("%s", err)
	}
}
