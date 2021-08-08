package email

import (
	"fmt"
	"net/smtp"
)

type Emailer interface {
	Send(to []string, msg string) error
}

type SimpleEmailer struct {
	From, Password, Host, Port string
}

func (e SimpleEmailer) Send(toList []string, body string) error {
	auth := smtp.PlainAuth("", e.From, e.Password, e.Host)
	from := fmt.Sprintf("From: <%s>\r\n", e.From)
	to := fmt.Sprintf("To: <%s>\r\n", toList[0])
	subject := "Subject: Basic Functioning Passwordless Stuff\r\n"
	b := fmt.Sprintf("%s\r\n<3 Swanky McSpanky", body)
	msg := from + to + subject + "\r\n" + b
	return smtp.SendMail(e.Host+":"+e.Port, auth, e.From, toList, []byte(msg))
}
