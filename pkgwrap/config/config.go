package config

import (
	"encoding/json"
	"io/ioutil"
)

type BuildJobTracker struct {
	Enabled   bool
	Datastore DatastoreConfig `json:"datastore"`
}

type DatastoreConfig struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Index       string `json:"index"`
	MappingFile string `json:"mapping_file"`
}

type RepoMountConfig struct {
	// Base dir for source repo
	SrcBase    string `json:"src_base"`
	MountPoint string `json:"mount_point"`
}

type BuilderConfig struct {
	ImagePrefix string            `json:"image_prefix"`
	Mounts      map[string]string `json:"mounts"`
	RepoMount   RepoMountConfig   `json:"repo_mount"`
}

type HttpEndpointsConfig struct {
	Gitlab  string `json:"gitlab"`
	Github  string `json:"github"`
	Builder string `json:"builder"`
	Repo    string `json:"repo"`
	Jobs    string `json:"jobs"`
}

type AppConfig struct {
	//Templates  string              `json:"templates"`
	//ImageFiles string              `json:"imagefiles"`
	Repository string              `json:"repository"`
	DataDir    string              `json:"data_dir"`
	Port       int                 `json:"port"`
	Endpoints  HttpEndpointsConfig `json:"endpoints"`
	Builder    BuilderConfig       `json:"builder"`
	JobTracker BuildJobTracker     `json:"job_tracker"`
}

func (a *AppConfig) TemplatesDir() string {
	return a.DataDir + "/" + "templates"
}
func (a *AppConfig) ImageFilesDir() string {
	return a.DataDir + "/" + "imagefiles"
}
func (a *AppConfig) BinDir() string {
	return a.DataDir + "/" + "bin"
}

func LoadConfigFromFile(cfgfile string) (*AppConfig, error) {
	cfg := AppConfig{}

	b, err := ioutil.ReadFile(cfgfile)
	if err != nil {
		return &cfg, err
	}
	if err = json.Unmarshal(b, &cfg); err != nil {
		return &cfg, err
	}

	return &cfg, nil
}
