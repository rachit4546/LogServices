package processor

import (
	"errors"
	"strings"
	"time"
)

type ParsedLine struct {
	Level string
	Hour  string
	Msg   string
}

var ErrIncorrectFormat = errors.New("incorrect log format")

func ParseLine(line string) (*ParsedLine, error) {
	parts := strings.SplitN(line, " ", 3)
	if len(parts) < 3 {
		return nil, ErrIncorrectFormat
	}

	ts := parts[0]
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return nil, ErrIncorrectFormat
	}

	level := parts[1]
	validLevels := map[string]bool{
		"INFO":  true,
		"WARN":  true,
		"ERROR": true,
		"DEBUG": true,
	}
	if !validLevels[level] {
		return nil, ErrIncorrectFormat
	}

	msg := parts[2]

	hour := t.Format("2006-01-02 15:00")
	return &ParsedLine{Level: level, Hour: hour, Msg: msg}, nil
}
