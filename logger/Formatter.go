package logger

import (
	"github.com/youpenglai/goutils/logger"
	"fmt"
)

func ApiXLoggerFormat(logMsg *logger.LoggerMsg, colorful bool) string {
	return fmt.Sprintf("%s %s %s\n", logMsg.Prefix, logMsg.Time.Format("2006-01-02 15:04:05.000"), logMsg.Msg)
}
