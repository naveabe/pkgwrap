package httphandlers

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/naveabe/pkgwrap/pkgwrap/config"
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
	"net/http"
)

type LogHandler struct {
	DefaultMethodHandler
	dclient *docker.Client
	logger  *logging.Logger
}

func NewLogHandler(dockerUri string, logger *logging.Logger) (*LogHandler, error) {
	var (
		lh  = LogHandler{logger: logger}
		err error
	)

	if lh.dclient, err = docker.NewClient(dockerUri); err != nil {
		return &lh, err
	}
	return &lh, nil
}

func (d *LogHandler) GET(w http.ResponseWriter, r *http.Request, args ...string) (map[string]string, interface{}, int) {
	if len(args) < 1 {
		return nil, `{"error":"Not found!"}`, 404
	}

	var (
		err  error
		opts = docker.LogsOptions{
			Container:    args[0],
			Stdout:       true,
			Stderr:       true,
			Timestamps:   true,
			OutputStream: w,
			ErrorStream:  w,
		}
	)

	if _, ok := r.URL.Query()["follow"]; ok {
		opts.Follow = true
	}

	if err = d.dclient.Logs(opts); err != nil {
		d.logger.Error.Printf("Error getting log (%s): %s\n", args[0], err)
	}
	return nil, nil, -1
}

func SetupLogHandler(cfg *config.AppConfig, dockerUri string, logger *logging.Logger) {
	// Log handler
	logHdlr, err := NewLogHandler(dockerUri, logger)
	if err != nil {
		logger.Error.Fatalf("%s\n", err)
	}
	NewRestHandler(cfg.Endpoints.Logs, logHdlr, logger)
	logger.Warning.Printf("Logs API: %s\n", cfg.Endpoints.Logs)
}
