package websvchooks

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	VERSION_RE, _ = regexp.Compile("([0-9\\.]+)")
)

func GetVersionFromRef(ref string) (string, error) {
	refParts := strings.Split(ref, "/")
	if refParts[1] == "tags" {
		mchArr := VERSION_RE.FindStringSubmatch(refParts[2])
		if len(mchArr) <= 0 {
			return "", fmt.Errorf("Could not determine version: %s", ref)
		}
		return mchArr[0], nil
	} else {
		return "", fmt.Errorf("Not tagged: %s", ref)
	}
}

func GetTagFromRef(ref string) (string, error) {

	refParts := strings.Split(ref, "/")
	if refParts[1] == "tags" {
		return refParts[2], nil
	} else {
		return "", fmt.Errorf("Not tagged: %s", ref)
	}
}
