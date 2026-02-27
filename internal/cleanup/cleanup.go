//go:build windows

package cleanup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"broom/internal/logger"
)

func Run() {
	var wg sync.WaitGroup

	// TEMP
	wg.Add(1)
	go func() {
		defer wg.Done()
		section("TEMP FILES")
		removeSafe("C:\\Windows\\Temp")
		removeSafe(os.TempDir())
	}()

	// PREFETCH
	wg.Add(1)
	go func() {
		defer wg.Done()
		section("PREFETCH")
		removeSafe("C:\\Windows\\Prefetch")
	}()

	// LOGS
	wg.Add(1)
	go func() {
		defer wg.Done()
		section("LOGS & REPORTS")
		removeSafe("C:\\Windows\\Logs")
		removeSafe("C:\\ProgramData\\Microsoft\\Windows\\WER\\ReportArchive")
		removeSafe("C:\\ProgramData\\Microsoft\\Windows\\WER\\ReportQueue")
	}()

	// THUMBNAIL CACHE
	wg.Add(1)
	go func() {
		defer wg.Done()
		section("THUMBNAIL CACHE")
		runCmd("taskkill", "/f", "/im", "explorer.exe")
		removeSafe(os.Getenv("LOCALAPPDATA") + "\\Microsoft\\Windows\\Explorer")
		startCmd("explorer.exe")
	}()

	// DELIVERY OPT
	wg.Add(1)
	go func() {
		defer wg.Done()
		section("DELIVERY OPTIMIZATION")
		removeSafe("C:\\Windows\\SoftwareDistribution\\DeliveryOptimization")
	}()

	// RECYCLE BIN
	wg.Add(1)
	go func() {
		defer wg.Done()
		section("RECYCLE BIN")
		removeSafe(os.Getenv("SystemDrive") + "\\$Recycle.Bin")
	}()

	// DNS
	wg.Add(1)
	go func() {
		defer wg.Done()
		section("DNS CACHE")
		runCmd("ipconfig", "/flushdns")
	}()

	// WINDOWS UPDATE CACHE
	wg.Add(1)
	go func() {
		defer wg.Done()
		section("WINDOWS UPDATE CACHE")
		runCmd("net", "stop", "wuauserv")
		runCmd("net", "stop", "bits")
		removeSafe("C:\\Windows\\SoftwareDistribution\\Download")
		runCmd("net", "start", "wuauserv")
		runCmd("net", "start", "bits")
	}()

	// WINDOWS OLD
	wg.Add(1)
	go func() {
		defer wg.Done()
		if exists("C:\\Windows.old") {
			section("WINDOWS.OLD REMOVAL")
			runCmd("takeown", "/F", "C:\\Windows.old", "/R", "/D", "Y")
			runCmd("icacls", "C:\\Windows.old", "/grant", os.Getenv("USERNAME")+":F", "/T")
			os.RemoveAll("C:\\Windows.old")
			logger.Step("[OK] Windows.old removed")
		}
	}()

	// MICROSOFT STORE RESET
	wg.Add(1)
	go func() {
		defer wg.Done()
		section("MICROSOFT STORE RESET")
		storeCache := filepath.Join(os.Getenv("LOCALAPPDATA"), "Packages", "Microsoft.WindowsStore_8wekyb3d8bbwe", "LocalCache")
		removeSafe(storeCache)
	}()

	// DISM CLEANUP
	wg.Add(1)
	go func() {
		defer wg.Done()
		section("DISM COMPONENT CLEANUP")
		runCmd("dism", "/online", "/cleanup-image", "/startcomponentcleanup")
	}()

	wg.Wait()
}

func section(name string) {
	logger.Step(fmt.Sprintf("\n========== %s ==========", name))
}

func removeSafe(path string) {
	if !exists(path) {
		logger.Step("[SKIPPED] Not found: " + path)
		return
	}

	logger.Step("[CLEANING] " + path)

	entries, err := os.ReadDir(path)
	if err != nil {
		logger.Step("[ERROR] Cannot read: " + path)
		return
	}

	var wg sync.WaitGroup
	var failedCount int
	var mu sync.Mutex

	for _, e := range entries {
		wg.Add(1)
		go func(e os.DirEntry) {
			defer wg.Done()
			fullPath := filepath.Join(path, e.Name())
			err := os.RemoveAll(fullPath)
			if err != nil {
				mu.Lock()
				failedCount++
				mu.Unlock()
			}
		}(e)
	}
	wg.Wait()

	if failedCount > 0 {
		logger.Step(fmt.Sprintf("[OK] Cleaned: %s (Skipped %d files in use)", path, failedCount))
	} else {
		logger.Step("[OK] Cleaned: " + path)
	}
}

func runCmd(name string, args ...string) {
	cmd := exec.Command(name, args...)
	err := cmd.Run()
	if err != nil {
		logger.Step(fmt.Sprintf("[WARN] Command failed: %s %v", name, args))
	}
}

func startCmd(name string, args ...string) {
	cmd := exec.Command(name, args...)
	err := cmd.Start()
	if err != nil {
		logger.Step(fmt.Sprintf("[WARN] Command start failed: %s %v", name, args))
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
