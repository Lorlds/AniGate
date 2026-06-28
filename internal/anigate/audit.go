package anigate

import (
	"time"
)

func (s *Service) auditSummary(args map[string]any) (map[string]any, error) {
	sinceSec := int64ArgDefault(args, "since_sec", 24*3600)
	if sinceSec <= 0 || sinceSec > 31*24*3600 {
		sinceSec = 24 * 3600
	}
	events, err := s.events.Tail(1000, EventFilter{})
	if err != nil {
		return nil, err
	}
	cutoff := time.Now().UTC().Add(-time.Duration(sinceSec) * time.Second)
	byKind := map[string]int{}
	byTool := map[string]int{}
	failures := 0
	var recentFailures []Event
	for _, ev := range events {
		if ev.Time.Before(cutoff) {
			continue
		}
		byKind[ev.Kind]++
		if ev.Tool != "" {
			byTool[ev.Tool]++
		}
		if !ev.OK {
			failures++
			if len(recentFailures) < 10 {
				recentFailures = append(recentFailures, ev)
			}
		}
	}
	return map[string]any{
		"since":           cutoff.Format(time.RFC3339),
		"events_scanned":  len(events),
		"by_kind":         byKind,
		"by_tool":         byTool,
		"failures":        failures,
		"recent_failures": recentFailures,
	}, nil
}
