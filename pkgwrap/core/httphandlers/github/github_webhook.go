package github

import (
	"encoding/json"
	"fmt"
	"github.com/naveabe/pkgwrap/pkgwrap/config"
	"github.com/naveabe/pkgwrap/pkgwrap/core/httphandlers"
	"github.com/naveabe/pkgwrap/pkgwrap/core/request"
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
	"io/ioutil"
	"net/http"
)

type GithubRepoOwner struct {
	Login   string
	Id      int64
	HtmlUrl string `json:"html_url"`
}

type GithubRepo struct {
	Name        string
	Id          string          `json:"id"`
	FullName    string          `json:"full_name"`
	Description string          `json:"description"`
	Owner       GithubRepoOwner `json:"owner"`
}

type GithubSender struct {
	Login string
	Id    int
	Type  string
}

type GithubCreateDeleteEvent struct {
	GenericGithubEvent
	Description string `json:"description"`
	Ref         string `json:"ref"`
	RefType     string `json:"ref_type"`
}

type GenericGithubEvent struct {
	EventType  string
	Repository GithubRepo   `json:"repository"`
	Sender     GithubSender `json:"sender"`
}

type GithubPushEvent struct {
	GenericGithubEvent
	Ref string `json:"ref"`
}

type GithubWebHook struct {
	Logger *logging.Logger
	// This channel will be read to get PackageRequests
	RequestChan chan request.PackageRequest
}

func (g *GithubWebHook) parsePushEvent(evtType string, payload []byte) (*request.PackageRequest, error) {
	var (
		pkgReq  *request.PackageRequest
		pushEvt = GithubPushEvent{GenericGithubEvent: GenericGithubEvent{EventType: evtType}}
		err     error
	)

	if err = json.Unmarshal(payload, &pushEvt); err != nil {
		return pkgReq, err
	}

	pkgReq = request.NewPackageRequest(pushEvt.Repository.Name)
	g.Logger.Trace.Printf("Push: %s\n", pushEvt)

	tagbranch, err := httphandlers.GetTagFromRef(pushEvt.Ref)
	if err != nil {
		g.Logger.Warning.Printf("Could not determine tag: %s\n", err)
	} else {
		pkgReq.Package.TagBranch = tagbranch
	}

	pkgReq.Package.URL = "https://github.com/" + pushEvt.Repository.FullName
	pkgReq.Package.Packager = pushEvt.Sender.Login

	version, err := httphandlers.GetVersionFromRef(pushEvt.Ref)
	if err != nil {
		g.Logger.Warning.Printf("Could not determine version: %s", pushEvt.Ref)
	} else {
		pkgReq.Version = version
		pkgReq.Package.Version = version
	}

	pkgReq.Package.Path = fmt.Sprintf("%s/%s/%s", pkgReq.Name, pkgReq.Version, pkgReq.Name)
	return pkgReq, pkgReq.Validate(false)
}

func (g *GithubWebHook) parseCreateEvent(evtType string, payload []byte) (*request.PackageRequest, error) {
	var (
		pkgReq      *request.PackageRequest
		createEvent = GithubCreateDeleteEvent{
			GenericGithubEvent: GenericGithubEvent{
				EventType: evtType,
			},
		}
		err error
	)

	if err = json.Unmarshal(payload, &createEvent); err != nil {
		return pkgReq, err
	} else if createEvent.RefType != "tag" {
		return pkgReq, fmt.Errorf("No tag on commit!")
	} else {
		g.Logger.Trace.Printf("Create: %s\n", createEvent)

		pkgReq := request.NewPackageRequest(createEvent.Repository.Name)
		pkgReq.Package.TagBranch = createEvent.Ref
		pkgReq.Package.URL = "https://github.com/" + createEvent.Repository.FullName
		pkgReq.Package.Packager = createEvent.Repository.Owner.Login

		mchArr := httphandlers.VERSION_RE.FindStringSubmatch(createEvent.Ref)
		if len(mchArr) <= 0 {
			g.Logger.Warning.Printf("Could not determine version: %s", createEvent.Ref)
		} else {
			pkgReq.Version = mchArr[0]
			pkgReq.Package.Version = mchArr[0]
		}
		pkgReq.Package.Path = fmt.Sprintf("%s/%s/%s", pkgReq.Name, pkgReq.Version, pkgReq.Name)

		return pkgReq, pkgReq.Validate(false)
	}
}

func (g *GithubWebHook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		pkgReq   *request.PackageRequest
		evtType  = r.Header.Get("X-Github-Event")
		respCode int
		err      error
	)
	// extra check
	if len(r.Header.Get("X-Github-Delivery")) != 36 {
		g.Logger.Error.Printf("X-Github-Delivery - not found!\n")
		w.WriteHeader(401)
		return
	}

	payloadBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		g.Logger.Error.Printf("%s\n", err)
		w.WriteHeader(400)
		return
	}

	switch evtType {
	case "create":
		pkgReq, err = g.parseCreateEvent(evtType, payloadBytes)
		if err != nil {
			if err.Error() == "No tag on commit!" {
				g.Logger.Warning.Printf("%s\n", err)
				respCode = 200
			} else {
				g.Logger.Error.Printf("%s\n", err)
				w.WriteHeader(400)
				return
			}
		} else {
			respCode = 200
		}
		break
	case "push":
		pkgReq, err = g.parsePushEvent(evtType, payloadBytes)
		if err != nil {
			g.Logger.Error.Printf("%s ==> %s\n", err, payloadBytes)
			w.WriteHeader(400)
			return
		} else {
			respCode = 200
		}
		break
	default:
		g.Logger.Debug.Printf("Skipped event (no case): %s\n", evtType)
		g.Logger.Trace.Printf("Skipped (%s): %s\n", evtType, payloadBytes)
		w.WriteHeader(200)
		return
	}

	g.Logger.Debug.Printf("Queueing request: %#v ...\n", pkgReq)
	g.RequestChan <- *pkgReq

	rslt, _ := json.MarshalIndent(pkgReq, "", "  ")
	g.Logger.Trace.Printf("%s\n", rslt)
	w.WriteHeader(respCode)
}

func SetupGithubWebHook(cfg *config.AppConfig, reqChan chan request.PackageRequest, logger *logging.Logger) {
	if cfg.Endpoints.Github != "" {
		ghHandle := GithubWebHook{logger, reqChan}
		http.Handle(cfg.Endpoints.Github, &ghHandle)
		logger.Warning.Printf("Github service: %s\n", cfg.Endpoints.Github)
	} else {
		logger.Warning.Printf("Github service disabled!\n")
	}
}

/*
	Helper function for github hook and oauth
*/
func SetupGithubHandlers(cfg *config.AppConfig, reqChan chan request.PackageRequest, logger *logging.Logger) {
	SetupGithubWebHook(cfg, reqChan, logger)
	SetupGithubOauthHandler(cfg.CodeRepos, logger)
}
