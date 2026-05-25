// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The go-semrel Authors

// Package plugin implements the GitHub Actions CI condition check.
package plugin

import (
	"fmt"
	"os"
	"strings"

	semrelv1 "github.com/SemRels/condition-github-actions/api/gen/v1"
)

// EnvProvider abstracts environment-variable access so tests can inject fakes.
type EnvProvider interface {
	Getenv(key string) string
}

// OSEnvProvider implements EnvProvider using the real process environment.
type OSEnvProvider struct{}

// Getenv delegates to os.Getenv.
func (o *OSEnvProvider) Getenv(key string) string {
	return os.Getenv(key)
}

// ConditionChecker verifies that the current environment is a GitHub Actions
// run that is authorised to publish a release.
type ConditionChecker struct {
	env EnvProvider
}

// NewConditionChecker returns a ConditionChecker backed by the supplied EnvProvider.
func NewConditionChecker(env EnvProvider) *ConditionChecker {
	return &ConditionChecker{env: env}
}

// VerifyConditions checks all required GitHub Actions environment variables.
// It maps its findings into a VerifyConditionsResponse:
//   - error_message is non-empty when the release must be aborted.
//   - details always contains the human-readable reason / remediation hint.
func (c *ConditionChecker) VerifyConditions(ctx *semrelv1.ReleaseContext) *semrelv1.VerifyConditionsResponse {
	var failures []string

	// 1. Must be running inside GitHub Actions.
	if c.env.Getenv("GITHUB_ACTIONS") != "true" {
		failures = append(failures, "GITHUB_ACTIONS is not 'true' – this plugin only runs inside GitHub Actions")
	}

	// 2. Must be a branch push, not a tag push.
	if refType := c.env.Getenv("GITHUB_REF_TYPE"); refType != "branch" {
		failures = append(failures, fmt.Sprintf(
			"GITHUB_REF_TYPE is %q – only branch pushes trigger a release", refType,
		))
	}

	// 3. Must not be triggered by a pull-request event.
	if c.env.Getenv("GITHUB_EVENT_NAME") == "pull_request" {
		failures = append(failures, "GITHUB_EVENT_NAME is 'pull_request' – releases are not published from pull-request runs")
	}

	// 4. A GitHub token must be available.
	if c.env.Getenv("GITHUB_TOKEN") == "" {
		failures = append(failures, "GITHUB_TOKEN is not set – a token is required to publish the release")
	}

	// 5. Optional: compare the workflow branch against the release context branch.
	if ctx != nil && ctx.Branch != "" {
		wfBranch := c.env.Getenv("GITHUB_REF_NAME")
		if wfBranch == "" {
			wfBranch = c.env.Getenv("GITHUB_HEAD_REF")
		}
		if wfBranch != "" && wfBranch != ctx.Branch {
			failures = append(failures, fmt.Sprintf(
				"branch mismatch: release context expects %q but workflow is running on %q",
				ctx.Branch, wfBranch,
			))
		}
	}

	if len(failures) == 0 {
		return &semrelv1.VerifyConditionsResponse{
			Details: "all GitHub Actions conditions satisfied",
		}
	}

	return &semrelv1.VerifyConditionsResponse{
		ErrorMessage: strings.Join(failures, "; "),
		Details:      buildRemediationHint(failures),
	}
}

// buildRemediationHint turns a list of failure strings into a single
// human-readable hint block.
func buildRemediationHint(failures []string) string {
	var b strings.Builder
	b.WriteString("The following conditions were not met:\n")
	for i, f := range failures {
		b.WriteString(fmt.Sprintf("  %d. %s\n", i+1, f))
	}
	b.WriteString("\nRemediation:\n")
	b.WriteString("  - Ensure the workflow runs on push to the release branch (not pull_request).\n")
	b.WriteString("  - Add 'permissions: contents: write' and set GITHUB_TOKEN in the workflow.\n")
	b.WriteString("  - Verify GITHUB_REF_TYPE is 'branch' (not 'tag').\n")
	return b.String()
}
