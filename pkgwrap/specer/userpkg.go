package specer

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/libgit2/git2go"
	"github.com/naveabe/pkgwrap/pkgwrap/initscript"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	DEFAULT_PACKAGER = "mock"
	DEFAULT_RELEASE  = 1
)

/*
type AugmentableProperties struct {
	BuildDeps []string `json:"build_deps" yaml:"build_deps"`
	Deps      []string `json:"deps"`

	BuildCmd []string `json:"build_cmd" yaml:"build_cmd"`

	PreInstall    []string `json:"pre_install,omitempty" yaml:"pre_install"`
	PostInstall   []string `json:"post_install,omitempty" yaml:"post_install"`
	PreUninstall  []string `json:"pre_uninstall,omitempty" yaml:"pre_uninstall"`
	PostUninstall []string `json:"post_uninstall,omitempty" yaml:"post_uninstall"`
}
*/
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
	// git tag/branch to checkout //
	TagBranch string `json:"tagbranch" yaml:"tagbranch"`

	BuildEnv  string       `json:"build_env" yaml:"build_env"`
	BuildType PkgBuildType `json:"build_type"`

	//Augmentable Properties
	BuildDeps []string `json:"build_deps,omitempty" yaml:"build_deps"`
	Deps      []string `json:"deps,omitempty"`

	BuildCmd []string `json:"build_cmd" yaml:"build_cmd"`

	PreInstall    []string `json:"pre_install,omitempty" yaml:"pre_install"`
	PostInstall   []string `json:"post_install,omitempty" yaml:"post_install"`
	PreUninstall  []string `json:"pre_uninstall,omitempty" yaml:"pre_uninstall"`
	PostUninstall []string `json:"post_uninstall,omitempty" yaml:"post_uninstall"`
}

func NewUserPackageWithName(name string) *UserPackage {
	uPkg := UserPackage{
		PackageMetadata: PackageMetadata{
			Name:     name,
			Release:  DEFAULT_RELEASE,
			Packager: DEFAULT_PACKAGER,
		},
		TagBranch: "master",
		//AugmentableProperties: AugmentableProperties{},
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
			//AugmentableProperties: AugmentableProperties{},
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

func NewUserPackageFromURL(url string) (*UserPackage, error) {
	var (
		upkg = UserPackage{
			PackageMetadata: PackageMetadata{
				Release: DEFAULT_RELEASE,
			},
			URL:       url,
			TagBranch: "master",
			BuildType: BUILDTYPE_SOURCE,
			//AugmentableProperties: AugmentableProperties{},
		}
		err error
	)

	parts := strings.Split(upkg.URL, "/")
	if len(parts) < 5 {
		return &upkg, fmt.Errorf("Invalid URL: %s", upkg.URL)
	}
	upkg.Name = parts[4]
	upkg.Packager = parts[3]

	if upkg.InitScript, err = initscript.NewBasicInitScript(upkg.Name); err != nil {
		return &upkg, err
	}

	return &upkg, nil
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

func (u *UserPackage) VersionBaseDir() string {
	repoName, _ := u.SourceRepoName()
	return fmt.Sprintf("%s/%s/%s/%s", repoName, u.Packager, u.Name, u.Version)
}

func (u *UserPackage) CloneRepo(dstDir string) error {
	// Clone branch
	gitRepo, err := git.Clone(u.URL+".git", dstDir, &git.CloneOptions{CheckoutBranch: u.TagBranch})
	if err != nil {
		// Clone repo on branch failure
		if gitRepo, err = git.Clone(u.URL+".git", dstDir, &git.CloneOptions{}); err != nil {
			return err
		}

		ref, err := gitRepo.DwimReference(u.TagBranch)
		if err != nil {
			// Tag not found
			return err
		}

		if err = gitRepo.SetHeadDetached(ref.Target(), nil, ""); err != nil {
			return err
		}
	}

	return nil
}

func (u *UserPackage) Uncompress(repoBase string) error {
	src := repoBase + "/" + u.Packager + "/" + u.Path
	dst := filepath.Dir(src)

	if strings.HasSuffix(src, ".tgz") || strings.HasSuffix(src, ".tar.gz") {
		// Gzip tarball //
		return u.unGzipTar(src, dst)
	} else if strings.HasSuffix(src, ".tbz2") || strings.HasSuffix(src, ".tar.bz2") {
		//  Bzip tarball //
		return fmt.Errorf("Compression not yet supported: %s", u.Path)
	} else if strings.HasSuffix(src, ".zip") {
		// Zip file //
		return fmt.Errorf("Compression not yet supported: %s", u.Path)
	} else {
		return fmt.Errorf("Compression not supported: %s", u.Path)
	}
}

func (u *UserPackage) unGzipTar(srcTarGz, dstDir string) error {
	irdr, err := os.Open(srcTarGz)
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
		// Start writing out uncompressed tarball
		if entry.FileInfo().IsDir() {
			os.MkdirAll(dst+string(os.PathSeparator)+entry.Name, entry.FileInfo().Mode())
		} else {
			// Remove pkg dirname from the prefix
			u.FileList = append(u.FileList, strings.TrimPrefix(entry.Name, u.Name))
			// Create file
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
