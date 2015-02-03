package specer

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const (
	DEFAULT_PKG_VERSION = "0.0.1"
)

type PackageRequest struct {
	Name    string
	Version string

	Package       *UserPackage   `yaml:"Package"`
	Distributions []Distribution `yaml:"Distributions"`
}

func NewPackageRequest(name string) *PackageRequest {
	return &PackageRequest{
		Name:          name,
		Version:       DEFAULT_PKG_VERSION,
		Distributions: make([]Distribution, 0),
		Package:       NewUserPackageWithName(name),
	}
}

func NewPackageRequestFromYamlConfig(yml string) (*PackageRequest, error) {
	pkgreq := NewPackageRequest("")

	b, err := ioutil.ReadFile(yml)
	if err != nil {
		return pkgreq, err
	}
	if err = yaml.Unmarshal(b, pkgreq); err != nil {
		return pkgreq, err
	}
	return pkgreq, nil
}

func (p *PackageRequest) Validate(deepInspection bool) error {
	if p.Package == nil {
		return fmt.Errorf("Package data not provided!")
	}
	if p.Package.URL == "" {
		return fmt.Errorf("Package url not provided!")
	}

	if deepInspection {
		if len(p.Distributions) <= 0 {
			return fmt.Errorf("Distribution/s not specified!")
		}
	}

	if p.Package.Packager == "" || p.Package.Packager == "mock" {
		pkgr, err := p.Package.PackagerFromURL()
		if err == nil {
			p.Package.Packager = pkgr
		}
	}

	if p.Package.Name == "" {
		p.Package.Name = p.Name
	}
	if p.Package.Version == "" {
		p.Package.Version = p.Version
	}

	if p.Package.Release == 0 {
		p.Package.Release = DEFAULT_RELEASE
	}

	return nil
}
