//go:build windows

package main

import (
	"fmt"
	"os"
	"time"

	"broom/internal/cleanup"
	"broom/internal/elevate"
	"broom/internal/install"
	"broom/internal/logger"
	"broom/internal/update"
)

func main() {
	// Handle "broom update" before anything else — no admin needed
	if len(os.Args) > 1 && os.Args[1] == "update" {
		if err := update.Run(); err != nil {
			fmt.Println("[ERROR]", err)
			os.Exit(1)
		}
		return
	}

	logger.Init()
	defer logger.Close()

	logger.Banner("SYSTEM CLEANUP STARTED")

	if !elevate.IsAdmin() {
		fmt.Println("[SECURITY] Not running as Administrator. Requesting elevation...")
		elevate.RunAsAdmin()
		return // RunAsAdmin exits the process; safety fallback
	}

	install.EnsureInPath()

	cleanup.Run()

	logger.Banner("CLEANUP COMPLETED SUCCESSFULLY")
	time.Sleep(3 * time.Second)
}
