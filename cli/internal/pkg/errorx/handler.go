package errorx

import (
	"log/slog"
	"os"
)

func CatchAndWrap(err error) {
	if err == nil {
		return
	}

	// TODO: mapping error with more details and exit code
	slog.Error("", "error", err.Error())
	os.Exit(2)
}
