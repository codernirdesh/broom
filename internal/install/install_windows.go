//go:build windows

package install

import (
	"os"
	"path/filepath"
	"strings"

	"broom/internal/logger"

	"golang.org/x/sys/windows/registry"
)

// EnsureInPath copies the running binary into %LOCALAPPDATA%\broom and adds
// that directory to the user's PATH if not already present. Runs once silently.
func EnsureInPath() {
	installDir := filepath.Join(os.Getenv("LOCALAPPDATA"), "broom")
	destExe := filepath.Join(installDir, "broom.exe")

	srcExe, err := os.Executable()
	if err != nil {
		logger.Step("[INSTALL] Could not resolve executable path: " + err.Error())
		return
	}

	// Skip if already running from the install directory
	if strings.EqualFold(filepath.Dir(srcExe), installDir) {
		return
	}

	// Create install directory
	if err := os.MkdirAll(installDir, 0755); err != nil {
		logger.Step("[INSTALL] Could not create install dir: " + err.Error())
		return
	}

	// Copy binary
	data, err := os.ReadFile(srcExe)
	if err != nil {
		logger.Step("[INSTALL] Could not read binary: " + err.Error())
		return
	}
	if err := os.WriteFile(destExe, data, 0755); err != nil {
		logger.Step("[INSTALL] Could not write binary: " + err.Error())
		return
	}

	// Add to user PATH via registry
	key, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		logger.Step("[INSTALL] Could not open registry: " + err.Error())
		return
	}
	defer key.Close()

	currentPath, _, err := key.GetStringValue("Path")
	if err != nil && err != registry.ErrNotExist {
		logger.Step("[INSTALL] Could not read PATH: " + err.Error())
		return
	}

	if !containsPath(currentPath, installDir) {
		newPath := currentPath
		if newPath != "" && !strings.HasSuffix(newPath, ";") {
			newPath += ";"
		}
		newPath += installDir

		if err := key.SetStringValue("Path", newPath); err != nil {
			logger.Step("[INSTALL] Could not update PATH: " + err.Error())
			return
		}
		logger.Step("[INSTALL] Added broom to PATH. Restart your terminal to use 'broom' from anywhere.")
	}
}

func containsPath(pathEnv, dir string) bool {
	for _, p := range strings.Split(pathEnv, ";") {
		if strings.EqualFold(strings.TrimSpace(p), dir) {
			return true
		}
	}
	return false
}
