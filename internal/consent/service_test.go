package consent

import "testing"

func TestWithdrawalReplacesLatestDecision(t *testing.T) {
    ledger := &Ledger{}
    _ = ledger.Record(Event{SubjectID: "user-1", Purpose: "analytics", Version: "v1", Decision: Granted})
    _ = ledger.Record(Event{SubjectID: "user-1", Purpose: "analytics", Version: "v1", Decision: Withdrawn})
    if ledger.IsGranted("user-1", "analytics", "v1") { t.Fatal("expected consent to be withdrawn") }
}
