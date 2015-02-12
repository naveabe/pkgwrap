package tracker

import (
	"encoding/json"
	"github.com/fsouza/go-dockerclient"
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
	"io"
	"net/http"
)

type DockerEvent struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

type DockerEventMonitor struct {
	URL    string
	logger *logging.Logger

	datastore *EssJobstore

	client *docker.Client
}

func NewDockerEventMonitor(eventUrl, dockerUri string, dstore *EssJobstore, logger *logging.Logger) (*DockerEventMonitor, error) {
	var (
		dem = DockerEventMonitor{
			URL:       eventUrl,
			logger:    logger,
			datastore: dstore,
		}
		err error
	)
	if dem.client, err = docker.NewClient(dockerUri); err != nil {
		return &dem, err
	}
	return &dem, nil
}

func (d *DockerEventMonitor) Start() error {
	var (
		bldJob *BuildJob
		status string
	)

	resp, err := http.Get(d.URL + "/events")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	for {
		var event DockerEvent
		if err := dec.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}
			d.logger.Error.Printf("Decode event failed: %s\n", err)
			continue
		}
		d.logger.Trace.Printf("%s - %s\n", event.Status, event.Id)
		// Get container info
		dCont, err := d.client.InspectContainer(event.Id)
		if err != nil {
			d.logger.Error.Printf("Failed to get info (%s): %s\n", event.Id, err)
			continue
		}

		status = event.Status
		switch status {
		case "die":
			bldJob, err = d.datastore.GetBuild(event.Id)
			break
		case "kill":
			bldJob, err = d.datastore.GetBuild(event.Id)
			break
		default:
			d.logger.Trace.Printf("Skipping event: %s\n", event.Status)
			continue
		}

		if err != nil {
			d.logger.Error.Printf("%s\n", err)
			continue
		}

		// Set status based on actual build
		if dCont.State.ExitCode != 0 {
			status = "failed"
		} else {
			status = "succeeded"
		}
		// Set new status for build job
		if !bldJob.SetJobStatus(event.Id, status) {
			d.logger.Error.Printf("Could not update job status: %s\n", event.Id)
			continue
		}
		// Update status
		if err = d.datastore.Update("job", bldJob.Id, bldJob); err != nil {
			d.logger.Error.Printf("Failed to update job: %s\n", bldJob)
			continue
		}
		d.logger.Debug.Printf("Updated build job: %s\n", bldJob)
	}

	return nil
}
