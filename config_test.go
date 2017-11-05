package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type testCaseConfig struct {
	name     string
	files    map[string]string
	expected map[string]string
}

func TestConfig(t *testing.T) {
	tests := []func(t *testing.T){
		testLoadConfig,
	}

	for _, test := range tests {
		setupConfigTest(t)
		test(t)
		cleanupConfigTest(t)
	}
}

func setupConfigTest(t *testing.T) {
	oldFiles, err := filepath.Glob(".env*")
	if err != nil {
		t.Fatalf("Unable to glob for .env* files: %v\n", err)
	}

	for _, oldFile := range oldFiles {
		t.Logf("renaming %s to %s", oldFile, ".tmp"+oldFile)
		os.Rename(oldFile, ".tmp"+oldFile)
	}
}

func cleanupConfigTest(t *testing.T) {
	createdFiles, err := filepath.Glob(".env*")
	if err == nil {
		for _, createdFile := range createdFiles {
			os.Remove(createdFile)
		}
	}

	tmpFiles, err := filepath.Glob(".tmp.env*")
	if err != nil {
		t.Fatalf("Unable to glob for .tmp.env* files: %v\n", err)
	}

	for _, tmpFile := range tmpFiles {
		t.Logf("restoring %s to %s", tmpFile, strings.TrimPrefix(tmpFile, ".tmp"))
		os.Rename(tmpFile, strings.TrimPrefix(tmpFile, ".tmp"))
	}
}

func runTests(t *testing.T, funcName string, testCasesConfig []testCaseConfig) {
	for _, tc := range testCasesConfig {
		for filename, contents := range tc.files {
			f, err := os.Create(filename)
			if err != nil {
				t.Fatalf("Cannot create file %s: %v\n", filename, err)
			}

			f.WriteString(contents)

			f.Close()
		}

		loadConfig()

		for key, value := range tc.expected {
			if os.Getenv(key) != value {
				t.Errorf("%s: %s: expected %s=%s, got %s=%s", funcName, tc.name, key, value, key, os.Getenv(key))
			}
			os.Setenv(key, "")
		}
	}
}

func testLoadConfig(t *testing.T) {
	testCasesConfig := []testCaseConfig{
		testCaseConfig{
			"Absence of all .env* files should load defaults",
			map[string]string{},
			map[string]string{
				"COMMENTO_PORT": "8080",
			},
		},

		testCaseConfig{
			".env should be loaded",
			map[string]string{
				".env": `
					COMMENTO_PORT=8081
					env1=val1
					env2=val2`,
			},
			map[string]string{
				"COMMENTO_PORT": "8081",
				"env1":          "val1",
				"env2":          "val2",
			},
		},

		testCaseConfig{
			".env should dominate .env.test",
			map[string]string{
				".env.test": `
					env1=val1
					env2=val2`,
				".env": `
					env2=val3
					env3=val4`,
			},
			map[string]string{
				"env1": "val1",
				"env2": "val3",
				"env3": "val4",
			},
		},
	}

	runTests(t, "loadConfig", testCasesConfig)
}
