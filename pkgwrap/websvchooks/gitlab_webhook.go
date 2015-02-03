package websvchooks

import (
	"encoding/json"
	//"fmt"
	"github.com/naveabe/pkgwrap/pkgwrap/initscript"
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
	"github.com/naveabe/pkgwrap/pkgwrap/specer"
	"io/ioutil"
	"net/http"
	//"strings"
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

func (g *GitlabWebHook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var gl GitlabTagEvent

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		g.Logger.Error.Printf("%s\n", err)
		w.WriteHeader(400)
		return
	}

	if err = json.Unmarshal(b, &gl); err != nil {
		g.Logger.Error.Printf("%s\n", err)
		w.WriteHeader(400)
		return
	}
	// Version may be in build config
	version, err := GetVersionFromRef(gl.Ref)
	if err != nil {
		g.Logger.Error.Printf("%s\n", err)
	}

	pkgReq := specer.NewPackageRequest(gl.Repository.Name)
	pkgReq.Version = version

	pkgReq.Package, err = specer.NewUserPackage(pkgReq.Name, pkgReq.Version,
		pkgReq.Name+"/"+pkgReq.Version+"/"+pkgReq.Name, initscript.BasicRunnable{})

	pkgReq.Package.URL = gl.Repository.Homepage
	pkgReq.Package.BuildType = specer.BUILDTYPE_SOURCE

	tagbranch, err := GetTagFromRef(gl.Ref)
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

	//g.RequestChan <- pkgReq

	rslt, _ := json.MarshalIndent(pkgReq, "", "  ")
	g.Logger.Trace.Printf("%s\n", rslt)

	w.WriteHeader(200)
}
