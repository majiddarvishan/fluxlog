# Fluxlog — independent release readiness (2026-10-19)

## Verified upstream release

The upstream `github.com/majiddarvishan/goconfig` **v1.7.0** annotated tag
exists, object `f3263dd3160cf226ec0f3371d7753615862d0c32`, pointing to
commit `213e87a1677af862fa9d6226f4e6ea8d0b0edfbc`.
Upstream main CI [39167180503](https://github.com/majiddarvishan/goconfig/actions/runs/39167180503)
passed. A GitHub Release object was not found; the Go module tag itself
exists and is resolvable.

In [module resolution CI 39167874434](https://github.com/majiddarvishan/fluxlog/actions/runs/39167874434),
Go 1.23 resolved `goconfig v1.7.0` and verified its VCS origin commit.
Go-generated module hashes are recorded in the nested `go.sum`:

```text
github.com/majiddarvishan/goconfig v1.7.0 h1:Utm3i4JkNtxBr0VYDmNf0xY1+0A+cMLlxsBllbKMkww=
github.com/majiddarvishan/goconfig v1.7.0/go.mod h1:oNFfQuyPt3xp2tdZSl3vutSXxnsHNvgOBBSnpqrMYuc=
```

## Migration and release state

- Official upstream module upgrade to `v1.7.0` is prepared in the
  `release/fluxlog-goconfig-v1.7.0-migration` branch, replacing the
  earlier development pseudo-version. Runtime Go source and root module
  dependency files are unchanged. The checksum-only temporary workflow
  is removed from the final diff.
- Existing Fluxlog `main@70b17302e5cc5b2974fff502227552ed8cf4d694`
  passed [CI 39167270820](https://github.com/majiddarvishan/fluxlog/actions/runs/39167270820)
  for both core and adapter before this migration.
- Full CI on the final migration PR commit, and then post-merge main,
  must pass before marking the migration complete.
- The current *published* Fluxlog root version is `v0.2.1`. The current
  *published* nested adapter is `adapters/goconfig/v0.1.0`. Neither is
  automatically updated by an upstream goconfig tag.
- The recommended next root `v0.3.0` and nested
  `adapters/goconfig/v0.2.0` tags are not published. Tagging them is a
  separate step after the gates in [RELEASING.md](../RELEASING.md).

## Remaining work

1. Verify the tagged upstream dependency migration under Go CI.
2. Publish Fluxlog root `v0.3.0` at the exact tested main commit using
   an authorized Git client.
3. Migrate the adapter's dependency on root Fluxlog `v0.2.1` to `v0.3.0`;
   run CI again and publish the nested adapter tag at that verified commit.
4. Remove merged working branches and create release objects where supported.

Do not invent checksums, tags or successful CI. Preserve the goconfig
first-precommit-handler limitation and Manager-lifetime listener retention.
