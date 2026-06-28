package anigate

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Event struct {
	Time      time.Time      `json:"time"`
	Kind      string         `json:"kind"`
	Tool      string         `json:"tool,omitempty"`
	JobID     string         `json:"job_id,omitempty"`
	Workspace string         `json:"workspace,omitempty"`
	Path      string         `json:"path,omitempty"`
	Preset    string         `json:"preset,omitempty"`
	OK        bool           `json:"ok"`
	Message   string         `json:"message,omitempty"`
	Fields    map[string]any `json:"fields,omitempty"`
}

type EventLog struct {
	path string
	mu   sync.Mutex
}

type EventFilter struct {
	Kind string
	Tool string
}

func NewEventLog(stateDir string) (*EventLog, error) {
	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		return nil, err
	}
	return &EventLog{path: filepath.Join(stateDir, "events.ndjson")}, nil
}

func (l *EventLog) Append(ev Event) {
	if l == nil {
		return
	}
	if ev.Time.IsZero() {
		ev.Time = time.Now().UTC()
	}
	b, err := json.Marshal(ev)
	if err != nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = f.Write(append(b, '\n'))
}

func (l *EventLog) Tail(limit int, filter EventFilter) ([]Event, error) {
	if l == nil {
		return nil, errors.New("event log is not configured")
	}
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	f, err := os.Open(l.path)
	if errors.Is(err, os.ErrNotExist) {
		return []Event{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var events []Event
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 4096), 1024*1024)
	for scanner.Scan() {
		var ev Event
		if err := json.Unmarshal(scanner.Bytes(), &ev); err != nil {
			continue
		}
		if filter.Kind != "" && ev.Kind != filter.Kind {
			continue
		}
		if filter.Tool != "" && ev.Tool != filter.Tool {
			continue
		}
		events = append(events, ev)
		if len(events) > limit {
			copy(events, events[len(events)-limit:])
			events = events[:limit]
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return events, nil
}
