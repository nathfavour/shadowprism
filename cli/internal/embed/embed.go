package embed

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed core_bin
var coreBinary []byte

// ExtractCore extracts the embedded Rust binary to a temporary location
// and returns the path to the executable.
func ExtractCore() (string, error) {
	tempDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	appDir := filepath.Join(tempDir, ".shadowprism", "bin")
	if err := os.MkdirAll(appDir, 0700); err != nil {
		return "", err
	}

	binPath := filepath.Join(appDir, "shadowprism-core")

	// We overwrite the binary every time to ensure it's up to date with the CLI version
	err = os.WriteFile(binPath, coreBinary, 0700)
	if err != nil {
		return "", fmt.Errorf("failed to write embedded binary: %w", err)
	}

	return binPath, nil
}
