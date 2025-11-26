package processor

import (
	"strings"
	"time"
)

type ParsedLine struct {
	Level string
	Hour  string
	Msg   string
}

func ParseLine(line string) (*ParsedLine, error) {
	parts := strings.SplitN(line, " ", 3)
	if len(parts) < 3 {
		return nil, nil
	}

	ts := parts[0]
	level := parts[1]
	msg := parts[2]
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return nil, err
	}

	hour := t.Format("2006-01-02 15:00")
	return &ParsedLine{Level: level, Hour: hour, Msg: msg}, nil
}
