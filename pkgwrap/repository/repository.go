package repository

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const PROJECT_CONFIG_NAME = ".pkgwrap.yml"

/*
	BuildRepository contains info and structure for
	where user packages will be stored.

	This is also the place where generated startup script,
	spec file/s and the uncompressed package will be stored
	along with the generated .rpm's and .deb's
*/
type BuildRepository struct {
	RepoDir string
}

func (b *BuildRepository) PackagePath(pkgr, shortPath string) string {
	return b.RepoDir + "/" + pkgr + "/" + shortPath
}

func (b *BuildRepository) ListPackageVersions(pkgr, pkgname string) ([]string, error) {
	files, err := ioutil.ReadDir(b.RepoDir + "/" + pkgr + "/" + pkgname)
	if err != nil {
		return make([]string, 0), err
	}

	flist := make([]string, len(files))
	for i, f := range files {
		flist[i] = f.Name()
	}
	return flist, nil
}

func (b *BuildRepository) ListPackages(pkgr, pkgname, pkgversion, distroLabel string) ([]string, error) {
	files, err := ioutil.ReadDir(b.RepoDir + "/" + pkgr + "/" + pkgname + "/" + pkgversion + "/" + distroLabel)
	if err != nil {
		return make([]string, 0), err
	}

	var flist []string
	if len(files) > 2 {
		flist = make([]string, len(files)-2)
		i := 0
		for _, f := range files {
			/* TODO: exclusion based on distro */
			if strings.HasSuffix(f.Name(), ".spec") || f.Name() == "RELEASE" {
				continue
			}
			flist[i] = f.Name()
			i++
		}
	} else {
		flist = make([]string, 0)
	}
	return flist, nil
}

func (b *BuildRepository) ListPackageDistros(pkgr, pkgname, pkgversion string) ([]string, error) {
	files, err := ioutil.ReadDir(b.RepoDir + "/" + pkgr + "/" + pkgname + "/" + pkgversion)
	if err != nil {
		return make([]string, 0), err
	}

	var flist []string
	if len(files) > 2 {
		flist = make([]string, len(files)-2)
		i := 0
		for _, f := range files {
			if strings.HasPrefix(f.Name(), pkgname) {
				continue
			}
			flist[i] = f.Name()
			i++
		}
	} else {
		flist = make([]string, 0)
	}
	return flist, nil
}

func (b *BuildRepository) GetPackagePathForDistro(pkgr, name, version, distroLabel, pkg string) (string, error) {
	pkgPath := b.RepoDir + "/" + pkgr + "/" + name + "/" + version + "/" + distroLabel + "/" + pkg
	if _, err := os.Stat(pkgPath); err != nil {
		return pkgPath, err
	}
	return pkgPath, nil
}

func (b *BuildRepository) BuildDir(pkgr, pkgName, pkgVersion string) string {
	return b.RepoDir + "/" + pkgr + "/" + pkgName + "/" + pkgVersion
}

func (b *BuildRepository) BuildConfig(pkgr, pkgName, pkgVersion string) string {
	return b.RepoDir + "/" + pkgr + "/" + pkgName + "/" + pkgVersion + "/" + pkgName + "/" + PROJECT_CONFIG_NAME
}

/*
	Last release stored in repository under RELEASE for each package
	type.

 	Returns:
 		Release number
 		-1 if unknown
*/
func (b *BuildRepository) LastRelease(pkgr, pkgName, pkgVersion, distroLabel string) int64 {
	if b, err := ioutil.ReadFile(b.RepoDir + "/" + pkgr + "/" + pkgName + "/" + pkgVersion + "/" + distroLabel + "/RELEASE"); err == nil {
		relstr := strings.TrimSpace(string(b))
		if val, err := strconv.ParseInt(relstr, 10, 64); err == nil {
			return val
		}
	}
	return -1
}

/*
	Next release version that will be built.

	Params:
		pkgName    : package name
		pkgVersion : package version
		pkgType    : "rpm" or "deb"

	Returns:
		-1 if unknown
*/
func (b *BuildRepository) NextRelease(pkgr, pkgName, pkgVersion, distroLabel string) int64 {
	lastRel := b.LastRelease(pkgr, pkgName, pkgVersion, distroLabel)
	if lastRel == -1 {
		return lastRel
	}
	return lastRel + 1
}

/*
	Remove cloned package repo or uncompressed tarball upload
	by the user
*/
func (b *BuildRepository) Clean(pkgr, pkgName, pkgVersion string) error {
	rmDir := b.RepoDir + "/" + pkgr + "/" + pkgName + "/" + pkgVersion + "/" + pkgName
	if _, err := os.Stat(rmDir); err == nil {
		return os.RemoveAll(rmDir)
	}
	return nil
}
