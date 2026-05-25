# condition-github-actions

GitHub Actions CI condition plugin for [SemRel](https://github.com/SemRels/semrel).

Verifies that a release is running inside an authorised GitHub Actions environment
before allowing the release pipeline to proceed.

## What It Checks

1. `GITHUB_ACTIONS=true` — confirms the pipeline runs in GitHub Actions
2. `GITHUB_TOKEN` or `GH_TOKEN` — a token must be present to create releases
3. Branch match (optional) — if a branch is configured in `.semrel.yaml`,
   the plugin compares it against `GITHUB_REF_NAME` / `GITHUB_REF`

## Configuration (`.semrel.yaml`)

```yaml
plugins:
  - name: condition-github-actions
    type: condition
    config:
      branch: main   # optional; defaults to no branch check
```

## Repository Layout

```
cmd/plugin/            Plugin entry point (go-plugin / gRPC serve)
internal/plugin/       Business logic (Condition struct)
internal/grpc/         Legacy stub (replaced by semrel-api/plugin)
```

## Development

```bash
go test ./...
go build ./cmd/plugin
```

## License

Apache-2.0 — see [LICENSE](LICENSE).

