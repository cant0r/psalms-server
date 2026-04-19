//go:build release

package configurations

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func NewLogger() *log.Logger {
	return log.NewWithOptions(os.Stdout, log.Options{
		TimeFormat:      time.RFC3339,
		ReportTimestamp: true,
		ReportCaller:    true,
		Level:           log.InfoLevel,
	})
}
