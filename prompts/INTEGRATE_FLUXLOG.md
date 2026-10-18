# Integrate Fluxlog — Reusable Coding-Agent Prompt

Use this prompt in the **consuming Go application's repository** when asking
Codex or another coding agent to integrate Fluxlog. It is designed to be
copied as a single request. Replace optional requirements with your real
application requirements; do not assume the remote library's `AGENTS.md`
is automatically loaded in the consuming project.

## Copy/paste prompt

```text
Integrate Fluxlog into this Go project.

Library:
https://github.com/majiddarvishan/fluxlog

Read the official documentation before editing code:
- https://github.com/majiddarvishan/fluxlog/blob/main/docs/INTEGRATION_GUIDE.md
- https://github.com/majiddarvishan/fluxlog/blob/main/docs/CONFIGURATION_REFERENCE.md
- https://github.com/majiddarvishan/fluxlog/blob/main/docs/RECIPES.md
- https://github.com/majiddarvishan/fluxlog/tree/main/examples
- https://github.com/majiddarvishan/fluxlog/tree/main/adapters/goconfig/examples

Inspect this project's architecture, existing logging, Go modules,
configuration provider, deployment environment, and shutdown lifecycle.
Choose the appropriate integration pattern instead of assuming one.

Use approved published versions as the initial baseline:
- Core: github.com/majiddarvishan/fluxlog@v0.3.0
- Optional adapter:
  github.com/majiddarvishan/fluxlog/adapters/goconfig@v0.3.0
Verify actual dependency and API compatibility before implementation.

Decide and explain:
1. Core-only vs. goconfig adapter: add the adapter only if the project
   already uses goconfig and the adapter's supported mapping meets its needs.
2. Package-level fluxlog.Info(...) vs. explicitly injected *fluxlog.Logger:
   choose based on shared ownership and independent component requirements.
3. Plain messages vs. structured logging for request IDs, components,
   status codes, durations, and other operational context.
4. Console vs. rotating file vs. both; human-readable vs. JSON output;
   use the deployment's existing log collection design.
5. Level updates through Logger.SetLevelString or goconfig Manager.Replace;
   do not assume the adapter hot-reloads any setting other than level.
6. Timestamp layout, caller and service fields, color handling,
   retention, and shutdown behavior based on actual requirements.

Implement the smallest correct and maintainable integration.
Use real Fluxlog APIs and supported configuration fields only; do not
invent adapter JSON keys or change the library's public API.
Respect Logger ownership: SetDefault does not transfer ownership;
restore the previous default before closing the active logger, and do
not close it while it remains in use. Do not introduce a goconfig
dependency into the core-only application without a real need.

Do not log passwords, tokens, authentication secrets, or sensitive
payloads. Preserve existing functional behavior and observability
contracts; migrate existing loggers only where justified.

After implementation, run available formatting, dependency checks,
go vet, relevant tests (including race checks when practical), and build.
Report which checks passed, which could not run, and why.

Finish with a short summary of:
- selected integration pattern and its rationale;
- configuration and output choices;
- files changed;
- verification results;
- supported limitations or follow-up work.

Do not bump dependency versions beyond the approved baseline, create
Git tags, or perform unrelated refactoring without explicit approval.
```

## Choosing a narrower request

The default prompt lets the agent inspect the project and choose. If the
project's deployment contract is known, prepend one or more explicit
constraints, for example:

- `The application already uses goconfig v1.7.0. Keep its Manager and
  support committed runtime changes to the logger level.`
- `The service is containerized. Emit JSON lines to stdout only; do not
  create local log files.`
- `Use human-readable console timestamps in local time and the Go layout
  "2006-01-02 15:04:05".`
- `Keep independent loggers per component; do not replace the package
  default.`
- `Use package-level fluxlog.Info for application-wide messages and
  structured instance events where request_id is needed.`

Consult the [Integration Guide](../docs/INTEGRATION_GUIDE.md),
[Configuration Reference](../docs/CONFIGURATION_REFERENCE.md), and
[Recipes](../docs/RECIPES.md) for feature support and runnable examples.

**Version note:** this prompt uses the published `v0.3.0` module API as its
baseline. Documentation on the repository's `main` branch may be newer than
the code at an immutable release tag; verify examples and APIs against the
selected dependency version.
