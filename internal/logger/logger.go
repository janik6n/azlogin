package logger

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/janik6n/azlogin/internal/configuration"
)

func SetupLogger(c configuration.Configuration) error {
	// This is the native way, but no log rotation etc.
	// flags := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	// file, err := os.OpenFile("azlogin.log", flags, 0666)
	// if err != nil {
	// 	return err
	// }
	// Redirecting logs to the file
	// log.SetOutput(file)
	// log.Println("Hello from LoggerTest() function")

	if !c.General.Logging {
		return nil
	}
	// Get user's home directory
	userHomeDirectory, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user's home directory, %v", err)
	}

	logFilename := userHomeDirectory + "/azlogin/azlogin.log"
	if os.Getenv("ENVIRONMENT") == "DEV" {
		logFilename = "./azlogin.log"
	}

	log.SetOutput(&lumberjack.Logger{
		Filename:   logFilename,
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})

	return nil
}

func LogInfo(info string, funcName string, c configuration.Configuration) {
	if c.General.Logging && c.General.LoggingLevel == "INFO" {
		log.Printf("INFO: %s: %s", funcName, info)
	}
}

func LogWarning(warning string, funcName string, c configuration.Configuration) {
	if c.General.Logging && (c.General.LoggingLevel == "INFO" || c.General.LoggingLevel == "WARNING") {
		log.Printf("WARNING: %s: %s", funcName, warning)
	}
}

func LogError(err error, funcName string, c configuration.Configuration) {
	if c.General.Logging && (c.General.LoggingLevel == "INFO" || c.General.LoggingLevel == "WARNING" || c.General.LoggingLevel == "ERROR") {
		log.Printf("ERROR: %s: %s", funcName, err)
	}
}

func LogFatal(err error, funcName string, c configuration.Configuration) {
	if c.General.Logging && (c.General.LoggingLevel == "INFO" || c.General.LoggingLevel == "WARNING" || c.General.LoggingLevel == "ERROR" || c.General.LoggingLevel == "FATAL") {
		message := fmt.Sprintf("FATAL: %s: %s", funcName, err)
		log.Fatal("FATAL:", message)
	}
}
