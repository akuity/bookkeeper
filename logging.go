package render

import log "github.com/sirupsen/logrus"

// LogLevel represents the level of detail logged by the Kargo Render service's
// internal logger.
type LogLevel log.Level

const (
	// LogLevelDebug represents DEBUG level logging.
	LogLevelDebug = LogLevel(log.DebugLevel)
	// LogLevelInfo represents INFO level logging. This is the default for the
	// Kargo Render service when no LogLevel is explicitly specified.
	LogLevelInfo = LogLevel(log.InfoLevel)
	// LogLevelError represents ERROR level logging.
	LogLevelError = LogLevel(log.ErrorLevel)
)
