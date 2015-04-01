package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/naveabe/pkgwrap/pkgwrap/builder"
	"github.com/naveabe/pkgwrap/pkgwrap/config"
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
	"github.com/naveabe/pkgwrap/pkgwrap/repository"
	"github.com/naveabe/pkgwrap/pkgwrap/specer"
	"github.com/naveabe/pkgwrap/pkgwrap/templater"
	"github.com/naveabe/pkgwrap/pkgwrap/tracker"
	"github.com/naveabe/pkgwrap/pkgwrap/websvc"
	"github.com/naveabe/pkgwrap/pkgwrap/websvchooks"
	"net/http"
	"os"
)

var (
	CONFIG_FILE = flag.String("c", "pkgwrapd.conf", "Configuration file")
	LOGLEVEL    = flag.String("l", "info", "Log level [ error | warning | info | debug | trace ]")
)

const (
	DOCKER_HOST_PORT = "localhost:5555"
	DOCKER_URI       = "tcp://" + DOCKER_HOST_PORT
)

func PrepTargetedBuild(bldrCfg config.BuilderConfig, repo repository.BuildRepository,
	pkgReq *specer.PackageRequest, tmplMgr *templater.TemplatesManager) (*builder.TargetedPackageBuild, error) {

	tBuild, err := builder.NewTargetedPackageBuild(bldrCfg, repo, pkgReq)
	if err != nil {
		return tBuild, err
	}

	if err = tBuild.SetupEnv(tmplMgr); err != nil {
		return tBuild, err
	}

	return tBuild, nil
}

func SetupGithubOauthHandler(cfgs []config.CodeRepoConfig, logger *logging.Logger) {
	for _, v := range cfgs {
		if v.Type == "github" {
			ghOauthHandler := websvc.NewGithubOauthHandler(v, logger)
			websvc.NewRestHandler(v.LocalEndpoint, ghOauthHandler, logger)
			logger.Warning.Printf("Github oauth API: %s\n", v.LocalEndpoint)
			return
		}
	}
	logger.Warning.Printf("Github config not found. Disabling...!\n")

}

func StartWebServices(cfg *config.AppConfig, repo repository.BuildRepository, logger *logging.Logger,
	reqChan chan specer.PackageRequest, datastore *tracker.TrackerStore) {

	methodHandler := websvc.PkgBuilderMethodHandler{
		Config:      cfg,
		Repository:  repo,
		Logger:      logger,
		RequestChan: reqChan, // testing
		Datastore:   datastore,
	}

	websvc.NewRestHandler(cfg.Endpoints.Builder, &methodHandler, logger)
	logger.Warning.Printf("Builder API: %s\n", cfg.Endpoints.Builder)

	// Log handler
	logHdlr, err := websvc.NewLogHandler(DOCKER_URI, logger)
	if err != nil {
		logger.Error.Fatalf("%s\n", err)
	}
	websvc.NewRestHandler(cfg.Endpoints.Logs, logHdlr, logger)
	logger.Warning.Printf("Logs API: %s\n", cfg.Endpoints.Logs)

	// Gitlab webhook
	if cfg.Endpoints.Gitlab != "" {
		glHandle := websvchooks.GitlabWebHook{logger, reqChan}
		http.Handle(cfg.Endpoints.Gitlab, &glHandle)
		logger.Warning.Printf("Gitlab service: %s\n", cfg.Endpoints.Gitlab)
	} else {
		logger.Warning.Printf("Gitlab service disabled!\n")
	}

	if cfg.Endpoints.Github != "" {
		ghHandle := websvchooks.GithubWebHook{logger, reqChan}
		http.Handle(cfg.Endpoints.Github, &ghHandle)
		logger.Warning.Printf("Github service: %s\n", cfg.Endpoints.Github)
	} else {
		logger.Warning.Printf("Github service disabled!\n")
	}

	SetupGithubOauthHandler(cfg.CodeRepos, logger)

	repoHandle := websvc.NewRepoHandler(repo, logger)
	websvc.NewRestHandler(cfg.Endpoints.Repo, repoHandle, logger)
	logger.Warning.Printf("Repository API: %s\n", cfg.Endpoints.Repo)

	if cfg.Webroot != "" {
		logger.Warning.Printf("HTTP root directory: %s\n", cfg.Webroot)
		http.Handle("/", http.FileServer(http.Dir(cfg.Webroot)))
	} else {
		logger.Warning.Printf("Web UI disabled!\n")
	}

	logger.Warning.Printf("Starting web service: http://%s:%d\n", cfg.Host, cfg.Port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), nil); err != nil {
		logger.Error.Printf("%s\n", err)
		os.Exit(2)
	}
}

func StartEventMonitor(dstore *tracker.TrackerStore, logger *logging.Logger) {
	if dem, err := tracker.NewDockerEventMonitor(DOCKER_URI, dstore, logger); err == nil {
		if err := dem.Start(); err != nil {
			logger.Error.Fatalf("%s\n", err)
		}
	} else {
		logger.Error.Fatalf("%s\n", err)
	}
}

func main() {
	flag.Parse()

	var (
		logger = logging.NewStdLogger()
		// channel receiving package requests
		pkgReqChan = make(chan specer.PackageRequest)
		cfg        *config.AppConfig
		// datastore
		datastore *tracker.TrackerStore
		err       error
		// Build repo with packages
		repo repository.BuildRepository
		// rpm and dep templates
		tmplMgr templater.TemplatesManager
	)

	logger.SetLogLevel(*LOGLEVEL)

	if cfg, err = config.LoadConfigFromFile(*CONFIG_FILE); err != nil {
		logger.Error.Printf("Failed to load config: %s\n", err)
		os.Exit(1)
	}

	repo = repository.BuildRepository{cfg.Repository}
	tmplMgr = templater.TemplatesManager{cfg.TemplatesDir()}

	if datastore, err = tracker.NewTrackerStore(&cfg.Tracker.Datastore, logger); err != nil {
		logger.Error.Printf("Datastore initialization failed: %s\n", err)
		os.Exit(2)
	}
	// HTTP server /api/builder
	go StartWebServices(cfg, repo, logger, pkgReqChan, datastore)
	// Used for updating state changes.
	go StartEventMonitor(datastore, logger)

	/* Prep and start builds */
	for {
		pkgReq := <-pkgReqChan

		logger.Info.Printf("Package request: name=%s version=%s release=%d build_type=%s\n",
			pkgReq.Name, pkgReq.Version, pkgReq.Package.Release, pkgReq.Package.BuildType)

		if pkgReq.Id, err = datastore.AddRequest(pkgReq); err != nil {
			logger.Error.Printf("Failed to add request: %s\n", err)
			continue
		}

		pReqBytes, _ := json.MarshalIndent(pkgReq, "", "  ")
		logger.Trace.Printf("Request added: %s\n", pReqBytes)

		tBld, err := PrepTargetedBuild(cfg.Builder, repo, &pkgReq, &tmplMgr)
		if err != nil {
			logger.Error.Printf("Failed to prep build: %s\n", err)
			continue
		}

		builds := tBld.StartBuilds(DOCKER_URI)
		logger.Info.Printf("Containers started: %d\n", len(builds))
		logger.Trace.Printf("Containers details: %s\n", builds)

		if err = datastore.UpdateRequest(tBld.BuildRequest.Id, *tBld.BuildRequest); err != nil {
			logger.Error.Printf("Failed to update loaded request: %s\n", err)
		} else {
			b, _ := json.MarshalIndent(tBld.BuildRequest, "", "  ")
			logger.Trace.Printf("Request loaded: %s\n", b)
		}
	}
}
