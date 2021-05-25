package mailUtils

import (
	"crypto/tls"
	"grpc-demo/utils"
	"github.com/go-gomail/gomail"
)

func Send(token string, target ...string) error {

	m := gomail.NewMessage()
	m.SetHeader("From", utils.GlobalConfig.Mail.Username)
	m.SetHeader("To", target...)

	m.SetHeader("Subject", "捣蒜官方邮件")
	m.SetBody("text/html", generateBody(token))
	d := gomail.NewDialer(
		utils.GlobalConfig.Mail.SmtpHost, 465,
		utils.GlobalConfig.Mail.Username, utils.GlobalConfig.Mail.Password)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
