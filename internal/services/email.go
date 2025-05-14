package services

import (
	"fmt"
	"net/smtp"
)

func SendVerificationCodeViaEmail(to string, code string) error {
	smtpHost := ""
	smtpPort := ""

	username := ""
	password := ""
	from := ""

	auth := smtp.PlainAuth("", username, password, smtpHost)

	subject := "Your login code"
	body := fmt.Sprintf("Your code is: %s", code)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\r\n",
		from, to, subject, body)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
