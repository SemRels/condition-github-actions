// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The SemRels Authors

package plugin

import (
	"context"
	"testing"

	semrelv1 "github.com/SemRels/semrel-api/api/gen/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func env(kv map[string]string) func(string) string {
	return func(key string) string { return kv[key] }
}

func validEnv() map[string]string {
	return map[string]string{
		"GITHUB_ACTIONS": "true",
		"GITHUB_TOKEN":   "ghp_test",
		"GITHUB_REF_NAME": "main",
	}
}

func TestVerifyConditions_HappyPath(t *testing.T) {
	t.Parallel()

	c := NewWithEnv(env(validEnv()))
	resp, err := c.VerifyConditions(context.Background(), &semrelv1.VerifyConditionsRequest{
		Ctx: &semrelv1.ReleaseContext{Branch: "main"},
	})

	require.NoError(t, err)
	assert.Empty(t, resp.GetErrorMessage())
}

func TestVerifyConditions_NotInGitHubActions(t *testing.T) {
	t.Parallel()

	e := validEnv()
	delete(e, "GITHUB_ACTIONS")

	c := NewWithEnv(env(e))
	resp, err := c.VerifyConditions(context.Background(), &semrelv1.VerifyConditionsRequest{})

	require.NoError(t, err)
	assert.Contains(t, resp.GetErrorMessage(), "GITHUB_ACTIONS")
}

func TestVerifyConditions_MissingToken(t *testing.T) {
	t.Parallel()

	e := validEnv()
	delete(e, "GITHUB_TOKEN")

	c := NewWithEnv(env(e))
	resp, err := c.VerifyConditions(context.Background(), &semrelv1.VerifyConditionsRequest{})

	require.NoError(t, err)
	assert.Contains(t, resp.GetErrorMessage(), "GITHUB_TOKEN")
}

func TestVerifyConditions_GHTokenAlternative(t *testing.T) {
	t.Parallel()

	e := validEnv()
	delete(e, "GITHUB_TOKEN")
	e["GH_TOKEN"] = "ghp_test"

	c := NewWithEnv(env(e))
	resp, err := c.VerifyConditions(context.Background(), &semrelv1.VerifyConditionsRequest{})

	require.NoError(t, err)
	assert.Empty(t, resp.GetErrorMessage())
}

func TestVerifyConditions_BranchMismatch(t *testing.T) {
	t.Parallel()

	e := validEnv()
	e["GITHUB_REF_NAME"] = "feature/something"

	c := NewWithEnv(env(e))
	resp, err := c.VerifyConditions(context.Background(), &semrelv1.VerifyConditionsRequest{
		Ctx: &semrelv1.ReleaseContext{Branch: "main"},
	})

	require.NoError(t, err)
	assert.Contains(t, resp.GetErrorMessage(), "branch mismatch")
}

func TestVerifyConditions_BranchFromGITHUB_REF(t *testing.T) {
	t.Parallel()

	e := validEnv()
	delete(e, "GITHUB_REF_NAME")
	e["GITHUB_REF"] = "refs/heads/main"

	c := NewWithEnv(env(e))
	resp, err := c.VerifyConditions(context.Background(), &semrelv1.VerifyConditionsRequest{
		Ctx: &semrelv1.ReleaseContext{Branch: "main"},
	})

	require.NoError(t, err)
	assert.Empty(t, resp.GetErrorMessage())
}

func TestVerifyConditions_MultipleErrors(t *testing.T) {
	t.Parallel()

	c := NewWithEnv(env(map[string]string{}))
	resp, err := c.VerifyConditions(context.Background(), &semrelv1.VerifyConditionsRequest{})

	require.NoError(t, err)
	assert.Contains(t, resp.GetErrorMessage(), "GITHUB_ACTIONS")
	assert.Contains(t, resp.GetErrorMessage(), "GITHUB_TOKEN")
}
