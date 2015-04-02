package notifications

import (
	"fmt"
	"github.com/thoj/go-ircevent"
	"strings"
)

const (
	IRC_BOT_NAME         = "ipkg-io"
	IRC_DEFAULT_TLS_PORT = 6697
	IRC_DEFAULT_PORT     = 6667
)

type IRCNotifier struct {
	Host    string
	Channel string
}

func NewIRCNotifierFromString(irc string) (*IRCNotifier, error) {
	var (
		ircN = &IRCNotifier{}
		idx  = strings.Index(irc, "#")
	)

	if idx < 0 {
		return ircN, fmt.Errorf("Invalid IRC string: %s", (irc))
	}

	if strings.Index(irc[:idx], ":") < 0 {
		ircN.Host = fmt.Sprintf("%s:%d", irc[:idx], IRC_DEFAULT_PORT)
	} else {
		ircN.Host = irc[:idx]
	}
	ircN.Channel = irc[idx:]

	return ircN, nil
}

func (n *IRCNotifier) Notify(msgs string) error {
	con := irc.IRC(IRC_BOT_NAME, IRC_BOT_NAME)
	//con.UseTLS = true

	/* on connect */
	con.AddCallback("001", func(e *irc.Event) {
		con.Join(n.Channel)
	})
	/* on join */
	con.AddCallback("JOIN", func(e *irc.Event) {
		// Send message line by line.
		for _, m := range strings.Split(msgs, "\n") {
			con.Privmsg(n.Channel, m)
		}
		con.Quit()
	})

	err := con.Connect(n.Host)
	if err != nil {
		return err
	}

	con.Loop()

	return nil
}
