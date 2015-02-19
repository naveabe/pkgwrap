package builder

import (
	"fmt"
	//"github.com/fsouza/go-dockerclient"
	"github.com/naveabe/pkgwrap/pkgwrap/config"
	"github.com/naveabe/pkgwrap/pkgwrap/initscript"
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
	"github.com/naveabe/pkgwrap/pkgwrap/repository"
	"github.com/naveabe/pkgwrap/pkgwrap/specer"
	"github.com/naveabe/pkgwrap/pkgwrap/templater"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

/*
	Manages overall build process including containers for a given package.
	i.e. one container per package type
*/
type TargetedPackageBuild struct {
	BuildRequest     *specer.PackageRequest
	Repository       repository.BuildRepository
	DistroContainers map[string]*ContainerRunner

	cfg    config.BuilderConfig
	logger *logging.Logger
}

func NewTargetedPackageBuild(cfg config.BuilderConfig, repo repository.BuildRepository, pkgReq *specer.PackageRequest) (*TargetedPackageBuild, error) {
	var (
		b = TargetedPackageBuild{
			BuildRequest: pkgReq,
			Repository:   repo,
			logger:       logging.NewStdLogger(),
			cfg:          cfg,
		}
		err error
	)
	// Create 1 container per distro build
	if err = b.buildDistroContainers(); err != nil {
		return &b, err
	}
	return &b, nil
}

func (b *TargetedPackageBuild) buildDistroContainers() error {

	b.DistroContainers = make(map[string]*ContainerRunner)
	for _, distro := range b.BuildRequest.Distributions {
		if err := b.Add(distro, b.BuildRequest.Package); err != nil {
			return err
		}
		b.logger.Trace.Printf("Distro initialized: %s %s\n", distro.Name, distro.Release)
	}
	return nil
}

/*
	Start containers for each distribution

	Returns:
		[]string : container id's

	Todo:
	uri : should be a pool of docker uri's
*/
func (b *TargetedPackageBuild) StartBuilds(uri string) []map[string]string {
	cntInfo := make([]map[string]string, 0)

	for id, dRun := range b.DistroContainers {

		if dkrCntr, err := dRun.Start(uri); err == nil {
			b.logger.Trace.Printf("Starting container: %s (%s)\n", id, b.BuildRequest.Package.Packager)
			b.logger.Trace.Printf("Config: %#v\n", dkrCntr.HostConfig)

			cntInfo = append(cntInfo, map[string]string{"id": dkrCntr.ID, "label": id})
			for i, d := range b.BuildRequest.Distributions {
				if d.Label() == id {
					b.BuildRequest.Distributions[i].Id = dkrCntr.ID
					break
				}
			}
		} else {
			b.logger.Error.Printf("Failed to start container %s: %s\n", id, err)
		}
	}
	return cntInfo
}

func (b *TargetedPackageBuild) SetupEnv(tmplMgr *templater.TemplatesManager) error {
	var (
		err error = nil
	)

	if err = b.Repository.Clean(b.BuildRequest.Package); err != nil {
		return err
	}

	if b.BuildRequest.Package.BuildType == specer.BUILDTYPE_BIN {
		// Uncompress if binary
		if err = b.BuildRequest.Package.Uncompress(b.Repository.RepoDir); err != nil {
			return err
		}
	} else {
		b.logger.Trace.Printf("Cloning: name: %s;  tagbranch: %s\n", b.BuildRequest.Name, b.BuildRequest.Package.TagBranch)

		// Git clone if build type source
		dstDir := b.Repository.BuildDir(b.BuildRequest.Package) + "/" + b.BuildRequest.Package.Name

		//b.logger.Trace.Printf("Cloning %s\n", b.BuildRequest.Package.TagBranch)
		if err = b.BuildRequest.Package.CloneRepo(dstDir); err != nil {
			//b.logger.Error.Printf("%s\n", err)
			return err
		}
		b.logger.Trace.Printf("Cloned: %s %s\n", b.BuildRequest.Name, b.BuildRequest.Version)
		// Check for .yml in project root - read, validate, re-evaluate distro
		if err = b.readProjectPkgwrapConfig(); err != nil {
			return err
		}
	}

	if err = b.buildInitScript(tmplMgr); err != nil {
		return err
	}
	//b.logger.Trace.Printf("%v\n", b.BuildRequest.Package)
	if err = b.prepPerDistroBuilds(tmplMgr); err != nil {
		return err
	}

	b.logger.Debug.Printf("Re-processing distros: %d\n", len(b.DistroContainers))
	if err = b.buildDistroContainers(); err != nil {
		return err
	}
	return nil
}

func (b *TargetedPackageBuild) buildInitScript(tmplMgr *templater.TemplatesManager) error {
	if b.BuildRequest.Package.InitScript != nil && b.BuildRequest.Package.InitScript.Runnable.Path != "" {
		b.logger.Trace.Printf("Generating init script: %s %s\n", b.BuildRequest.Name, b.BuildRequest.Version)
		return initscript.BuildInitScript(tmplMgr, *b.BuildRequest.Package.InitScript,
			b.Repository.BuildDir(b.BuildRequest.Package))
	} else {
		b.logger.Info.Printf("Not creating startup script. No runnable path or initscript options specified!\n")
	}
	return nil
}

/*
	Read .pkgwrap.yml from the project root and re-evaluate
	distro containers.
*/
func (b *TargetedPackageBuild) readProjectPkgwrapConfig() error {
	bldConf := b.Repository.BuildConfig(b.BuildRequest.Package)
	b.logger.Trace.Printf("Reading project config: %s\n", bldConf)

	cBytes, err := ioutil.ReadFile(bldConf)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(cBytes, b.BuildRequest); err != nil {
		return err
	}

	if err = b.BuildRequest.Validate(true); err != nil {
		return err
	}
	b.logger.Debug.Printf("Project config read: %s\n", bldConf)

	return nil
}

/*
	Setup rpm/deb build structure needed to make the package.
*/
func (b *TargetedPackageBuild) prepPerDistroBuilds(tmplMgr *templater.TemplatesManager) error {
	var ptype specer.OSPackageType

	// 'distro' is by value
	for i, distro := range b.BuildRequest.Distributions {
		// Set next release
		nextRelease := b.Repository.NextRelease(b.BuildRequest.Package, distro.Label())
		if nextRelease < 1 {
			b.BuildRequest.Distributions[i].PkgRelease = 1
		} else if nextRelease > distro.PkgRelease {
			b.BuildRequest.Distributions[i].PkgRelease = nextRelease
		}
		b.logger.Trace.Printf("Set distro (%s) release: %d\n",
			distro.Label(), b.BuildRequest.Distributions[i].PkgRelease)
		b.logger.Debug.Printf("Processing distro: %s", distro.Label())

		// Process based on package type
		ptype = distro.PackageType()
		switch ptype {
		case specer.OS_PKG_TYPE_RPM:
			if err := b.setupRPMBuild(b.BuildRequest.Distributions[i], tmplMgr); err != nil {
				return err
			}
			break
		case specer.OS_PKG_TYPE_DEB:
			if err := b.setupDEBBuild(b.BuildRequest.Distributions[i], tmplMgr); err != nil {
				return err
			}
			break
		default:
			return fmt.Errorf("Package type not supported: %s", ptype)
		}
	}
	return nil
}

/*
	All pre-build setup before the .rpm build can start
*/
func (b *TargetedPackageBuild) setupRPMBuild(distro specer.Distribution, tmplMgr *templater.TemplatesManager) error {
	specDst := b.Repository.BuildDir(b.BuildRequest.Package) + "/" + distro.Label()
	// Write spec to repository
	_, err := specer.BuildRPMSpec(tmplMgr, b.BuildRequest.Package, distro, specDst)
	return err
}

/*
	All pre-build setup before the .deb build can start
*/
func (b *TargetedPackageBuild) setupDEBBuild(distro specer.Distribution, tmplMgr *templater.TemplatesManager) error {
	dstDir := b.Repository.BuildDir(b.BuildRequest.Package) + "/" + distro.Label()
	return specer.BuildDebStructure(tmplMgr, b.BuildRequest.Package, distro, dstDir)
}

func (b *TargetedPackageBuild) Add(distro specer.Distribution, pkg *specer.UserPackage) error {

	cRunner, err := NewContainerRunner(b.cfg, distro, pkg)
	if err != nil {
		return err
	}

	b.DistroContainers[distro.Label()] = cRunner
	return nil
}

func (b *TargetedPackageBuild) ListContainers() []string {
	clist := make([]string, len(b.DistroContainers))
	i := 0
	for k, _ := range b.DistroContainers {
		clist[i] = k
		i++
	}
	return clist
}
