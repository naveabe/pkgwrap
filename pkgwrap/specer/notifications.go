package specer

type BuildNotifications struct {
	// e.g. chat.freenode.net#ipkgio
	IRC   []string `yaml:"irc"`
	Email []string `yaml:"email"`
}
