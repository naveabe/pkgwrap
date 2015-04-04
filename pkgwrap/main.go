package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/naveabe/pkgwrap/pkgwrap/builder"
	"github.com/naveabe/pkgwrap/pkgwrap/config"
	"github.com/naveabe/pkgwrap/pkgwrap/core/httphandlers"
	ghhandlers "github.com/naveabe/pkgwrap/pkgwrap/core/httphandlers/github"
	glhandlers "github.com/naveabe/pkgwrap/pkgwrap/core/httphandlers/gitlab"
	"github.com/naveabe/pkgwrap/pkgwrap/core/request"
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
	"github.com/naveabe/pkgwrap/pkgwrap/notifications"
	"github.com/naveabe/pkgwrap/pkgwrap/repository"
	"github.com/naveabe/pkgwrap/pkgwrap/templater"
	"github.com/naveabe/pkgwrap/pkgwrap/tracker"
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

func RunBuildRequest(bldCfg config.BuilderConfig, datastore *tracker.TrackerStore,
	repo repository.BuildRepository, pkgReq *request.PackageRequest, tmplMgr *templater.TemplatesManager, logger *logging.Logger) error {

	tBld, err := builder.NewTargetedPackageBuild(bldCfg, repo, pkgReq)
	if err != nil {
		return err
	}

	if err = tBld.SetupEnv(tmplMgr); err != nil {
		return err
	}

	builds := tBld.StartBuilds(DOCKER_URI)
	logger.Info.Printf("Containers started: %d\n", len(builds))
	logger.Trace.Printf("Containers details: %s\n", builds)

	if err = datastore.UpdateRequest(tBld.BuildRequest.Id, *tBld.BuildRequest); err != nil {
		return fmt.Errorf("Failed to update build state: %s", err)
	} else {
		b, _ := json.MarshalIndent(tBld.BuildRequest, "", "  ")
		logger.Trace.Printf("Build started: %s\n", b)
	}

	return nil
}

func StartWebServices(cfg *config.AppConfig, repo repository.BuildRepository,
	logger *logging.Logger, reqChan chan request.PackageRequest, datastore *tracker.TrackerStore) {

	httphandlers.SetupBuildHandler(cfg, datastore, repo, reqChan, logger)
	httphandlers.SetupRepoHandler(cfg, repo, logger)
	httphandlers.SetupLogHandler(cfg, DOCKER_URI, logger)

	glhandlers.SetupGitlabHandler(cfg, reqChan, logger)
	ghhandlers.SetupGithubHandlers(cfg, reqChan, logger)

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

func StartUserNotifier(dstore *tracker.TrackerStore, logger *logging.Logger) *notifications.NotificationProcessor {
	np := notifications.NewNotificationProcessor(dstore, logger)
	go np.Start()
	return np
}

func main() {
	flag.Parse()

	var (
		logger = logging.NewStdLogger()
		// channel receiving package requests
		pkgReqChan = make(chan request.PackageRequest)
		// package requests will get sent on this channel for builds
		//bldReqChan = make(chan *request.PackageRequest)
		notifier *notifications.NotificationProcessor
		// global config
		cfg *config.AppConfig
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

	// send user email/irc after build exit.
	notifier = StartUserNotifier(datastore, logger)

	// Used for updating state changes.
	go tracker.StartEventMonitor(DOCKER_URI, datastore, notifier.Listener, logger)

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

		/* this is what gets queued */
		//bldReqChan <- &pkgReq
		if err = RunBuildRequest(cfg.Builder, datastore, repo, &pkgReq, &tmplMgr, logger); err != nil {
			logger.Error.Printf("%s", err)
		}
	}
}
