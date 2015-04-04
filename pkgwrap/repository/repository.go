package repository

import (
	//"fmt"
	"github.com/naveabe/pkgwrap/pkgwrap/core/request"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const PROJECT_CONFIG_NAME = ".pkgwrap.yml"

/*
type ProjectUrlInfo struct {
	URL     string
	Repo    string
	User    string
	Project string
}

func NewProjectUrlInfo(url string) (*ProjectUrlInfo, error) {
	pp := ProjectUrlInfo{URL: url}

	parts := strings.Split(url, "/")
	if len(parts) < 4 {
		return &pp, fmt.Errorf("Invalid URL: %s", url)
	}
	pp.Repo = parts[2]
	pp.User = parts[3]
	pp.Project = parts[4]
	return &pp, nil
}

func (p *ProjectUrlInfo) CreatePath(baseDir string) {
	makedir := fmt.Sprintf("%s/%s/%s/%s", baseDir, p.Repo, p.User, p.Project)
	os.MkdirAll(makedir, 0755)
}
*/
type RepoPackage struct {
	Name    string  `json:"name"`
	ModTime float64 `json:"mtime"`
	Size    int64   `json:"size"`
}

/*
	BuildRepository contains info and structure for
	where user packages will be stored.

	This is also the place where generated startup script,
	spec file/s and the uncompressed package will be stored
	along with the generated .rpm's and .deb's
*/
type BuildRepository struct {
	// Repo base directory
	RepoDir string
}

func (b *BuildRepository) PackagePath(pkgr, shortPath string) string {
	return b.RepoDir + "/" + pkgr + "/" + shortPath
}

func (b *BuildRepository) ListUserProjects(repo, pkgr string) ([]string, error) {
	flist := make([]string, 0)

	files, err := ioutil.ReadDir(b.RepoDir + "/" + repo + "/" + pkgr)
	if err != nil {
		return make([]string, 0), err
	}

	for _, f := range files {
		if f.IsDir() {
			flist = append(flist, f.Name())
		}
	}

	return flist, nil
}

func (b *BuildRepository) ListPackageVersions(repo, pkgr, pkgname string) ([]string, error) {
	files, err := ioutil.ReadDir(b.RepoDir + "/" + repo + "/" + pkgr + "/" + pkgname)
	if err != nil {
		return make([]string, 0), err
	}

	flist := make([]string, len(files))
	for i, f := range files {
		flist[i] = f.Name()
	}
	return flist, nil
}

func (b *BuildRepository) ListPackages(repo, pkgr, pkgname, pkgversion, distroLabel string) ([]RepoPackage, error) {
	//flist := make([]string, 0)
	flist := make([]RepoPackage, 0)

	files, err := ioutil.ReadDir(b.RepoDir + "/" + repo + "/" + pkgr +
		"/" + pkgname + "/" + pkgversion + "/" + distroLabel)
	if err != nil {
		return flist, err
	}

	for _, f := range files {

		if strings.HasSuffix(f.Name(), ".rpm") || strings.HasSuffix(f.Name(), ".deb") {
			flist = append(flist,
				RepoPackage{
					ModTime: float64(f.ModTime().UnixNano()) / 1000000000,
					Name:    f.Name(),
					Size:    f.Size(), // bytes
				})

			//flist = append(flist, f.Name())
		}
	}
	return flist, nil
}

func (b *BuildRepository) ListPackageDistros(repo, pkgr, pkgname, pkgversion string) ([]string, error) {
	flist := make([]string, 0)

	files, err := ioutil.ReadDir(b.RepoDir + "/" + repo + "/" + pkgr + "/" + pkgname + "/" + pkgversion)
	if err != nil {
		return flist, err
	}

	for _, f := range files {
		if strings.HasPrefix(f.Name(), pkgname) {
			continue
		}
		flist = append(flist, f.Name())
	}
	return flist, nil
}

func (b *BuildRepository) GetPackagePathForDistro(repo, pkgr, name, version, distroLabel, pkg string) (string, error) {
	pkgPath := b.RepoDir + "/" + repo + "/" + pkgr + "/" + name + "/" + version + "/" + distroLabel + "/" + pkg
	if _, err := os.Stat(pkgPath); err != nil {
		return pkgPath, err
	}
	return pkgPath, nil
}

func (b *BuildRepository) BuildDir(pkg *request.UserPackage) string {
	return b.RepoDir + "/" + pkg.VersionBaseDir()
	//return b.RepoDir + "/" + pkg.Packager + "/" + pkg.Name + "/" + pkg.Version
}

func (b *BuildRepository) BuildConfig(pkg *request.UserPackage) string {
	return b.BuildDir(pkg) + "/" + pkg.Name + "/" + PROJECT_CONFIG_NAME
}

/*
	Last release stored in repository under RELEASE for each package
	type.
*/
func (b *BuildRepository) LastRelease(pkg *request.UserPackage, distroLabel string) int64 {
	if b, err := ioutil.ReadFile(b.BuildDir(pkg) + "/" + distroLabel + "/RELEASE"); err == nil {

		relstr := strings.TrimSpace(string(b))
		if val, err := strconv.ParseInt(relstr, 10, 64); err == nil {
			return val
		}
	}
	return -1
}

/*
	Next release version that will be built.
*/
func (b *BuildRepository) NextRelease(pkg *request.UserPackage, distroLabel string) int64 {
	lastRel := b.LastRelease(pkg, distroLabel)
	if lastRel == -1 {
		return lastRel
	}
	return lastRel + 1
}

/*
	Remove cloned package repo or uncompressed tarball upload
	by the user
*/
func (b *BuildRepository) Clean(pkg *request.UserPackage) error {
	rmDir := b.BuildDir(pkg) + "/" + pkg.Name
	if _, err := os.Stat(rmDir); err == nil {
		return os.RemoveAll(rmDir)
	}
	return nil
}
