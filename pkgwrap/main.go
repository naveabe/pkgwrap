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

func PrepTargetedBuild(bldrCfg config.BuilderConfig, repo repository.BuildRepository, pkgReq *specer.PackageRequest, tmplMgr *templater.TemplatesManager) (*builder.TargetedPackageBuild, error) {

	tBuild, err := builder.NewTargetedPackageBuild(bldrCfg, repo, pkgReq)
	if err != nil {
		return tBuild, err
	}

	if err = tBuild.SetupEnv(tmplMgr); err != nil {
		return tBuild, err
	}

	return tBuild, nil
}

func StartWebServices(cfg *config.AppConfig, repo repository.BuildRepository, logger *logging.Logger, reqChan chan specer.PackageRequest) {

	methodHandler := websvc.PkgBuilderMethodHandler{
		Config:      cfg,
		Repository:  repo,
		Logger:      logger,
		RequestChan: reqChan, // testing
	}

	websvc.NewRestHandler(cfg.Endpoints.Builder, &methodHandler, logger)
	logger.Warning.Printf("Builder API: %s\n", cfg.Endpoints.Builder)

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

	repoHandle := websvc.NewRepoHandler(repo, logger)
	websvc.NewRestHandler(cfg.Endpoints.Repo, repoHandle, logger)
	logger.Warning.Printf("Repository API: %s\n", cfg.Endpoints.Repo)

	logger.Warning.Printf("Starting service: http://0.0.0.0:%d\n", cfg.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil); err != nil {
		logger.Error.Printf("%s\n", err)
		os.Exit(2)
	}
}

func StartEventMonitor(dstore *tracker.EssDatastore, logger *logging.Logger) {
	dem := tracker.NewDockerEventMonitor("http://localhost:5555", dstore, logger)
	if err := dem.Start(); err != nil {
		logger.Error.Fatalf("%s\n", err)
	}
}

func main() {
	flag.Parse()

	var (
		logger     = logging.NewStdLogger()
		pkgReqChan = make(chan specer.PackageRequest)
		cfg        *config.AppConfig
		datastore  *tracker.EssDatastore
		err        error
	)

	logger.SetLogLevel(*LOGLEVEL)

	if cfg, err = config.LoadConfigFromFile(*CONFIG_FILE); err != nil {
		logger.Error.Printf("%s\n", err)
		os.Exit(1)
	}

	repo := repository.BuildRepository{cfg.Repository}
	tmplMgr := templater.TemplatesManager{cfg.TemplatesDir()}

	if cfg.JobTracker.Enabled {
		logger.Warning.Printf("Job tracker ENABLED!\n")

		datastore, err = tracker.NewEssDatastore(&cfg.JobTracker.Datastore, logger)
		if err != nil {
			logger.Error.Printf("Failed to init datastore: %s\n", err)
			os.Exit(2)
		}

		jobsHandle := websvc.NewJobsHandler(datastore, logger)
		websvc.NewRestHandler(cfg.Endpoints.Jobs, jobsHandle, logger)
		logger.Warning.Printf("Logs API: %s\n", cfg.Endpoints.Jobs)
	}

	// HTTP server /api/builder
	go StartWebServices(cfg, repo, logger, pkgReqChan)

	go StartEventMonitor(datastore, logger)

	// Avoid extra if statement in busy loop
	if cfg.JobTracker.Enabled {

		for {
			pkgReq := <-pkgReqChan

			logger.Info.Printf("Package request: name=%s version=%s release=%d build_type=%s\n",
				pkgReq.Name, pkgReq.Version, pkgReq.Package.Release, pkgReq.Package.BuildType)

			tBld, err := PrepTargetedBuild(cfg.Builder, repo, &pkgReq, &tmplMgr)
			if err != nil {
				logger.Error.Printf("%s\n", err)
				continue
			}

			buildIds := tBld.StartBuilds(DOCKER_URI)
			logger.Info.Printf("Containers started: %d %s\n", len(buildIds), buildIds)

			bJob := tracker.NewBuildJob(&pkgReq, buildIds, DOCKER_HOST_PORT)
			for i, _ := range bJob.Jobs {
				bJob.Jobs[i].Status = "start"
			}

			if err = bJob.Record(datastore); err != nil {
				logger.Error.Printf("%s\n", err)
			}

			b, _ := json.MarshalIndent(bJob, "", "  ")
			logger.Trace.Printf("Wrote job: %s\n", b)
		}
	} else {
		logger.Warning.Printf("Job tracker DISABLED!\n")

		for {
			pkgReq := <-pkgReqChan

			logger.Info.Printf("Package request: name=%s version=%s release=%d build_type=%s\n",
				pkgReq.Name, pkgReq.Version, pkgReq.Package.Release, pkgReq.Package.BuildType)

			tBld, err := PrepTargetedBuild(cfg.Builder, repo, &pkgReq, &tmplMgr)
			if err != nil {
				logger.Error.Printf("%s\n", err)
				continue
			}

			buildIds := tBld.StartBuilds(DOCKER_URI)
			logger.Info.Printf("Containers started: %d %s\n", len(buildIds), buildIds)

			bJob := tracker.NewBuildJob(&pkgReq, buildIds, DOCKER_HOST_PORT)

			b, _ := json.MarshalIndent(bJob, "", "  ")
			logger.Trace.Printf("Build job: %s\n", b)
		}
	}
}
