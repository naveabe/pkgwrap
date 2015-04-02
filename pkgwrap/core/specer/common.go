package specer

const (
	DEFAULT_RELEASE  = 1
	DEFAULT_PACKAGER = "mock"
)

type PackageMetadata struct {
	Name     string `json:"name" yaml:"name"`
	Version  string `json:"version" yaml:"version"`
	Release  int64  `json:"release" yaml:"release"`
	Packager string `json:"packager" yaml:"packager"`

	Description string `json:"description"`
}
