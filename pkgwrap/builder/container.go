package builder

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/naveabe/pkgwrap/pkgwrap/config"
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
	"github.com/naveabe/pkgwrap/pkgwrap/specer"
	"strings"
)

type ContainerRunner struct {
	Distro  specer.Distribution
	Package *specer.UserPackage

	client *docker.Client

	ContainerConfig docker.CreateContainerOptions

	dockerCont *docker.Container

	cfg config.BuilderConfig

	logger *logging.Logger
}

func NewContainerRunner(builderCfg config.BuilderConfig, distro specer.Distribution, pkg *specer.UserPackage) (*ContainerRunner, error) {
	var (
		c   ContainerRunner
		err error = nil
	)
	c = ContainerRunner{
		Distro:  distro,
		Package: pkg,
		cfg:     builderCfg,
		logger:  logging.NewStdLogger(),
	}

	c.ContainerConfig, err = c.initContainerConfig()

	return &c, err
}

func (c *ContainerRunner) initContainerConfig() (docker.CreateContainerOptions, error) {
	var (
		opts = docker.CreateContainerOptions{}
	)

	opts.HostConfig = &docker.HostConfig{
		Binds: c.getMounts(),
	}

	// Container image and entry point
	opts.Config = &docker.Config{
		Image: c.ContainerImage(),
		Cmd:   []string{c.Distro.BuildCommand(), c.Package.Name, c.Package.TagBranch},
	}

	repoName, err := c.Package.SourceRepoName()
	if err != nil {
		return opts, err
	}

	opts.Config.Env = []string{
		"REPO=" + repoName,
		"BUILD_USER=" + c.Package.Packager,
		"BUILD_ENV=" + c.Package.BuildEnv,
		"BUILD_CMD=" + strings.Join(c.Distro.BuildCmd, " ; "),
		"BUILD_DEPS=" + strings.Join(c.Distro.BuildDeps, " "),
		"PKG_DEPS=" + strings.Join(c.Distro.Deps, " "),
		"PKG_DISTRO=" + c.Distro.Label(),
		fmt.Sprintf("PKG_TYPE=%s", c.Distro.PackageType()),
		fmt.Sprintf("PKG_RELEASE=%d", c.Distro.PkgRelease),
		fmt.Sprintf("PKG_VERSION=%s", c.Package.Version),
		fmt.Sprintf("BUILD_TYPE=%s", c.Package.BuildType),
	}

	c.logger.Trace.Printf("Container env: %v\n", opts.Config.Env)

	return opts, nil
}

func (c *ContainerRunner) ContainerImage() string {
	if c.Distro.Release == "" {
		return fmt.Sprintf("%s%s", c.cfg.ImagePrefix, c.Distro.Name)
	} else {
		return fmt.Sprintf("%s%s:%s", c.cfg.ImagePrefix, c.Distro.Name, c.Distro.Release)
	}
}

func (c *ContainerRunner) Start(uri string) (*docker.Container, error) {
	var err error

	if c.client, err = docker.NewClient(uri); err != nil {
		return c.dockerCont, err
	}

	if c.dockerCont, err = c.client.CreateContainer(c.ContainerConfig); err != nil {
		return c.dockerCont, err
	}

	if err = c.client.StartContainer(c.dockerCont.ID, c.ContainerConfig.HostConfig); err != nil {
		return c.dockerCont, err
	}

	return c.dockerCont, nil
}

/*
	Mounts
*/
func (c *ContainerRunner) getMounts() []string {

	out := make([]string, len(c.cfg.Mounts)+1)
	i := 0
	for k, v := range c.cfg.Mounts {
		out[i] = k + ":" + v
		i++
	}

	out[len(out)-1] = c.cfg.RepoMount.SrcBase + "/" + c.Package.VersionBaseDir() +
		":" + c.cfg.RepoMount.MountPoint
	return out
}
