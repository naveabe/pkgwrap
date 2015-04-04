package notifications

import (
	"bytes"
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
	c, err := smtp.Dial(fmt.Sprintf("%s:%d", DEFAULT_MAIL_HOST, DEFAULT_MAIL_PORT))
	if err != nil {
		return err
	}
	// Set the sender and recipient.
	c.Mail(e.From)
	c.Rcpt(e.To)

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		return err
	}
	defer wc.Close()

	// write subject
	subject := bytes.NewBufferString(fmt.Sprintf("Subject: %s\r\n\r\n", e.Subject))
	if _, err = subject.WriteTo(wc); err != nil {
		return err
	}

	buf := bytes.NewBufferString(e.Body)
	if _, err = buf.WriteTo(wc); err != nil {
		return err
	}

	// Set up authentication information.
	//auth := smtp.PlainAuth("", e.From, "", DEFAULT_MAIL_HOST)
	//auth := smtp.CRAMMD5Auth("", "")
	//if err := smtp.SendMail(,
	//	auth, e.From, []string{e.To}, []byte(e.Body)); err != nil {

	//	return err
	//}

	return nil
}
