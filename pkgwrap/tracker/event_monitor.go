package tracker

import (
	"encoding/json"
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

	datastore *EssDatastore
}

func NewDockerEventMonitor(url string, dstore *EssDatastore, logger *logging.Logger) *DockerEventMonitor {
	return &DockerEventMonitor{
		URL: url, logger: logger, datastore: dstore}
}

func (d *DockerEventMonitor) Start() error {

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
			d.logger.Error.Printf("%s\n", err)
			continue
		}
		d.logger.Trace.Printf("%s - %s\n", event.Status, event.Id)
		if event.Status == "die" || event.Status == "kill" {
			bldJob, err := d.datastore.GetBuild(event.Id)
			if err != nil {
				if err.Error() == "Not found" {
					d.logger.Warning.Printf("Build job not found: %s", event.Id)
					continue
				}
				d.logger.Error.Printf("%s\n", err)
				continue
			}
			for i, j := range bldJob.Jobs {
				if j.Id == event.Id {
					bldJob.Jobs[i].Status = event.Status
					break
				}
			}
			//d.datastore.Update(bldJob)
			d.logger.Debug.Printf("Updating: %s\n", bldJob)
		} else {
			d.logger.Trace.Printf("Skipping event: %s\n", event.Status)
			continue
		}

		// Process 'die' status

	}

	return nil
}
