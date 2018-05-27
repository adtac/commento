package main

import (
	"testing"
)

func TestCreateLoggerBasics(t *testing.T) {
	logger = nil

	if err := createLogger(); err != nil {
		t.Errorf("unexpected error creating logger: %v", err)
		return
	}

	if logger == nil {
		t.Errorf("logger null after createLogger()")
		return
	}

	logger.Debugf("test message please ignore")
}
