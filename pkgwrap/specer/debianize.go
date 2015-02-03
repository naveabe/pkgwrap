package specer

import (
	"fmt"
	"io/ioutil"
	"os"
)

type DEBSpec struct {
	PackageMetadata

	BuildDeps string
	Deps      string
}

func (d *DEBSpec) WriteDirStructure(dstDir string) error {
	debDir := dstDir + "/" + "debian"
	os.MkdirAll(debDir, 0755)

	err := d.writeCompat(debDir, 9) // debian compatibility version. should be 9 in almost all cases //
	if err != nil {
		return err
	}

	if err = d.writeRules(debDir); err != nil {
		return err
	}
	//return d.writeInstallFiles(debDir)
	return nil
}

func (d *DEBSpec) writeCompat(dstDir string, version int) error {
	return ioutil.WriteFile(dstDir+"/"+"compat",
		[]byte(fmt.Sprintf("%d", version)), 0755)
}

/*
 * Not to be changed as building will happen externally.
 */
func (d *DEBSpec) writeRules(dstDir string) error {
	rules := `#!/usr/bin/make -f
%:
    dh $@`

	return ioutil.WriteFile(dstDir+"/"+"rules",
		[]byte(rules), 0755)
}

/*
func (d *DEBSpec) writeInstallFiles(dstDir string) error {
	var (
		fh *os.File
	)

	fh, _ = os.OpenFile(fmt.Sprintf("%s/%s.install", dstDir, d.Name),
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)

	defer fh.Close()

	return nil
}
*/
/*
func (d *DEBSpec) WriteControl() {}
*/
