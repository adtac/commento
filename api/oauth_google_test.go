package main

import (
	"os"
	"testing"
)

func resetGoogleVars() {
	for _, env := range []string{"GOOGLE_KEY", "GOOGLE_SECRET"} {
		os.Setenv(env, "")
	}
}

func TestGoogleOauthConfigureBasics(t *testing.T) {
	resetGoogleVars()

	os.Setenv("GOOGLE_KEY", "google-key")
	os.Setenv("GOOGLE_SECRET", "google-secret")

	if err := googleOauthConfigure(); err != nil {
		t.Errorf("unexpected error configuring google oauth: %v", err)
		return
	}

	if googleConfig == nil {
		t.Errorf("expected googleConfig!=nil got googleConfig=nil")
		return
	}
}

func TestGoogleOauthConfigureEmpty(t *testing.T) {
	resetGoogleVars()

	os.Setenv("GOOGLE_KEY", "google-key")

	if err := googleOauthConfigure(); err == nil {
		t.Errorf("expected error not found when configuring google oauth with empty COMMENTO_GOOGLE_SECRET")
		return
	}

	if googleConfig != nil {
		t.Errorf("expected googleConfig=nil got googleConfig=%v", googleConfig)
		return
	}
}

func TestGoogleOauthConfigureEmpty2(t *testing.T) {
	resetGoogleVars()

	if err := googleOauthConfigure(); err != nil {
		t.Errorf("unexpected error configuring google oauth with empty everything: should be disabled")
		return
	}

	if googleConfig != nil {
		t.Errorf("expected googleConfig=nil got googleConfig=%v", googleConfig)
		return
	}
}
