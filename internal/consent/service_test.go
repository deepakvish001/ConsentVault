package consent

import (
    "testing"
    "time"
)

func TestWithdrawalReplacesLatestDecision(t *testing.T) {
    ledger := &Ledger{}
    _ = ledger.Record(Event{SubjectID: "user-1", Purpose: "analytics", Version: "v1", Decision: Granted})
    _ = ledger.Record(Event{SubjectID: "user-1", Purpose: "analytics", Version: "v1", Decision: Withdrawn})
    if ledger.IsGranted("user-1", "analytics", "v1") { t.Fatal("expected consent to be withdrawn") }
}

func TestLedgerCreatesVerifiableHashChain(t *testing.T) {
    ledger := &Ledger{}
    occurred := time.Date(2026, time.May, 1, 10, 30, 0, 0, time.UTC)
    if err := ledger.Record(Event{ID: "evt-1", SubjectID: "user-1", Purpose: "analytics", Version: "v1", Decision: Granted, Occurred: occurred}); err != nil {
        t.Fatal(err)
    }
    if err := ledger.Record(Event{ID: "evt-2", SubjectID: "user-1", Purpose: "analytics", Version: "v1", Decision: Withdrawn, Occurred: occurred.Add(time.Minute)}); err != nil {
        t.Fatal(err)
    }
    events := ledger.Events()
    if events[0].Hash == "" || events[1].PreviousHash != events[0].Hash {
        t.Fatal("expected events to form a hash chain")
    }
    if err := ledger.Verify(); err != nil {
        t.Fatalf("expected valid ledger: %v", err)
    }
}

func TestLedgerDetectsTampering(t *testing.T) {
    ledger := &Ledger{}
    _ = ledger.Record(Event{ID: "evt-1", SubjectID: "user-1", Purpose: "email", Version: "v1", Decision: Granted})
    ledger.events[0].Purpose = "advertising"
    if err := ledger.Verify(); err == nil {
        t.Fatal("expected tampering to be detected")
    }
}
