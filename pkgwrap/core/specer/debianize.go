package specer

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/naveabe/pkgwrap/pkgwrap/core/request"
	"github.com/naveabe/pkgwrap/pkgwrap/templater"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const DEB_BIN_VERSION = "2.0"
const DEFAULT_SHELL = "#!/bin/bash"

type DEBSpec struct {
	PackageMetadata

	Shell string

	Url string

	Arch string

	Summary string

	BuildDeps string
	Deps      string

	PreInstall    []string
	PostInstall   []string
	PreUninstall  []string
	PostUninstall []string
	//CurrentDateTime string // used for changelog
}

func NewDEBSpec(name, version string) *DEBSpec {
	return &DEBSpec{
		PackageMetadata: PackageMetadata{
			Name:    name,
			Version: version,
		},
		Shell: DEFAULT_SHELL,
		//CurrentDateTime: time.Now().Format(time.RFC1123Z),
	}
}

func (d *DEBSpec) WriteDirStructure(dstDir string) error {
	debDir := dstDir + "/" + "debian"
	os.MkdirAll(debDir, 0755)
	return d.writeDebianBinary(debDir)

}

func (d *DEBSpec) writeDebianBinary(dstDir string) error {
	return ioutil.WriteFile(dstDir+"/"+"debian-binary",
		[]byte(fmt.Sprintf("%s\n", DEB_BIN_VERSION)), 0755)
}

func (d *DEBSpec) WriteScripts(tmplMgr *templater.TemplatesManager, dstDir string) ([]string, error) {
	var (
		bldr    *templater.TemplateBuilder
		err     error
		scripts = make([]string, 0)
	)

	if d.PreInstall != nil && len(d.PreInstall) > 0 {
		if bldr, err = tmplMgr.DebScriptTemplateBuilder("debian", "preinst"); err != nil {
			return scripts, err
		}
		if err = bldr.WriteNormalizedFile(d, dstDir+"/preinst"); err != nil {
			return scripts, err
		}
		scripts = append(scripts, "preinst")
	}
	if d.PostInstall != nil && len(d.PostInstall) > 0 {
		if bldr, err = tmplMgr.DebScriptTemplateBuilder("debian", "postinst"); err != nil {
			return scripts, err
		}
		if err = bldr.WriteNormalizedFile(d, dstDir+"/postinst"); err != nil {
			return scripts, err
		}
		scripts = append(scripts, "postinst")
	}
	if d.PreUninstall != nil && len(d.PreUninstall) > 0 {
		if bldr, err = tmplMgr.DebScriptTemplateBuilder("debian", "prerm"); err != nil {
			return scripts, err
		}
		if err = bldr.WriteNormalizedFile(d, dstDir+"/prerm"); err != nil {
			return scripts, err
		}
		scripts = append(scripts, "prerm")
	}
	if d.PostUninstall != nil && len(d.PostUninstall) > 0 {
		if bldr, err = tmplMgr.DebScriptTemplateBuilder("debian", "postrm"); err != nil {
			return scripts, err
		}
		if err = bldr.WriteNormalizedFile(d, dstDir+"/postrm"); err != nil {
			return scripts, err
		}
		scripts = append(scripts, "postrm")
	}
	return scripts, nil
}

func WriteControlArchive(dstDir string, scripts []string) error {
	archFiles := append([]string{"control"}, scripts...)

	outFh, err := os.OpenFile(dstDir+"/control.tar.gz", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer outFh.Close()

	//var fWriter io.WriteCloser = outFh
	fWriter := gzip.NewWriter(outFh)
	defer fWriter.Close()

	tw := tar.NewWriter(fWriter)
	defer tw.Close()

	for _, s := range archFiles {
		fInfo, err := os.Stat(dstDir + "/" + s)
		if err != nil {
			return err
		}

		hdr := tar.Header{
			Name:    fInfo.Name(),
			Size:    fInfo.Size(),
			Mode:    int64(fInfo.Mode()),
			ModTime: fInfo.ModTime(),
		}
		if err = tw.WriteHeader(&hdr); err != nil {
			return err
		}

		fh, err := os.Open(dstDir + "/" + s)
		if err != nil {
			return err
		}
		defer fh.Close()

		if _, err = io.Copy(tw, fh); err != nil {
			return err
		}
	}

	return nil
}

func WriteDebControlFile(tmplMgr *templater.TemplatesManager, data interface{}, dstDir string) error {
	tbldr, err := tmplMgr.DebControlTemplateBuilder("debian")
	if err != nil {
		return err
	}

	outFile := dstDir + "/control"
	fh, err := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer fh.Close()

	return tbldr.Build(data, fh)
}

func BuildDebStructure(tmplMgr *templater.TemplatesManager, uPkg *request.UserPackage, distro request.Distribution, dstDir string) error {
	var (
		dspec   = NewDEBSpec(uPkg.Name, uPkg.Version)
		err     error
		scripts []string
	)
	// TODO: auto configure
	dspec.Arch = "amd64"
	dspec.Packager = uPkg.Packager
	dspec.Url = uPkg.URL

	//dspec.Release = uPkg.Release
	dspec.Release = distro.PkgRelease

	dspec.Description = uPkg.Description
	dspec.Summary = uPkg.Name + " " + uPkg.Version
	dspec.BuildDeps = strings.Join(distro.BuildDeps, ", ")
	dspec.Deps = strings.Join(distro.Deps, ", ")

	dspec.PreInstall = distro.PreInstall
	dspec.PostInstall = distro.PostInstall
	dspec.PreUninstall = distro.PreUninstall
	dspec.PostUninstall = distro.PostUninstall

	if err = dspec.WriteDirStructure(dstDir); err != nil {
		return err
	}

	if scripts, err = dspec.WriteScripts(tmplMgr, dstDir+"/debian"); err != nil {
		return err
	}

	if err = WriteDebControlFile(tmplMgr, dspec, dstDir+"/debian"); err != nil {
		return err
	}

	return WriteControlArchive(dstDir+"/debian", scripts)
}
