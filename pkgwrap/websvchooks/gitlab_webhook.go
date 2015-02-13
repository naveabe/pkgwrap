package websvchooks

import (
	"encoding/json"
	"github.com/naveabe/pkgwrap/pkgwrap/initscript"
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
	"github.com/naveabe/pkgwrap/pkgwrap/specer"
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
	RequestChan chan specer.PackageRequest
}

func (g *GitlabWebHook) parseTagEvent(payload []byte) (*specer.PackageRequest, error) {
	var (
		glEvt  GitlabTagEvent
		err    error
		pkgReq *specer.PackageRequest
	)

	if err = json.Unmarshal(payload, &glEvt); err != nil {
		return pkgReq, err
	}
	// Version may be in build config
	version, err := GetVersionFromRef(glEvt.Ref)
	if err != nil {
		g.Logger.Error.Printf("%s\n", err)
	}

	pkgReq = specer.NewPackageRequest(glEvt.Repository.Name)
	pkgReq.Version = version

	pkgReq.Package, err = specer.NewUserPackage(pkgReq.Name, pkgReq.Version,
		pkgReq.Name+"/"+pkgReq.Version+"/"+pkgReq.Name, initscript.BasicRunnable{})

	pkgReq.Package.URL = glEvt.Repository.Homepage
	pkgReq.Package.BuildType = specer.BUILDTYPE_SOURCE

	tagbranch, err := GetTagFromRef(glEvt.Ref)
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

	g.Logger.Debug.Printf("Queueing request: %#v ...\n", pkgReq)
	g.RequestChan <- *pkgReq

	rslt, _ := json.MarshalIndent(pkgReq, "", "  ")
	g.Logger.Trace.Printf("%s\n", rslt)
	w.WriteHeader(200)
}
