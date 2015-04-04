package notifications

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
	"github.com/naveabe/pkgwrap/pkgwrap/tracker"
)

func IrcNotify(ircs []string, msg string, logger *logging.Logger) {
	var (
		err error
		nt  *IRCNotifier
	)

	for _, v := range ircs {
		if nt, err = NewIRCNotifierFromString(v); err != nil {
			logger.Error.Printf("IRC notification request: %s\n", err)
			continue
		}
		if err = nt.Notify(msg); err != nil {
			logger.Error.Printf("Failed to send IRC notification: %s\n", err)
		}
	}
}

func EmailNotify(recipients []string, subject, msg string, logger *logging.Logger) {
	var err error
	for _, r := range recipients {
		nen := NewEmailNotifier(r)
		nen.Subject = subject
		nen.Body = msg

		if err = nen.Notify(); err != nil {
			logger.Error.Printf("Failed to send email notification: %s\n", err)
		}
	}
}

type NotificationProcessor struct {
	dstore   *tracker.TrackerStore
	Listener chan *docker.APIEvents

	logger *logging.Logger
}

func NewNotificationProcessor(dstore *tracker.TrackerStore, logger *logging.Logger) *NotificationProcessor {
	return &NotificationProcessor{
		dstore:   dstore,
		logger:   logger,
		Listener: make(chan *docker.APIEvents),
	}
}

/*
	IN PROGRESS
*/
func (n *NotificationProcessor) Start() {
	var (
		msg     string
		subject string
		status  string
	)
	for {
		dEvt := <-n.Listener

		// container info
		cont, err := n.dstore.GetContainer(dEvt.ID)
		if err != nil {
			n.logger.Error.Printf("%s", err)
			continue
		}
		// original request
		preq, err := n.dstore.GetRequestByContainerId(dEvt.ID)
		if err != nil {
			n.logger.Error.Printf("%s", err)
			continue
		}
		// distro
		distro, err := preq.GetDistribution(dEvt.ID)
		if err != nil {
			n.logger.Error.Printf("Distro - %s", err)
			continue
		}

		if preq.Notifications == nil {
			n.logger.Trace.Printf("No notifications specified!\n")
			continue
		}

		n.logger.Trace.Printf("Notifying: %s", preq.Notifications)

		if cont.State.ExitCode == 0 {
			status = "Succeeded"
		} else {
			status = "Failed"
		}

		msg = GetNotificationMessage(preq.Name, preq.Version, status, string(distro.Name), distro.Release)
		go IrcNotify(preq.Notifications.IRC, msg, n.logger)

		subject = GetNotificationSubject(preq.Name, preq.Version, string(distro.Name), distro.Release, status)
		go EmailNotify(preq.Notifications.Email, subject, msg, n.logger)
	}
}
