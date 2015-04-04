package gitlab

import (
	"encoding/json"
	"github.com/naveabe/pkgwrap/pkgwrap/config"
	"github.com/naveabe/pkgwrap/pkgwrap/core/httphandlers"
	"github.com/naveabe/pkgwrap/pkgwrap/core/initscript"
	"github.com/naveabe/pkgwrap/pkgwrap/core/request"
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
	"io/ioutil"
	"net/http"
)

type GitlabRepo struct {
	Homepage    string
	URL         string
	Name        string
	Description string
}

type GitlabAuthor struct {
	Email string
	Name  string
}

type GitlabCommit struct {
	Author GitlabAuthor
}

type GitlabTagEvent struct {
	Before     string
	After      string
	Ref        string
	UserId     int64  `json:"user_id"`
	Username   string `json:"user_name"`
	ProjectId  int64  `json:"project_id"`
	Repository GitlabRepo
	Commits    []GitlabCommit
}

type GitlabWebHook struct {
	Logger *logging.Logger
	// This channel will be read to get PackageRequests
	RequestChan chan request.PackageRequest
}

func (g *GitlabWebHook) parseTagEvent(payload []byte) (*request.PackageRequest, error) {
	var (
		glEvt  GitlabTagEvent
		err    error
		pkgReq *request.PackageRequest
	)

	if err = json.Unmarshal(payload, &glEvt); err != nil {
		return pkgReq, err
	}

	pkgReq = request.NewPackageRequest(glEvt.Repository.Name)
	// Version may be in build config
	pkgReq.Version, err = httphandlers.GetVersionFromRef(glEvt.Ref)
	if err != nil {
		pkgReq.Version = request.DEFAULT_PKG_VERSION
		g.Logger.Warning.Printf("Using default version - %s\n", err)
	} else {
		g.Logger.Debug.Printf("Using version (Gitlab): %s\n", pkgReq.Version)
	}

	pkgReq.Package, err = request.NewUserPackage(pkgReq.Name, pkgReq.Version,
		pkgReq.Name+"/"+pkgReq.Version+"/"+pkgReq.Name, initscript.BasicRunnable{})

	pkgReq.Package.URL = glEvt.Repository.Homepage
	pkgReq.Package.BuildType = request.BUILDTYPE_SOURCE

	tagbranch, err := httphandlers.GetTagFromRef(glEvt.Ref)
	if err != nil {
		g.Logger.Warning.Printf("%s - using default: %s\n", err, pkgReq.Package.TagBranch)
	} else {
		pkgReq.Package.TagBranch = tagbranch
	}

	pkgr, err := pkgReq.Package.PackagerFromURL()
	if err != nil {
		g.Logger.Error.Printf("%s\n", err)
	}
	pkgReq.Package.Packager = pkgr

	return pkgReq, nil
}

func (g *GitlabWebHook) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		g.Logger.Error.Printf("%s\n", err)
		w.WriteHeader(400)
		return
	}
	g.Logger.Trace.Printf("Gitlab event: %s\n", b)

	pkgReq, err := g.parseTagEvent(b)
	if err != nil {
		g.Logger.Error.Printf("%s\n", err)
		w.WriteHeader(400)
		return
	}

	if err = pkgReq.Validate(false); err != nil {
		g.Logger.Error.Printf("Validation failed: %s\n", err)
		w.WriteHeader(400)
		return
	}

	g.Logger.Debug.Printf("Queueing request: %#v ...\n", pkgReq)
	g.RequestChan <- *pkgReq

	rslt, _ := json.MarshalIndent(pkgReq, "", "  ")
	g.Logger.Trace.Printf("%s\n", rslt)
	w.WriteHeader(200)
}

func SetupGitlabHandler(cfg *config.AppConfig, reqChan chan request.PackageRequest, logger *logging.Logger) {
	if cfg.Endpoints.Gitlab != "" {
		glHandle := GitlabWebHook{logger, reqChan}
		http.Handle(cfg.Endpoints.Gitlab, &glHandle)
		logger.Warning.Printf("Gitlab service: %s\n", cfg.Endpoints.Gitlab)
	} else {
		logger.Warning.Printf("Gitlab service disabled!\n")
	}
}
