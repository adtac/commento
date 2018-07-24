package main

import (
	"os"
	"testing"
)

func smtpVarsClean() {
	for _, env := range []string{"SMTP_USERNAME", "SMTP_PASSWORD", "SMTP_HOST", "SMTP_PORT", "SMTP_FROM_ADDRESS"} {
		os.Setenv(env, "")
	}
}

func TestSmtpConfigureBasics(t *testing.T) {
	failTestOnError(t, setupTestEnv())
	smtpVarsClean()

	os.Setenv("SMTP_USERNAME", "test@example.com")
	os.Setenv("SMTP_PASSWORD", "hunter2")
	os.Setenv("SMTP_HOST", "smtp.commento.io")
	os.Setenv("SMTP_FROM_ADDRESS", "no-reply@commento.io")

	if err := smtpConfigure(); err != nil {
		t.Errorf("unexpected error when configuring SMTP: %v", err)
		return
	}
}

func TestSmtpConfigureEmptyHost(t *testing.T) {
	failTestOnError(t, setupTestEnv())
	smtpVarsClean()

	os.Setenv("SMTP_USERNAME", "test@example.com")
	os.Setenv("SMTP_PASSWORD", "hunter2")
	os.Setenv("SMTP_FROM_ADDRESS", "no-reply@commento.io")

	if err := smtpConfigure(); err != nil {
		t.Errorf("unexpected error when configuring SMTP: %v", err)
		return
	}

	if smtpConfigured {
		t.Errorf("SMTP configured when it should not be due to empty COMMENTO_SMTP_HOST")
		return
	}
}

func TestSmtpConfigureEmptyAddress(t *testing.T) {
	failTestOnError(t, setupTestEnv())
	smtpVarsClean()

	os.Setenv("SMTP_USERNAME", "test@example.com")
	os.Setenv("SMTP_PASSWORD", "hunter2")
	os.Setenv("SMTP_HOST", "smtp.commento.io")
	os.Setenv("SMTP_PORT", "25")

	if err := smtpConfigure(); err == nil {
		t.Errorf("expected error not found; SMTP should not be configured when COMMENTO_SMTP_FROM_ADDRESS is empty")
		return
	}
}
