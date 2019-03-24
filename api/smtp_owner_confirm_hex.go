package main

import (
	"bytes"
	"os"
)

type ownerConfirmHexPlugs struct {
	Origin     string
	ConfirmHex string
}

func smtpOwnerConfirmHex(to string, toName string, confirmHex string) error {
	var header bytes.Buffer
	headerTemplate.Execute(&header, &headerPlugs{FromAddress: os.Getenv("SMTP_FROM_ADDRESS"), ToAddress: to, ToName: toName, Subject: "Please confirm your email address"})

	var body bytes.Buffer
	templates["confirm-hex"].Execute(&body, &ownerConfirmHexPlugs{Origin: os.Getenv("ORIGIN"), ConfirmHex: confirmHex})

	err := sendMail([]string{to}, concat(header, body))
	if err != nil {
		logger.Errorf("cannot send confirmation email: %v", err)
		return errorCannotSendEmail
	}

	return nil
}
