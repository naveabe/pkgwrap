package notifications

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
	"github.com/naveabe/pkgwrap/pkgwrap/tracker"
)

func EmailNotify(msg string) {

}

func IrcNotify(msg string) {
	/*
	   if b.IRC != nil {

	       var (
	           nt  *IRCNotifier
	           err error
	       )

	       for _, v := range b.IRC {
	           if nt, err = NewIRCNotifierFromString(v); err != nil {
	               fmt.Println(err)
	               continue
	           }
	           if err = nt.Notify(msgs); err != nil {
	               fmt.Println(err)
	           }
	       }
	   }
	*/
}

func buildNotification() {
	// message to send
}

func Notify(msg string) {
	//if != nil
	go EmailNotify(msg)
	go IrcNotify(msg)
}

type NotificationProcessor struct {
	dstore     *tracker.TrackerStore
	listenChan chan *docker.APIEvents

	logger *logging.Logger
}

func (n *NotificationProcessor) Start() {
	for {
		dEvt := <-n.listenChan

		preq, err := n.dstore.GetRequestByContainerId(dEvt.ID)
		if err != nil {
			n.logger.Error.Printf("%s", err)
			continue
		}
		//
		if preq.Notifications == nil {
			n.logger.Trace.Printf("No notifications specified!\n")
			continue
		}

		// Notify(message)
	}
}
