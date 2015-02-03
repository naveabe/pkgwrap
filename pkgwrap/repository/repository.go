package repository

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

/*
	BuildRepository contains info and structure for
	where user packages will be stored.

	This is also the place where generated startup script,
	spec file/s and the uncompressed package will be stored.
*/
type BuildRepository struct {
	RepoDir string
}

func (b *BuildRepository) PackagePath(shortPath string) string {
	return b.RepoDir + "/" + shortPath
}

func (b *BuildRepository) BuildDir(pkgName, pkgVersion string) string {
	return b.RepoDir + "/" + pkgName + "/" + pkgVersion
}

func (b *BuildRepository) BuildConfig(pkgName, pkgVersion string) string {
	return b.RepoDir + "/" + pkgName + "/" + pkgVersion + "/" + pkgName + "/.pkgwrap.yml"
}

/*
	Last release stored in repository under RELEASE for each package
	type.

 	Returns:
 		Release number
 		-1 if unknown
*/
func (b *BuildRepository) LastRelease(pkgName, pkgVersion, pkgType string) int64 {
	if b, err := ioutil.ReadFile(b.RepoDir + "/" + pkgName + "/" + pkgVersion + "/" + pkgType + "/RELEASE"); err == nil {
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
func (b *BuildRepository) NextRelease(pkgName, pkgVersion, pkgType string) int64 {
	lastRel := b.LastRelease(pkgName, pkgVersion, pkgType)
	if lastRel == -1 {
		return lastRel
	}
	return lastRel + 1
}

func (b *BuildRepository) Clean(pkgName, pkgVersion string) error {
	rmDir := b.RepoDir + "/" + pkgName + "/" + pkgVersion + "/" + pkgName
	if _, err := os.Stat(rmDir); err == nil {
		return os.RemoveAll(rmDir)
	}
	return nil
}
