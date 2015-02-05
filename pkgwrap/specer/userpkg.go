package specer

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/libgit2/git2go"
	"github.com/naveabe/pkgwrap/pkgwrap/initscript"
	"github.com/naveabe/pkgwrap/pkgwrap/repository"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	DEFAULT_PACKAGER = "mock"
	DEFAULT_RELEASE  = 1
)

type PkgBuildType string

const (
	BUILDTYPE_BIN    PkgBuildType = "binary"
	BUILDTYPE_SOURCE PkgBuildType = "source"
)

type PackageMetadata struct {
	Name     string `json:"name" yaml:"name"`
	Version  string `json:"version" yaml:"version"`
	Release  int64  `json:"release" yaml:"release"`
	Packager string `json:"packager" yaml:"packager"`

	Description string `json:"description"`
}

type UserPackage struct {
	PackageMetadata

	Path string `json:"path"`

	InitScript *initscript.BasicInitScript `json:"init_script,omitempty"`

	FileList []string `json:"files,omitempty"`

	URL string `json:"url" yaml:"url"`

	BuildEnv  string       `json:"build_env" yaml:"build_env"`
	BuildType PkgBuildType `json:"build_type"`

	// git tag/branch to checkout //
	TagBranch string `json:"tagbranch" yaml:"tagbranch"`
}

func NewUserPackageWithName(name string) *UserPackage {
	uPkg := UserPackage{
		PackageMetadata: PackageMetadata{
			Name:     name,
			Release:  DEFAULT_RELEASE,
			Packager: DEFAULT_PACKAGER,
		},
		TagBranch: "master",
	}
	uPkg.InitScript, _ = initscript.NewBasicInitScript(uPkg.Name)

	return &uPkg
}

func NewUserPackage(name, version, pkgpath string, runnable initscript.BasicRunnable) (*UserPackage, error) {
	var (
		uPkg = UserPackage{
			PackageMetadata: PackageMetadata{
				Name:     name,
				Version:  version,
				Release:  DEFAULT_RELEASE,
				Packager: DEFAULT_PACKAGER,
			},
			TagBranch: "master",
			Path:      pkgpath,
			BuildType: BUILDTYPE_BIN,
		}
		err error
	)

	if runnable.Path != "" {
		if uPkg.InitScript, err = initscript.NewBasicInitScript(uPkg.Name); err != nil {
			return &uPkg, err
		}
		uPkg.InitScript.Runnable = runnable
	}

	if !strings.HasSuffix(uPkg.Path, ".tgz") && !strings.HasSuffix(uPkg.Path, ".tar.gz") {
		return &uPkg, fmt.Errorf("Only .tgz and .tar.gz supported!")
	}
	return &uPkg, nil
}

func (u *UserPackage) SourceRepoName() (string, error) {
	p := strings.Split(u.URL, "/")
	if len(p) < 3 {
		return "", fmt.Errorf("Invalid URL: %s", u.URL)
	}
	return p[2], nil
}

func (u *UserPackage) PackagerFromURL() (string, error) {
	parts := strings.Split(u.URL, "/")
	if len(parts) > 2 {
		return parts[len(parts)-2], nil
	} else {
		return "", fmt.Errorf("Could not determine packager from url: %s\n", u.URL)
	}
}

func (u *UserPackage) CloneRepo(repoBase string) error {
	repoUrl := u.URL + ".git"
	//b.logger.Trace.Printf("Cloning: %s\n", repoUrl)
	copts := git.CloneOptions{
		CheckoutBranch: u.TagBranch,
	}
	_, err := git.Clone(repoUrl, repoBase+"/"+u.Path, &copts)
	if err != nil {
		return err
	}
	return nil
}

/*
	Params:
		repo : Repository
		distroLabel : e.g. centos, centos-6, ubuntu-12.04 ...
*/
func (u *UserPackage) AutoSetRelease(repo repository.BuildRepository, distroLabel string) {
	nextRelease := repo.NextRelease(u.Name, u.Version, distroLabel)
	if nextRelease > u.Release {
		u.Release = nextRelease
	}
}

func (u *UserPackage) Uncompress(repoBase string) error {
	dst := filepath.Dir(repoBase + "/" + u.Path)

	if strings.HasSuffix(u.Path, ".tgz") || strings.HasSuffix(u.Path, ".tar.gz") {
		// Gzip tarball //
		return u.unGzipTar(repoBase, dst)
	} else if strings.HasSuffix(u.Path, ".tbz2") || strings.HasSuffix(u.Path, ".tar.bz2") {
		//  Bzip tarball //
		return fmt.Errorf("Compression not yet supported: %s", u.Path)
	} else if strings.HasSuffix(u.Path, ".zip") {
		// Zip file //
		return fmt.Errorf("Compression not yet supported: %s", u.Path)
	} else {
		return fmt.Errorf("Compression not supported: %s", u.Path)
	}
}

func (u *UserPackage) unGzipTar(repoBase, dstDir string) error {
	irdr, err := os.Open(repoBase + "/" + u.Path)
	if err != nil {
		return err
	}
	// ungzip
	gz, err := gzip.NewReader(irdr)
	if err != nil {
		return err
	}
	defer gz.Close()

	return u.untar(gz, dstDir)
}

func (u *UserPackage) untar(r io.Reader, dst string) error {
	tr := tar.NewReader(r)
	for {
		entry, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		// start writing out uncompressed tarball //
		if entry.FileInfo().IsDir() {
			os.MkdirAll(dst+string(os.PathSeparator)+entry.Name, entry.FileInfo().Mode())
		} else {
			/* Remove pkg dirname from the prefix */
			u.FileList = append(u.FileList, strings.TrimPrefix(entry.Name, u.Name))
			/* Create file */
			fw, err := os.OpenFile(dst+string(os.PathSeparator)+entry.Name,
				os.O_CREATE|os.O_WRONLY|os.O_TRUNC, entry.FileInfo().Mode())
			if err != nil {
				return err
			}
			defer fw.Close()

			if _, err = io.Copy(fw, tr); err != nil {
				return err
			}
		}
	}
	return nil
}
