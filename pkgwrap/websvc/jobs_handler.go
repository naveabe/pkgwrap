package websvc

import (
	"bufio"
	"fmt"
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
	"github.com/naveabe/pkgwrap/pkgwrap/tracker"
	"io"
	"net/http"
)

type JobsHandler struct {
	DefaultMethodHandler

	logger *logging.Logger

	datastore *tracker.EssDatastore
}

func NewJobsHandler(dstore *tracker.EssDatastore, logger *logging.Logger) *JobsHandler {
	lgh := JobsHandler{datastore: dstore}
	if logger != nil {
		lgh.logger = logger
	} else {
		lgh.logger = logging.NewStdLogger()
	}
	return &lgh
}

func (l *JobsHandler) proxyLogStream(w http.ResponseWriter, r *http.Request, id string) error {

	bldJob, err := l.datastore.GetBuild(id)
	if err != nil {
		return err
	}
	l.logger.Trace.Printf("Getting sub job: %s\n", id)
	job, err := bldJob.GetSubJob(id)
	if err != nil {
		return err
	}

	//logUrl := fmt.Sprintf("http://%s/containers/%s/logs?stderr=1&stdout=1&follow=1", job.Uri, job.Id)
	logUrl := fmt.Sprintf("http://%s/containers/%s/logs?stderr=1&stdout=1", job.Uri, job.Id)
	resp, err := http.Get(logUrl)
	if err != nil {
		return err
	} else if resp.StatusCode < 200 || resp.StatusCode > 304 {
		l.logger.Warning.Printf("Could not retrieve log: %s\n", logUrl)
		return fmt.Errorf("%s", resp.Status)
	}

	l.logger.Trace.Printf("Log response: %#v\n", resp)

	bRdr := bufio.NewReader(resp.Body)

	defer resp.Body.Close()

	l.logger.Trace.Printf("Tailing log: %s...\n", id)

	for {
		lineBytes, err := bRdr.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				w.WriteHeader(200)
				break
			}
			l.logger.Warning.Printf("%s\n", err)
			continue
		}

		_, err = w.Write(lineBytes)
		if err != nil {
			l.logger.Warning.Printf("%s\n", err)
			continue
		}
	}
	return nil
}

func (l *JobsHandler) GET(w http.ResponseWriter, r *http.Request, args ...string) (map[string]string, interface{}, int) {
	var (
		bJobs interface{}
		err   error
	)

	switch len(args) {
	case 1:
		bJobs, err = l.datastore.GetBuildsForUser(args[0])
		break
	case 2:
		bJobs, err = l.datastore.GetBuildsForPackage(args[0], args[1])
		break
	case 3:
		bJobs, err = l.datastore.GetBuildsForPackageVersion(args[0], args[1], args[2])
		break
	case 4:
		// container id
		bJobs, err = l.datastore.GetBuild(args[3])
		// Add more details.
		break
	case 5:
		if args[4] == "log" {
			l.logger.Trace.Printf("Streaming log %s...\n", args[3])
			if err = l.proxyLogStream(w, r, args[3]); err == nil {
				return nil, nil, -1
			}
		} else {
			return nil, map[string]string{"error": "Not found!"}, 404
		}
		break
	default:
		err = fmt.Errorf("Invalid request")
		break
	}

	if err != nil {
		return nil, map[string]string{"error": err.Error()}, 400
	}
	//l.logger.Trace.Printf("%v\n", args)
	return nil, bJobs, 200
}
