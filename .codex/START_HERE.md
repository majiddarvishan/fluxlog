# Fluxlog — Start Here

> Persistent engineering knowledge for **github.com/majiddarvishan/fluxlog**.
> This documentation was established on **2026-10-18** from the complete tracked source tree of `main` at commit `f0b79fcd7afeb0406e75a1c2ca5fe1e15243e881`. It is a **source-based snapshot**, not a substitute for verifying later changes.

## Thirty-second orientation

Fluxlog is a reusable **Go structured-logging library**. The core Go module wraps `rs/zerolog` with instance-scoped, atomically changeable severity filtering, synchronized output lifecycle, custom console formatting, and optional size-rotated logs through `lumberjack`. It has both a `*fluxlog.Logger` API returning `*zerolog.Event` and convenience package functions. A separate Go module under `adapters/goconfig` that follows lockstep release versions maps a legacy configuration shape and subscribes to runtime level replacements.

This is **not** an application, network service, remote log collector, asynchronous queue, or general-purpose configuration framework.

## Read these files in order

1. [PRODUCT_OVERVIEW.md](PRODUCT_OVERVIEW.md) — capabilities, boundaries, intended use.
2. [ARCHITECTURE.md](ARCHITECTURE.md) — event flow, concurrency and ownership.
3. [API_REFERENCE.md](API_REFERENCE.md) — exact public API semantics and pitfalls.
4. [CONFIGURATION.md](CONFIGURATION.md) — defaults and optional adapter.
5. [CODEBASE_MAP.md](CODEBASE_MAP.md) — source files and responsibilities.
6. [TESTING.md](TESTING.md) — how to verify both Go modules and evidence status.
7. [DECISIONS.md](DECISIONS.md) — existing design invariants inferred from implementation.
8. [CURRENT_STATUS.md](CURRENT_STATUS.md), [TASKS.md](TASKS.md), and [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md) — actual status, ranked tasks, bounded execution gates.
9. [SESSION.md](SESSION.md) and [NEXT_CHAT_PROMPT.md](NEXT_CHAT_PROMPT.md) — reproducible handoff.

## Source of truth and refresh protocol

- **Code, tests, `go.mod`, and CI definitions are authoritative**. These notes record verified facts at the baseline commit.
- Before changing anything: read `CURRENT_STATUS.md`; inspect the **actual current HEAD and diff** from the documented baseline. Do not automatically repeat a full codebase audit.
- If HEAD has changed, inspect only changed areas first, updating relevant docs and the snapshot SHA. Perform a broader audit only when the changed area or uncertainty warrants it.
- Keep `.codex` in **English**; distinguish implemented behavior, inferred behavior, and proposals.
- Do not silently alter public APIs, output formatting, concurrency/lifecycle guarantees, or adapter configuration compatibility.
- Root and `adapters/goconfig` are **separate Go modules**: validate both.

## Repository entry points

- Root: [README.md](../README.md), [logger.go](../logger.go), [config.go](../config.go), [writer.go](../writer.go).
- Adapter: [adapters/goconfig/adapter.go](../adapters/goconfig/adapter.go).
- Tests: `*_test.go` at root and `adapters/goconfig/adapter_test.go`.
- CI: [`.github/workflows/ci.yml`](../.github/workflows/ci.yml).
- Development: [Makefile](../Makefile).

**Published release checkpoint (2026-10-19):** both core `v0.3.0` and the independent nested `adapters/goconfig/v0.3.0` tag are published and CI-verified. See `CURRENT_STATUS.md`, `RELEASE_READINESS_2026-10-19.md` and `RELEASING.md` for exact commits, tests and versioning policy.
