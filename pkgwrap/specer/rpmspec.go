package specer

import (
	"fmt"
	"github.com/naveabe/pkgwrap/pkgwrap/templater"
	//"path/filepath"
	"strings"
	"time"
)

const (
	DEFAULT_LICENSE = "GPLv2"
	DEFAULT_GROUP   = "Development/Tools"
)

var (
	DEFAULT_PREP = []string{
		`[ -f "$RPM_BUILD_DIR/%{NAME}" ] && rm -rf $RPM_BUILD_DIR/%{NAME}`,
		`cp -a $RPM_SOURCE_DIR/%{NAME} $RPM_BUILD_DIR`,
		`cd $RPM_BUILD_DIR/%{NAME}`,
		`chmod -Rf a+rX,u+w,g-w,o-w .`,
	}

	DEFAULT_CLEAN = []string{""}

	/* Copy init script if one does not exist in the user provided tarball */
	DEFAULT_BUILD = []string{""}

	/* Copy everything retaining perms. */
	DEFAULT_INSTALL = []string{
		"rsync -vrlptH %{NAME}/ $RPM_BUILD_ROOT/",
	}
)

type RPMSpec struct {
	PackageMetadata

	Year int `json:"year"`

	Url     string `json:"url"`
	Summary string `json:"summary"`

	License string `json:"license"`
	Group   string `json:"group"`

	Source string `json:"source"`

	BuildRequires string `json:"build_requires"`
	Requires      string `json:"requires"`

	Prep    []string `json:"omit"`
	Build   []string `json:"omit"`
	Install []string `json:"omit"`
	Clean   []string `json:"omit"`

	PreInstall    []string `json:"pre_install,omitempty"`
	PostInstall   []string `json:"post_install,omitempty"`
	PreUninstall  []string `json:"pre_uninstall,omitempty"`
	PostUninstall []string `json:"post_uninstall,omitempty"`

	// TODO: source - analyze package //
	Files     []string `json:"omit"`
	Changelog []string `json:"changelog"`
}

func NewRPMSpec(name, version string) (*RPMSpec, error) {
	var (
		defaultRpmString = fmt.Sprintf("%s %s", name, version)
		//err              error
	)
	rspec := RPMSpec{
		PackageMetadata: PackageMetadata{
			Name:        name,
			Version:     version,
			Release:     DEFAULT_RELEASE,
			Description: defaultRpmString,
			Packager:    DEFAULT_PACKAGER,
		},
		Year:    time.Now().Year(),
		Summary: defaultRpmString,
		License: DEFAULT_LICENSE,

		Prep:      DEFAULT_PREP,
		Clean:     DEFAULT_CLEAN,
		Group:     DEFAULT_GROUP,
		Build:     DEFAULT_BUILD,
		Install:   DEFAULT_INSTALL,
		Files:     make([]string, 0),
		Changelog: make([]string, 0),
	}
	return &rspec, nil
}

func BuildRPMSpec(tmplMgr *templater.TemplatesManager, pkgReq *UserPackage, distro Distribution, dstDir string) (*RPMSpec, error) {
	spec, err := NewRPMSpec(pkgReq.Name, pkgReq.Version)
	if err != nil {
		return spec, err
	}

	spec.BuildRequires = strings.Join(distro.BuildDeps, " ")
	spec.Requires = strings.Join(distro.Deps, " ")

	//spec.Release = pkgReq.Release
	spec.Release = distro.PkgRelease

	spec.Source = pkgReq.Name
	spec.Url = pkgReq.URL
	spec.Files = pkgReq.FileList
	spec.Packager = pkgReq.Packager

	// Exclude as build happens outside of RPM env.
	//spec.Build = distro.UserBuildCmd

	spec.PreInstall = distro.PreInstall
	spec.PostInstall = distro.PostInstall
	spec.PreUninstall = distro.PreUninstall
	spec.PostUninstall = distro.PostUninstall

	if err = tmplMgr.WriteSpecFile(pkgReq.Name, string(distro.Name), spec, dstDir); err != nil {
		return spec, err
	}

	return spec, nil
}
