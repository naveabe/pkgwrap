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
	Runnable BasicRunnable `json:"runnable"`
	Logfile  string        `json:"logfile"`
}

func NewBasicInitScript(name string) (*BasicInitScript, error) {
	var (
		//err error
		bis = BasicInitScript{
			Name:     name,
			Runnable: BasicRunnable{},
		}
	)
	bis.Logfile = DEFAULT_LOG_DIR + "/" + bis.Name + ".log"

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
