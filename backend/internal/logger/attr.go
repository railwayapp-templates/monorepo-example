package logger

import (
	"log/slog"
	"strings"
	"time"
)

func ErrAttr(err error) slog.Attr {
	return slog.String("err", strings.TrimSpace(err.Error()))
}

func ErrorsAttr(errors ...error) slog.Attr {
	stringErrors := []string{}

	for _, err := range errors {
		stringErrors = append(stringErrors, strings.TrimSpace(err.Error()))
	}

	return slog.Any("errors", stringErrors)
}

func DurationAttr(d time.Duration, groupName ...string) slog.Attr {
	key := "duration"
	if len(groupName) > 0 {
		key = groupName[0]
	}

	return slog.Group(key,
		slog.String("pretty", d.String()),
		slog.Int64("nanoseconds", d.Nanoseconds()),
		slog.Int64("microseconds", d.Microseconds()),
		slog.Int64("milliseconds", d.Milliseconds()),
		slog.Float64("seconds", d.Seconds()),
	)
}
