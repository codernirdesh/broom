//go:build windows

package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var (
	logFile *os.File
	mu      sync.Mutex
)

func Init() {
	f, err := os.OpenFile("cleanup_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Failed to create log file:", err)
	}
	logFile = f
	log.SetOutput(logFile)
	log.Println("=======================================")
	log.Println("Cleanup started at:", time.Now())
	log.Println("=======================================")
}

func Close() {
	if logFile != nil {
		logFile.Close()
	}
}

func Step(msg string) {
	mu.Lock()
	defer mu.Unlock()
	log.Println(msg)
	fmt.Println(msg)
}

func Banner(title string) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Println("=====================================")
	fmt.Println(" ", title)
	fmt.Println("=====================================")
	log.Println(title)
}
