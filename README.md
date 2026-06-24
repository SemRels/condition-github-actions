# condition-github-actions

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
  --certificate-identity-regexp 'https://github.com/SemRels/condition-github-actions/.github/workflows/release.yml.*' \
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
| _None_ | - | This plugin does not require any `SEMREL_PLUGIN_*` variables. It relies on CI-provided environment state. | - |

## `SEMREL_*` release context used

This plugin does not consume any `SEMREL_*` release context variables directly.

## Example behavior

The plugin checks the CI environment and succeeds when `GITHUB_ACTIONS=true`. Outside GitHub Actions it exits non-zero to stop the release.

## License

Apache-2.0
