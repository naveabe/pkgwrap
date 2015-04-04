package tracker

import (
	"encoding/json"
	"github.com/fsouza/go-dockerclient"
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
)

type DockerEventMonitor struct {
	logger *logging.Logger
	// elasticsearch datastore
	datastore *TrackerStore
	// docker client
	client *docker.Client
	// notification channel for post build tasks.
	notifications chan *docker.APIEvents
}

func NewDockerEventMonitor(dockerUri string, dstore *TrackerStore,
	nChan chan *docker.APIEvents, logger *logging.Logger) (*DockerEventMonitor, error) {

	var (
		dem = DockerEventMonitor{
			logger:        logger,
			datastore:     dstore,
			notifications: nChan,
		}
		err error
	)
	if dem.client, err = docker.NewClient(dockerUri); err != nil {
		return &dem, err
	}
	return &dem, nil
}

/*
	Start listening to events and update container information accordingly.
*/
func (d *DockerEventMonitor) Start() error {
	var (
		err      error
		dCont    *docker.Container
		listener = make(chan *docker.APIEvents)
	)

	if err = d.client.AddEventListener(listener); err != nil {
		return err
	}

	for {
		event := <-listener
		//d.logger.Trace.Printf("Event: %s - %s\n", event.Status, event.ID)
		// Only update datastore on these status'
		switch event.Status {
		case "create":
			break
		case "start":
			break
		case "die":
			// TODO: copy logs (useful when docker gets cleaned up)
			break
		case "kill":
			// TODO: copy logs
			break
		default:
			d.logger.Trace.Printf("Skipping event: %s\n", event.Status)
			continue
		}

		// Get container info
		if dCont, err = d.client.InspectContainer(event.ID); err != nil {
			d.logger.Error.Printf("Failed to get container info (%s): %s\n", event.ID, err)
			continue
		}
		//d.logger.Trace.Printf("inspect container: %#v\n", dCont.State)
		b, _ := json.MarshalIndent(dCont.State, "", "  ")
		d.logger.Trace.Printf("inspect container: %s\n", b)

		// Update datastore
		if err = d.datastore.UpdateContainer(event.ID, dCont); err != nil {
			d.logger.Error.Printf("Failed to update container info (%s): %s\n", event.ID, err)
			continue
		}
		d.logger.Debug.Printf("Updated container (%s): %s\n", event.Status, event.ID)

		// Only send on die and kill events.
		if d.notifications != nil {
			if event.Status == "die" || event.Status == "kill" {
				d.notifications <- event
			}
		}
	}

	return nil
}

/*
	Params:
		dockerUri : docker uri
		dstore : datastore struct
		notifChan: channel so send notifications on
		logger : global logger
*/
func StartEventMonitor(dockerUri string, dstore *TrackerStore, notifChan chan *docker.APIEvents, logger *logging.Logger) {
	if dem, err := NewDockerEventMonitor(dockerUri, dstore, notifChan, logger); err == nil {
		if err := dem.Start(); err != nil {
			logger.Error.Fatalf("%s\n", err)
		}
	} else {
		logger.Error.Fatalf("%s\n", err)
	}
}
