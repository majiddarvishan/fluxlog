# Continue Fluxlog official release pipeline

Act as a **Senior Go Software Architect / Release Engineer**.
Explain in **Persian**. All repo docs, comments and code must be **English**.

## Canonical baseline (2026-10-19)

- Fluxlog main before migration: `70b17302e5cc5b2974fff502227552ed8cf4d694`, both CI jobs passed in [39167270820](https://github.com/majiddarvishan/fluxlog/actions/runs/39167270820).
- goconfig `v1.7.0` **exists as an annotated Git tag**, pointing to
  `213e87a1677af862fa9d6226f4e6ea8d0b0edfbc`. The upstream post-merge CI
  passed in [39167180503](https://github.com/majiddarvishan/goconfig/actions/runs/39167180503).
- Fluxlog migration branch `release/fluxlog-goconfig-v1.7.0-migration`
  changes only nested adapter `go.mod/go.sum` (and English docs) from
  the previously verified goconfig pseudo-version to **official `v1.7.0`**.
- Official Go module origin/checksums were obtained through Go 1.23 in
  [39167874434](https://github.com/majiddarvishan/fluxlog/actions/runs/39167874434).
  All temporary checksum workflows should be absent from final diff.
- Strict CFG-01/CFG-02/CFG-03 tests and both modules' checks **PASSED with tagged v1.7.0** in [CI 39168203987](https://github.com/majiddarvishan/fluxlog/actions/runs/39168203987); do not confuse that PR-run evidence with post-merge CI.

## Remaining steps

1. Verify current migration PR's **final** root+adapter CI (fmt, tidy,
   vet, race, examples) and its merge status. If successful, merge and
   verify both post-merge jobs; update REL-03 evidence. Inspect actual
   GitHub HEAD before using any cached SHA.
2. A future official root tag `v0.3.0` is proposed: create only on
   verified main HEAD after both jobs pass. The connector does not
   support directly creating tags, so use an authorized Git client.
3. Once root `v0.3.0` is published and Go-resolvable, update the nested
   adapter module's Fluxlog core dependency from `v0.2.1` to `v0.3.0`,
   regenerate `go.sum` with Go, test both modules, merge, and tag the
   adapter independently as `adapters/goconfig/v0.2.0` at its tested SHA.
4. Do not invent or move tags, do not make the root module depend on
   goconfig, and don't claim adapter v0.1.0 contains new behavior.
5. Release instructions: `RELEASING.md`; tracking: `.codex/TASKS.md`,
   `.codex/RELEASE_READINESS_2026-10-19.md` and `.codex/CURRENT_STATUS.md`.

Keep work scoped and verifiable, without speculative features.
