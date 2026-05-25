// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The go-semrel Authors

package grpc

import (
	"context"
	"testing"

	semrelv1 "github.com/SemRels/condition-github-actions/api/gen/v1"
	"github.com/SemRels/condition-github-actions/internal/plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockEnv map[string]string

func (m mockEnv) Getenv(key string) string { return m[key] }

func happyEnv() mockEnv {
	return mockEnv{
		"GITHUB_ACTIONS":    "true",
		"GITHUB_REF_TYPE":   "branch",
		"GITHUB_EVENT_NAME": "push",
		"GITHUB_TOKEN":      "ghs_test",
		"GITHUB_REF_NAME":   "main",
	}
}

func TestVerifyConditions_Success(t *testing.T) {
	t.Parallel()

	checker := plugin.NewConditionChecker(happyEnv())
	srv := NewCIConditionServer(checker)

	resp, err := srv.VerifyConditions(context.Background(), &semrelv1.VerifyConditionsRequest{})
	require.NoError(t, err)
	assert.Empty(t, resp.ErrorMessage)
	assert.Contains(t, resp.Details, "satisfied")
}

func TestVerifyConditions_Failure(t *testing.T) {
	t.Parallel()

	checker := plugin.NewConditionChecker(mockEnv{})
	srv := NewCIConditionServer(checker)

	resp, err := srv.VerifyConditions(context.Background(), &semrelv1.VerifyConditionsRequest{})
	require.NoError(t, err)
	assert.NotEmpty(t, resp.ErrorMessage)
	assert.Contains(t, resp.Details, "Remediation")
}

func TestVerifyConditions_WithContext(t *testing.T) {
	t.Parallel()

	env := happyEnv()
	checker := plugin.NewConditionChecker(env)
	srv := NewCIConditionServer(checker)

	req := &semrelv1.VerifyConditionsRequest{
		Ctx: &semrelv1.ReleaseContext{Branch: "main"},
	}
	resp, err := srv.VerifyConditions(context.Background(), req)
	require.NoError(t, err)
	assert.Empty(t, resp.ErrorMessage)
}
