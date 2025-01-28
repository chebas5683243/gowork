package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func GetProjectRoot() (string, error) {
	const PROJECT_MARKER = "go.mod"

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get the current file path")
	}

	dir := filepath.Dir(currentFile)
	for {
		markerPath := filepath.Join(dir, PROJECT_MARKER)
		if _, err := os.Stat(markerPath); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("project root not found (marker: %s)", PROJECT_MARKER)
}
