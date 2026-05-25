// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The go-semrel Authors

package plugin

import (
	"testing"

	semrelv1 "github.com/SemRels/condition-github-actions/api/gen/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockEnv implements EnvProvider using a simple map for testing.
type mockEnv map[string]string

func (m mockEnv) Getenv(key string) string {
	return m[key]
}

// happyEnv returns a fully configured, valid GitHub Actions environment.
func happyEnv() mockEnv {
	return mockEnv{
		"GITHUB_ACTIONS":    "true",
		"GITHUB_REF_TYPE":   "branch",
		"GITHUB_EVENT_NAME": "push",
		"GITHUB_TOKEN":      "ghs_test",
		"GITHUB_REF_NAME":   "main",
	}
}

func TestVerifyConditions_AllPassing(t *testing.T) {
	t.Parallel()

	checker := NewConditionChecker(happyEnv())
	resp := checker.VerifyConditions(nil)

	require.Empty(t, resp.ErrorMessage, "should pass with valid environment")
	assert.Contains(t, resp.Details, "satisfied")
}

func TestVerifyConditions_NotGitHubActions(t *testing.T) {
	t.Parallel()

	env := happyEnv()
	delete(env, "GITHUB_ACTIONS")

	checker := NewConditionChecker(env)
	resp := checker.VerifyConditions(nil)

	require.NotEmpty(t, resp.ErrorMessage)
	assert.Contains(t, resp.ErrorMessage, "GITHUB_ACTIONS")
}

func TestVerifyConditions_GitHubActionsFalse(t *testing.T) {
	t.Parallel()

	env := happyEnv()
	env["GITHUB_ACTIONS"] = "false"

	checker := NewConditionChecker(env)
	resp := checker.VerifyConditions(nil)

	require.NotEmpty(t, resp.ErrorMessage)
	assert.Contains(t, resp.ErrorMessage, "GITHUB_ACTIONS")
}

func TestVerifyConditions_RefTypeTag(t *testing.T) {
	t.Parallel()

	env := happyEnv()
	env["GITHUB_REF_TYPE"] = "tag"

	checker := NewConditionChecker(env)
	resp := checker.VerifyConditions(nil)

	require.NotEmpty(t, resp.ErrorMessage)
	assert.Contains(t, resp.ErrorMessage, "GITHUB_REF_TYPE")
}

func TestVerifyConditions_PullRequestEvent(t *testing.T) {
	t.Parallel()

	env := happyEnv()
	env["GITHUB_EVENT_NAME"] = "pull_request"

	checker := NewConditionChecker(env)
	resp := checker.VerifyConditions(nil)

	require.NotEmpty(t, resp.ErrorMessage)
	assert.Contains(t, resp.ErrorMessage, "pull_request")
}

func TestVerifyConditions_MissingToken(t *testing.T) {
	t.Parallel()

	env := happyEnv()
	delete(env, "GITHUB_TOKEN")

	checker := NewConditionChecker(env)
	resp := checker.VerifyConditions(nil)

	require.NotEmpty(t, resp.ErrorMessage)
	assert.Contains(t, resp.ErrorMessage, "GITHUB_TOKEN")
}

func TestVerifyConditions_MultipleFailures(t *testing.T) {
	t.Parallel()

	checker := NewConditionChecker(mockEnv{})
	resp := checker.VerifyConditions(nil)

	require.NotEmpty(t, resp.ErrorMessage)
	// All four mandatory checks should be present.
	assert.Contains(t, resp.ErrorMessage, "GITHUB_ACTIONS")
	assert.Contains(t, resp.ErrorMessage, "GITHUB_REF_TYPE")
	assert.Contains(t, resp.ErrorMessage, "GITHUB_TOKEN")
	assert.Contains(t, resp.Details, "Remediation")
}

func TestVerifyConditions_BranchMatch(t *testing.T) {
	t.Parallel()

	env := happyEnv()
	env["GITHUB_REF_NAME"] = "main"

	checker := NewConditionChecker(env)
	ctx := &semrelv1.ReleaseContext{Branch: "main"}
	resp := checker.VerifyConditions(ctx)

	require.Empty(t, resp.ErrorMessage)
}

func TestVerifyConditions_BranchMismatch(t *testing.T) {
	t.Parallel()

	env := happyEnv()
	env["GITHUB_REF_NAME"] = "feature/x"

	checker := NewConditionChecker(env)
	ctx := &semrelv1.ReleaseContext{Branch: "main"}
	resp := checker.VerifyConditions(ctx)

	require.NotEmpty(t, resp.ErrorMessage)
	assert.Contains(t, resp.ErrorMessage, "branch mismatch")
}

func TestVerifyConditions_BranchMismatch_FallbackToHeadRef(t *testing.T) {
	t.Parallel()

	env := happyEnv()
	delete(env, "GITHUB_REF_NAME")
	env["GITHUB_HEAD_REF"] = "feature/y"

	checker := NewConditionChecker(env)
	ctx := &semrelv1.ReleaseContext{Branch: "main"}
	resp := checker.VerifyConditions(ctx)

	require.NotEmpty(t, resp.ErrorMessage)
	assert.Contains(t, resp.ErrorMessage, "branch mismatch")
}

func TestVerifyConditions_ContextNoBranch_NoBranchCheck(t *testing.T) {
	t.Parallel()

	env := happyEnv()
	env["GITHUB_REF_NAME"] = "some-branch"

	checker := NewConditionChecker(env)
	// ctx.Branch is empty → branch comparison is skipped.
	ctx := &semrelv1.ReleaseContext{Branch: ""}
	resp := checker.VerifyConditions(ctx)

	require.Empty(t, resp.ErrorMessage)
}

func TestOSEnvProvider(t *testing.T) {
	t.Parallel()

	p := &OSEnvProvider{}
	// Should not panic; exact value depends on the test environment.
	_ = p.Getenv("PATH")
}
