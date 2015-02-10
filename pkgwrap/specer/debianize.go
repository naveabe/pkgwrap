package specer

import (
	"fmt"
	"github.com/naveabe/pkgwrap/pkgwrap/templater"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type DEBSpec struct {
	PackageMetadata

	Url string

	Summary string

	BuildDeps string
	Deps      string

	CurrentDateTime string // used for changelog
}

func NewDEBSpec(name, version string) *DEBSpec {
	return &DEBSpec{
		PackageMetadata: PackageMetadata{
			Name:    name,
			Version: version,
		},
		CurrentDateTime: time.Now().Format(time.RFC1123Z),
	}
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

func BuildDebStructure(tmplMgr *templater.TemplatesManager, uPkg *UserPackage, distro Distribution, dstDir string) error {
	var (
		dspec = NewDEBSpec(uPkg.Name, uPkg.Version)
		err   error
	)

	dspec.Packager = uPkg.Packager
	dspec.Url = uPkg.URL
	dspec.Release = uPkg.Release
	dspec.Description = uPkg.Description
	dspec.Summary = uPkg.Name + " " + uPkg.Version
	dspec.BuildDeps = strings.Join(distro.BuildDeps, " ")
	dspec.Deps = strings.Join(distro.Deps, " ")

	if err := dspec.WriteDirStructure(dstDir); err != nil {
		return err
	}

	tbldr, err := tmplMgr.DebControlTemplateBuilder("debian")
	if err != nil {
		return err
	}

	outFile := dstDir + "/debian/control"
	fh, err := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer fh.Close()

	if err = tbldr.Build(dspec, fh); err != nil {
		return err
	}

	cbldr, err := tmplMgr.DebChangelogTemplateBuilder("debian")
	if err != nil {
		return err
	}

	outClog := dstDir + "/debian/changelog"
	cfh, err := os.OpenFile(outClog, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer cfh.Close()

	return cbldr.Build(dspec, cfh)
}
