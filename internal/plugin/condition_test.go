// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The SemRels Authors

package plugin

import (
	"strings"
	"testing"
)

func env(kv map[string]string) func(string) string {
	return func(key string) string { return kv[key] }
}

func TestCheck_HappyPath(t *testing.T) {
	t.Parallel()

	c := NewWithEnv(env(map[string]string{
		"GITHUB_ACTIONS": "true",
		"GITHUB_TOKEN":   "ghp_test",
	}))

	if err := c.Check(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheck_AllowsGHToken(t *testing.T) {
	t.Parallel()

	c := NewWithEnv(env(map[string]string{
		"GITHUB_ACTIONS": "true",
		"GH_TOKEN":       "ghp_test",
	}))

	if err := c.Check(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheck_NotInGitHubActions(t *testing.T) {
	t.Parallel()

	c := NewWithEnv(env(map[string]string{"GITHUB_TOKEN": "ghp_test"}))
	err := c.Check()
	if err == nil || !strings.Contains(err.Error(), "GITHUB_ACTIONS") {
		t.Fatalf("expected GITHUB_ACTIONS error, got: %v", err)
	}
}

func TestCheck_MissingToken(t *testing.T) {
	t.Parallel()

	c := NewWithEnv(env(map[string]string{"GITHUB_ACTIONS": "true"}))
	err := c.Check()
	if err == nil || !strings.Contains(err.Error(), "GITHUB_TOKEN") {
		t.Fatalf("expected token error, got: %v", err)
	}
}

func TestCheck_BranchMatchWithRefName(t *testing.T) {
	t.Parallel()

	c := NewWithEnv(env(map[string]string{
		"GITHUB_ACTIONS":       "true",
		"GITHUB_TOKEN":         "ghp_test",
		"SEMREL_PLUGIN_BRANCH": "main",
		"GITHUB_REF_NAME":      "main",
	}))

	if err := c.Check(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheck_BranchMatchWithRefFallback(t *testing.T) {
	t.Parallel()

	c := NewWithEnv(env(map[string]string{
		"GITHUB_ACTIONS":       "true",
		"GITHUB_TOKEN":         "ghp_test",
		"SEMREL_PLUGIN_BRANCH": "main",
		"GITHUB_REF":           "refs/heads/main",
	}))

	if err := c.Check(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheck_BranchMismatch(t *testing.T) {
	t.Parallel()

	c := NewWithEnv(env(map[string]string{
		"GITHUB_ACTIONS":       "true",
		"GITHUB_TOKEN":         "ghp_test",
		"SEMREL_PLUGIN_BRANCH": "main",
		"GITHUB_REF_NAME":      "develop",
	}))

	err := c.Check()
	if err == nil || !strings.Contains(err.Error(), "branch mismatch") {
		t.Fatalf("expected branch mismatch, got: %v", err)
	}
}

func TestCheck_MultipleErrors(t *testing.T) {
	t.Parallel()

	err := NewWithEnv(env(map[string]string{})).Check()
	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), "GITHUB_ACTIONS") || !strings.Contains(err.Error(), "GITHUB_TOKEN") {
		t.Fatalf("expected combined errors, got: %v", err)
	}
}
