//go:build windows

package main

import (
	"fmt"
	"time"

	"broom/internal/cleanup"
	"broom/internal/elevate"
	"broom/internal/logger"
)

func main() {
	logger.Init()
	defer logger.Close()

	logger.Banner("SYSTEM CLEANUP STARTED")

	if !elevate.IsAdmin() {
		fmt.Println("[SECURITY] Not running as Administrator. Requesting elevation...")
		elevate.RunAsAdmin()
		return // RunAsAdmin exits the process; safety fallback
	}

	cleanup.Run()

	logger.Banner("CLEANUP COMPLETED SUCCESSFULLY")
	time.Sleep(3 * time.Second)
}
