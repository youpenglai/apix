package logger

import (
	"github.com/youpenglai/goutils/logger"
	"os"
	"path/filepath"
)

const(
	AccessLog = "access.log"
	ErrorLog = "error.log"
	RunLog = "run.log"

	PrefixAccess = "ApiXAccess"
	PrefixRun = "ApiXRun"
	PrefixError = "ApiXError"
)

var (
	LogPath = "logs"
)

func initLoggers() {
	runLevel := os.Getenv("APIX_RUN_LEVEL")
	if runLevel == "" {
		runLevel = "info"
	}

	accessLogPath := filepath.Join(LogPath, AccessLog)
	errorLogPath := filepath.Join(LogPath, ErrorLog)
	runLogPath := filepath.Join(LogPath, RunLog)

	logger.InitLogger(runLevel,
	//	logger.LoggerOpts{
	//	Type: logger.LoggerConsole,
	//},
	logger.LoggerOpts{
		Type:logger.LoggerFile,
		Rotate: logger.LogRotateDaily,
		FileName: accessLogPath,
		Prefixs: []string{PrefixAccess},
	}, logger.LoggerOpts{
		Type: logger.LoggerFile,
		Rotate: logger.LogRotateDaily,
		FileName: errorLogPath,
		Prefixs: []string{PrefixError},
	}, logger.LoggerOpts{
		Type: logger.LoggerFile,
		Rotate: logger.LogRotateDaily,
		FileName: runLogPath,
		Prefixs: []string{PrefixRun},
	})

	logger.GetLogger(PrefixAccess).SetFormatter(ApiXLoggerFormat)
	logger.GetLogger(PrefixError).SetFormatter(ApiXLoggerFormat)
	logger.GetLogger(PrefixRun).SetFormatter(ApiXLoggerFormat)
}

func init() {
	initLoggers()
}
