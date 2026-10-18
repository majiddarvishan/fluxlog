# Current status — Fluxlog

**Updated:** 2026-10-19  
**Base main before PR #18:** `3c3b38352e44eb311f7841f670c8071ae7cf41a7`  
**Development PR:** [#18 — Post-commit shared level binding](https://github.com/majiddarvishan/fluxlog/pull/18)  
**Final pre-merge test evidence:** [GitHub Actions run 39165664802](https://github.com/majiddarvishan/fluxlog/actions/runs/39165664802) **PASS** for both Go modules.

## Completed implementation

The optional `adapters/goconfig` module now pins a reproducible upstream
pseudo-version `v1.6.1-0.20261018202246-20a58eb9a426` (exact commit
`20a58eb9a426d0ed7d522b1e5848a0971ce58198`). This is a deliberate
development integration, **not** a newly tagged upstream version.
Checksum provenance: [Go module resolution run 39165338586](https://github.com/majiddarvishan/fluxlog/actions/runs/39165338586).

`BindLevel` validates incoming levels before persistence and updates all
bound loggers only after a successful committed configuration change, using
the new upstream `Manager.OnReplaceCommitted`. The historical known-issue
CFG-02 and CFG-03 tests are now strict **PASS** assertions, also covering
successful/failed Batch, invalid values and real FileSource failure injection.

`gofmt`, `go mod tidy`, `go vet`, `go test -race`, and example builds
passed for both the root and adapter jobs in run 39165664802. No root
runtime source or root `go.mod` was changed.

## Tracking

**Complete in development source:** UP-01, UP-02, CFG-01, CFG-02, CFG-03,
LVL-01, FMT-01, FMT-02, LIFE-01, FILE-01, ERR-01, DEP-01/02, REL-01 and DOC-01.

**Not yet done:** Publish a stable new goconfig tag and, separately, a new
nested adapter version — neither was authorized or created. Do not say the
published adapter `v0.1.0` contains these changes.

**Supported limitation:** Fluxlog must be the first pre-commit handler
registered for a managed level path to validate every candidate level.
The legacy goconfig `OnReplace` still uses first-match selection.
Post-commit observers are retained for the lifetime of the Manager.
The adapter README documents these restrictions.

**Final process checks:** Verify final PR #18 head CI, merge status and
post-merge main SHA/CI. Clean up temporary branches when possible. All
committed files are English and explanations to the maintainer are Persian.
