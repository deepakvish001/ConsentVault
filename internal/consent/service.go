package consent

import (
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "fmt"
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
    ID        string    `json:"id"`
    SubjectID string    `json:"subjectId"`
    Purpose   string    `json:"purpose"`
    Version   string    `json:"version"`
    Decision  Decision  `json:"decision"`
    Occurred  time.Time `json:"occurredAt"`
    PreviousHash string  `json:"previousHash,omitempty"`
    Hash      string     `json:"hash"`
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
    if event.ID == "" {
        event.ID = fmt.Sprintf("evt-%d", event.Occurred.UnixNano())
    }
    if len(l.events) > 0 {
        event.PreviousHash = l.events[len(l.events)-1].Hash
    }
    event.Hash = eventHash(event)
    l.events = append(l.events, event)
    return nil
}

func eventHash(event Event) string {
    canonical := strings.Join([]string{
        event.ID,
        event.SubjectID,
        event.Purpose,
        event.Version,
        string(event.Decision),
        event.Occurred.UTC().Format(time.RFC3339Nano),
        event.PreviousHash,
    }, "\x1f")
    digest := sha256.Sum256([]byte(canonical))
    return hex.EncodeToString(digest[:])
}

func (l *Ledger) Events() []Event {
    l.mu.RLock()
    defer l.mu.RUnlock()
    result := make([]Event, len(l.events))
    copy(result, l.events)
    return result
}

func (l *Ledger) Verify() error {
    l.mu.RLock()
    defer l.mu.RUnlock()
    previous := ""
    for index, event := range l.events {
        if event.PreviousHash != previous {
            return fmt.Errorf("consent chain broken at event %d", index)
        }
        if event.Hash != eventHash(event) {
            return fmt.Errorf("consent event %d failed integrity verification", index)
        }
        previous = event.Hash
    }
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
