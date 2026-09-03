package consent

import (
    "errors"
    "strings"
    "sync"
    "time"
)

type Decision string

const (
    Granted Decision = "granted"
    Withdrawn Decision = "withdrawn"
)

type Event struct {
    SubjectID string    `json:"subjectId"`
    Purpose   string    `json:"purpose"`
    Version   string    `json:"version"`
    Decision  Decision  `json:"decision"`
    Occurred  time.Time `json:"occurredAt"`
}

type Ledger struct {
    mu     sync.RWMutex
    events []Event
}

func (l *Ledger) Record(event Event) error {
    if strings.TrimSpace(event.SubjectID) == "" || strings.TrimSpace(event.Purpose) == "" || strings.TrimSpace(event.Version) == "" {
        return errors.New("subject, purpose, and version are required")
    }
    if event.Decision != Granted && event.Decision != Withdrawn {
        return errors.New("unsupported consent decision")
    }
    if event.Occurred.IsZero() { event.Occurred = time.Now().UTC() }
    l.mu.Lock()
    defer l.mu.Unlock()
    l.events = append(l.events, event)
    return nil
}

func (l *Ledger) IsGranted(subject, purpose, version string) bool {
    l.mu.RLock()
    defer l.mu.RUnlock()
    for index := len(l.events) - 1; index >= 0; index-- {
        event := l.events[index]
        if event.SubjectID == subject && event.Purpose == purpose && event.Version == version {
            return event.Decision == Granted
        }
    }
    return false
}
