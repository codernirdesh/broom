//go:build windows

package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const repoAPI = "https://api.github.com/repos/codernirdesh/broom/releases/latest"

type ghAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type ghRelease struct {
	TagName string    `json:"tag_name"`
	Assets  []ghAsset `json:"assets"`
}

// Run downloads the latest broom.exe from GitHub releases and replaces the
// current binary. The old binary is kept as broom.old.exe as a fallback.
func Run() error {
	fmt.Println("[UPDATE] Checking for latest release...")

	resp, err := http.Get(repoAPI)
	if err != nil {
		return fmt.Errorf("could not reach GitHub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub API returned %s", resp.Status)
	}

	var release ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return fmt.Errorf("could not parse release info: %w", err)
	}

	// Find the right asset — prefer broom.exe (GUI build)
	var downloadURL string
	for _, a := range release.Assets {
		if strings.EqualFold(a.Name, "broom.exe") {
			downloadURL = a.BrowserDownloadURL
			break
		}
	}
	if downloadURL == "" {
		return fmt.Errorf("no broom.exe found in release %s", release.TagName)
	}

	fmt.Printf("[UPDATE] Downloading %s ...\n", release.TagName)

	dlResp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer dlResp.Body.Close()

	if dlResp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned %s", dlResp.Status)
	}

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not resolve own path: %w", err)
	}
	exePath, _ = filepath.EvalSymlinks(exePath)

	// Write to a temp file in the same directory first
	tmpPath := exePath + ".tmp"
	out, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("could not create temp file: %w", err)
	}

	if _, err := io.Copy(out, dlResp.Body); err != nil {
		out.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("could not write update: %w", err)
	}
	out.Close()

	// Rename current binary to .old, then move new binary into place
	oldPath := exePath + ".old"
	os.Remove(oldPath) // clean up any previous .old file
	if err := os.Rename(exePath, oldPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("could not back up current binary: %w", err)
	}
	if err := os.Rename(tmpPath, exePath); err != nil {
		// Try to restore the old binary
		os.Rename(oldPath, exePath)
		return fmt.Errorf("could not replace binary: %w", err)
	}

	fmt.Printf("[UPDATE] Updated to %s. Restart broom to use the new version.\n", release.TagName)
	return nil
}
