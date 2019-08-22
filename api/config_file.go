package main

import (
	"bufio"
	"os"
	"strings"
)

func configFileLoad(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}

	defer file.Close()

	num := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		num += 1

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		i := strings.Index(line, "=")
		if i == -1 {
			logger.Errorf("%s: line %d: neither a comment nor a valid setting", filepath, num)
			return errorInvalidConfigFile
		}

		key := line[:i]
		value := line[i+1:]

		if !strings.HasPrefix(key, "COMMENTO_") {
			continue
		}

		if os.Getenv(key) != "" {
			// Config files have lower precedence.
			continue
		}

		os.Setenv(key, value)
	}

	return nil
}
