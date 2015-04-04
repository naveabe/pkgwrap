package notifications

import (
	"fmt"
	"net/smtp"
)

const (
	DEFAULT_MAIL_HOST = "localhost"
	DEFAULT_MAIL_PORT = 25
	SUBJECT_TMPL      = "%s-%s (%s %s) Status: %s"
	FROM_ADDR_DEFAULT = "packager@ipkg.io"
	NOTIFICATION_TMPL = `
%s-%s [ Status: %s ]
Distro: %s Version: %s
`
)

func GetNotificationSubject(name, version, distro, release, status string) string {
	return fmt.Sprintf(SUBJECT_TMPL, name, version, distro, release, status)
}

func GetNotificationMessage(pkgName, pkgVersion, status, distro, distroVersion string) string {
	return fmt.Sprintf(NOTIFICATION_TMPL,
		pkgName, pkgVersion, status, distro, distroVersion)
}

type EmailNotifier struct {
	To      string
	From    string
	Subject string
	Body    string
}

func NewEmailNotifier(toAddr string) *EmailNotifier {
	return &EmailNotifier{
		To:   toAddr,
		From: FROM_ADDR_DEFAULT,
	}
}

func (e *EmailNotifier) Notify() error {
	// Set up authentication information.
	auth := smtp.PlainAuth("", e.From, "", DEFAULT_MAIL_HOST)
	if err := smtp.SendMail(fmt.Sprintf("%s:%d", DEFAULT_MAIL_HOST, DEFAULT_MAIL_PORT),
		auth, e.From, []string{e.To}, []byte(e.Body)); err != nil {

		return err
	}

	return nil
}
