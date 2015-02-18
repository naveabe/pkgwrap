package tracker

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
)

type DockerEventMonitor struct {
	logger *logging.Logger

	datastore *EssJobstore

	client *docker.Client
}

func NewDockerEventMonitor(dockerUri string, dstore *EssJobstore, logger *logging.Logger) (*DockerEventMonitor, error) {
	var (
		dem = DockerEventMonitor{
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
		bldJob   *BuildJob
		status   string
		err      error
		listener = make(chan *docker.APIEvents)
	)

	if err = d.client.AddEventListener(listener); err != nil {
		return err
	}

	for {
		event := <-listener

		d.logger.Trace.Printf("Event: %s - %s\n", event.Status, event.ID)

		status = event.Status
		switch status {
		case "die":
			bldJob, err = d.datastore.GetBuild(event.ID)
			break
		case "kill":
			bldJob, err = d.datastore.GetBuild(event.ID)
			break
		default:
			d.logger.Trace.Printf("Skipping event: %s\n", event.Status)
			continue
		}

		if err != nil {
			d.logger.Error.Printf("%s\n", err)
			continue
		}

		// Get container info
		dCont, err := d.client.InspectContainer(event.ID)
		if err != nil {
			d.logger.Error.Printf("Failed to get info (%s): %s\n", event.ID, err)
			continue
		}
		// Set status based on actual build
		if dCont.State.ExitCode != 0 {
			status = "failed"
		} else {
			status = "succeeded"
		}
		// Set new status for build job
		if !bldJob.SetJobStatus(event.ID, status) {
			d.logger.Error.Printf("Could not update job status: %s\n", event.ID)
			continue
		}
		// Update status
		if err = d.datastore.Update("job", bldJob.Id, bldJob); err != nil {
			d.logger.Error.Printf("Failed to update job: %s\n", bldJob)
			continue
		}
		d.logger.Debug.Printf("Updated build job: %s\n", bldJob)
		// TODO: event
		// TODO: add end timestamp
	}

	return nil
}
