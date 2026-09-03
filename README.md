# ConsentVault

A consent and preference ledger for applications that need clear, auditable privacy choices.

ConsentVault records versioned purposes and policies, captures explicit user decisions, supports withdrawal, and provides a tamper-evident event history. It is designed as a focused service, not a legal-compliance claim.

## Core capabilities

- Versioned consent purposes and notices
- Explicit grant and withdrawal records
- Current-consent evaluation
- Append-only audit events
- Idempotent API operations
- Retention and export workflows

## Technology

Go 1.24, the standard `net/http` server, PostgreSQL, structured logging, table-driven tests, and GitHub Actions.

## Local setup

1. Install Go 1.24.
2. Copy `.env.example` to `.env`.
3. Run `go test ./...`.
4. Start with `go run ./cmd/api`.
5. Check `http://localhost:8080/healthz`.

## Quality commands

```bash
go test -race ./...
go vet ./...
```

## License

MIT
