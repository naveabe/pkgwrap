package templater

import (
	"io"
	"os"
	"text/template"
)

type TemplateBuilder struct {
	Name     string
	Filepath string
	tmpl     *template.Template
}

/* TODO: load all templates on startup */
func NewTemplateBuilder(name, filepath string) (*TemplateBuilder, error) {
	var (
		t = TemplateBuilder{
			Name:     name,
			Filepath: filepath,
		}
		err error
	)
	if t.tmpl, err = template.New(name).ParseFiles(filepath); err != nil {
		return &t, err
	}

	return &t, nil
}

func (t *TemplateBuilder) Build(input interface{}, output io.Writer) error {
	return t.tmpl.Execute(output, input)
}

func (t *TemplateBuilder) WriteNormalizedFile(data interface{}, dstFile string) error {
	fh, err := os.OpenFile(dstFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer fh.Close()

	return t.Build(data, fh)
}

type TemplatesManager struct {
	TemplatesDir string
}

/*
func (tm *TemplatesManager) DebRulesTemplateBuilder(distroName string) (*TemplateBuilder, error) {
	return NewTemplateBuilder("rules", tm.TemplatesDir+"/"+distroName+"/rules")
}

func (tm *TemplatesManager) DebChangelogTemplateBuilder(distroName string) (*TemplateBuilder, error) {
	return NewTemplateBuilder("changelog", tm.TemplatesDir+"/"+distroName+"/changelog")
}
*/
func (tm *TemplatesManager) DebScriptTemplateBuilder(distroName, script string) (*TemplateBuilder, error) {
	return NewTemplateBuilder(script, tm.TemplatesDir+"/"+distroName+"/"+script)
}

func (tm *TemplatesManager) DebControlTemplateBuilder(distroName string) (*TemplateBuilder, error) {
	return NewTemplateBuilder("control", tm.TemplatesDir+"/"+distroName+"/control")
}

func (tm *TemplatesManager) SpecTemplateBuilder(distroName string) (*TemplateBuilder, error) {
	return NewTemplateBuilder(distroName+".spec", tm.TemplatesDir+"/spec/"+distroName+".spec")
}

func (tm *TemplatesManager) WriteSpecFile(specName, distroName string, spec interface{}, outDir string) error {
	os.MkdirAll(outDir, 0755)

	tbldr, err := tm.SpecTemplateBuilder(distroName)
	if err != nil {
		return err
	}

	outFile := outDir + "/" + specName + ".spec"
	fh, err := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer fh.Close()

	return tbldr.Build(spec, fh)
}

func (tm *TemplatesManager) StartupTemplateBuilder() (*TemplateBuilder, error) {
	return NewTemplateBuilder("startup.sh", tm.TemplatesDir+"/startup/startup.sh")
}

func (tm *TemplatesManager) WriteStartupFile(startupName string, startupData interface{}, outDir string) error {
	os.MkdirAll(outDir, 0755)

	tbldr, err := tm.StartupTemplateBuilder()
	if err != nil {
		return err
	}

	fh, err := os.OpenFile(outDir+"/"+startupName+".service", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer fh.Close()

	return tbldr.Build(startupData, fh)
}
