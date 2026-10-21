# condition-github-actions

[![Latest Release](https://img.shields.io/github/v/release/SemRels/condition-github-actions?label=version&color=blue)](https://github.com/SemRels/condition-github-actions/releases/latest)

Allows releases only when semrel is running inside GitHub Actions.

This plugin is distributed as the standalone Go binary `semrel-plugin-condition-github-actions`. Semrel executes the binary as a subprocess, provides plugin configuration through `SEMREL_PLUGIN_*` environment variables, provides release context through `SEMREL_*` environment variables, reads standard output, and treats exit code `0` as success and any non-zero exit code as failure. Install the binary in `~/.semrel/plugins/` or anywhere on your `$PATH`.

## Installation

### Binary

```bash
go install github.com/SemRels/condition-github-actions/cmd/plugin@latest
```

### Docker

Pre-built, multi-platform images (linux/amd64, linux/arm64) are published to the GitHub Container Registry on every release:

```bash
docker pull ghcr.io/semrels/condition-github-actions:latest
```

Images are signed with [cosign](https://github.com/sigstore/cosign) and include a full SBOM attestation. Verify the signature:

```bash
cosign verify ghcr.io/semrels/condition-github-actions:latest \
  --certificate-identity-regexp 'https://github.com/SemRels/condition-github-actions/.github/workflows/container.yml.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com
```


## Configuration

```yaml
plugins:
  - name: condition-github-actions
    path: ~/.semrel/plugins/semrel-plugin-condition-github-actions
    env:
      {}
```

## `SEMREL_PLUGIN_*` variables

| Name | Required | Description | Default |
| --- | --- | --- | --- |
| `SEMREL_PLUGIN_BRANCH` | No | Optional branch name override. When set, requires `GITHUB_REF_NAME` (or `GITHUB_REF` stripped of `refs/heads/`) to match this value. | `""` |

## GitHub Actions environment variables

| Name | Required | Description |
| --- | --- | --- |
| `GITHUB_ACTIONS` | Yes | Must be `"true"`. |
| `GITHUB_TOKEN` / `GH_TOKEN` | Yes | At least one token variable must be non-empty. |
| `GITHUB_REF_NAME` | No | Checked when `SEMREL_PLUGIN_BRANCH` is set; falls back to `GITHUB_REF` (`refs/heads/<branch>`) on older runners. |
| `GITHUB_REF` | No | Fallback branch reference when `GITHUB_REF_NAME` is unset. |
| `GITHUB_SERVER_URL` | No | Optional GitHub server URL (for example `https://github.com`, `https://ghe.example.com`, or `https://octocorp.ghe.com`). Validated as an `http`/`https` URL when set. |
| `GITHUB_API_URL` | No | Optional GitHub API URL (for example `https://api.github.com`, `https://ghe.example.com/api/v3`, or `https://api.octocorp.ghe.com`). Validated as an `http`/`https` URL when set. |

## GitHub Enterprise Server (GHES) & Data Residency Support

`condition-github-actions` works out of the box on:
- **github.com**
- **GitHub Enterprise Server (GHES)** (for example `GITHUB_SERVER_URL=https://ghe.example.com` and `GITHUB_API_URL=https://ghe.example.com/api/v3`)
- **GitHub Enterprise Cloud with data residency** (`*.ghe.com`, for example `GITHUB_SERVER_URL=https://octocorp.ghe.com` and `GITHUB_API_URL=https://api.octocorp.ghe.com`)

No `github.com` hostnames are hardcoded, and older GHES runners that only populate `GITHUB_REF` without `GITHUB_REF_NAME` are supported via the `refs/heads/` fallback.

## `SEMREL_*` release context used

This plugin does not consume any `SEMREL_*` release context variables directly.

## Example behavior

The plugin checks the CI environment and succeeds when `GITHUB_ACTIONS=true`. Outside GitHub Actions it exits non-zero to stop the release.

## License

Apache-2.0
