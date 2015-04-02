package initscript

import (
	"github.com/naveabe/pkgwrap/pkgwrap/templater"
)

const (
	DEFAULT_LOG_DIR = "/var/log"
)

type InitScriptTemplateData struct {
	Name         string
	RunnablePath string
	RunnableArgs string
	Logfile      string
}

type BasicInitScript struct {
	Name     string        `json:"name"`
	Runnable BasicRunnable `json:"runnable" yaml:"runnable"`
	Logfile  string        `json:"logfile"`
}

func (b *BasicInitScript) SetName(name string) {
	b.Name = name
	b.Logfile = DEFAULT_LOG_DIR + "/" + b.Name + ".log"
}

func NewBasicInitScript(name string) (*BasicInitScript, error) {
	bis := BasicInitScript{Runnable: BasicRunnable{}}
	bis.SetName(name)
	return &bis, nil
}

func BuildInitScript(tmplMgr *templater.TemplatesManager, bis BasicInitScript, outDir string) error {
	tmplData := InitScriptTemplateData{
		Name:         bis.Name,
		RunnablePath: bis.Runnable.Path,
		RunnableArgs: bis.Runnable.Args,
		Logfile:      DEFAULT_LOG_DIR + "/" + bis.Name + ".log",
	}

	return tmplMgr.WriteStartupFile(bis.Name, tmplData, outDir)
}
