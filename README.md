# condition-github-actions

GitHub Actions CI condition plugin for SemRel.

Checks GitHub Actions workflow state before SemRel publishes a release.

## Documentation

- SemRel docs (planned): <https://github.com/SemRels/semrel/tree/main/docs/plugins/condition-github-actions>
- Plugin template: <https://github.com/SemRels/plugin-template>
- Registry: <https://registry.semrel.io>

## Repository Layout

~~~text
cmd/plugin/              Plugin entry point
internal/plugin/         Business logic scaffold
internal/grpc/           gRPC transport scaffold
proto/v1                 Symlink to the SemRel protobuf contract
.github/workflows/       CI, release, and security automation
~~~

## Development

~~~bash
go build ./cmd/plugin
go test ./...
~~~

## Configuration Example

~~~yaml
plugins:
  - name: condition-github-actions
    type: condition
    config:
      required_workflows:
        - CI
        - Security
      required_conclusion: success
      token: ${GITHUB_TOKEN}
~~~

## Status

This repository is bootstrapped from SemRels/plugin-template and is ready for implementation.
