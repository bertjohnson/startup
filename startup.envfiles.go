package startup

import (
	"io/ioutil"
	"os"
	"strings"
)

// readEnvFiles reads local environment variable overrides.
func readEnvFiles() error {
	// Check if .env file exists.
	fi, err := os.Stat(".env")
	if err != nil {
		return nil
	}
	if fi.IsDir() {
		return nil
	}

	// Read .env.
	fileBytes, err := ioutil.ReadFile(".env")
	if err != nil {
		return err
	}

	// Parse each line.
	lines := strings.Split(string(fileBytes), "\n")
	for _, line := range lines {
		// Trim whitespace.
		line = strings.TrimSpace(line)

		// Ignore comments.
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Split parameters.
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			err = os.Setenv(strings.Trim(parts[0], " \t'\""), strings.Trim(parts[1], " \t'\""))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
