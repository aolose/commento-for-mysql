package main

import (
	"bytes"
	"os"
)

type ownerResetHexPlugs struct {
	Origin   string
	ResetHex string
}

func smtpOwnerResetHex(to string, toName string, resetHex string) error {
	var header bytes.Buffer
	headerTemplate.Execute(&header, &headerPlugs{FromAddress: os.Getenv("SMTP_FROM_ADDRESS"), ToAddress: to, ToName: toName, Subject: "Reset your password"})

	var body bytes.Buffer
	templates["reset-hex"].Execute(&body, &ownerResetHexPlugs{Origin: os.Getenv("ORIGIN"), ResetHex: resetHex})

	err := sendMail([]string{to}, concat(header, body))
	if err != nil {
		logger.Errorf("cannot send reset email: %v", err)
		return errorCannotSendEmail
	}

	return nil
}
