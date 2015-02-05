package specer

import (
	"fmt"
)

/* Todo: this should come from config */
const DISTRO_BIN_PREFIX = "/opt/pkgwrap/bin"

type OSPackageType string

const (
	OS_PKG_TYPE_RPM OSPackageType = "rpm"
	OS_PKG_TYPE_DEB OSPackageType = "deb"
)

type OSDistro string

const (
	DISTRO_CENTOS OSDistro = "centos"
	DISTRO_ORACLE OSDistro = "oracle"
	DISTRO_REDHAT OSDistro = "redhat"

	DISTRO_UBUNTU OSDistro = "ubuntu"
	DISTRO_DEBIAN OSDistro = "debian"

	DISTRO_AGNOSTIC OSDistro = "agnostic"
)

/*
 * 'latest' is implied
 */
var SUPPORTED_DISTROS = map[OSDistro]map[string]string{
	DISTRO_CENTOS: map[string]string{
		"6": "6",
	},
	DISTRO_ORACLE: map[string]string{
		"6": "6",
	},
	DISTRO_REDHAT: map[string]string{
		"6": "6",
	},
	DISTRO_UBUNTU: map[string]string{
		"12.04": "12.04",
		"14.04": "14.04",
	},
	DISTRO_DEBIAN: map[string]string{},
}

var DISTRO_BUILD_DIRS = map[OSDistro]string{
	DISTRO_CENTOS: "/home/%s/rpmbuild/SOURCES",
	DISTRO_ORACLE: "/home/%s/rpmbuild/SOURCES",
	DISTRO_REDHAT: "/home/%s/rpmbuild/SOURCES",

	DISTRO_UBUNTU: "/home/%s/debuild",
	DISTRO_DEBIAN: "/home/%s/debuild",
}

var DISTRO_BUILD_SCRIPT = map[OSDistro]string{
	DISTRO_CENTOS: DISTRO_BIN_PREFIX + "/run-rpm-build.sh",
	DISTRO_ORACLE: DISTRO_BIN_PREFIX + "/run-rpm-build.sh",
	DISTRO_REDHAT: DISTRO_BIN_PREFIX + "/run-rpm-build.sh",

	DISTRO_UBUNTU: DISTRO_BIN_PREFIX + "/build-deb.sh",
	DISTRO_DEBIAN: DISTRO_BIN_PREFIX + "/build-deb.sh",
}

var DISTRO_PKG_TYPE = map[OSDistro]OSPackageType{
	DISTRO_CENTOS: OS_PKG_TYPE_RPM,
	DISTRO_ORACLE: OS_PKG_TYPE_RPM,
	DISTRO_REDHAT: OS_PKG_TYPE_RPM,
	DISTRO_UBUNTU: OS_PKG_TYPE_DEB,
	DISTRO_DEBIAN: OS_PKG_TYPE_DEB,
}

type Distribution struct {
	Name    OSDistro `json:"name"`
	Release string   `json:"release"`

	// distro specific
	BuildDeps []string `json:"build_deps" yaml:"build_deps"`
	Deps      []string `json:"deps"`

	// Only used if source is built/compiled
	UserBuildCmd []string `json:"build_cmd" yaml:"build_cmd"`

	PreInstall    []string `json:"pre_install,omitempty" yaml:"pre_install"`
	PostInstall   []string `json:"post_install,omitempty" yaml:"post_install"`
	PreUninstall  []string `json:"pre_uninstall,omitempty" yaml:"pre_uninstall"`
	PostUninstall []string `json:"post_uninstall,omitempty" yaml:"post_uninstall"`

	buildDir string
}

func NewDistribution(name OSDistro, release string) (Distribution, error) {
	var (
		d   Distribution
		err error = nil
	)
	_, ok := SUPPORTED_DISTROS[name][release]
	// implied latest
	if release == "" || ok {
		d = Distribution{
			Name:    name,
			Release: release,
		}
		d.buildDir = DISTRO_BUILD_DIRS[d.Name]
	} else {
		err = fmt.Errorf("Not supported: %s %s", name, release)
	}

	return d, err
}

func (d *Distribution) Label() string {
	if d.Release == "" {
		return fmt.Sprintf("%s", d.Name)
	} else {
		return fmt.Sprintf("%s-%s", d.Name, d.Release)
	}
}

func (d *Distribution) PackageType() OSPackageType {
	return DISTRO_PKG_TYPE[d.Name]
}

func (d *Distribution) BuildDir(packager string) string {
	return fmt.Sprintf(d.buildDir, packager)
}

func (d *Distribution) BuildCommand() string {
	return DISTRO_BUILD_SCRIPT[d.Name]
}
