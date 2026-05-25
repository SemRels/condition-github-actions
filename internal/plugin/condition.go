// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The SemRels Authors

// Package plugin implements the CIConditionPlugin for GitHub Actions.
// It verifies that the release pipeline is running inside an authorised
// GitHub Actions workflow with the required environment variables present.
package plugin

import (
	"context"
	"fmt"
	"os"
	"strings"

	semrelv1 "github.com/SemRels/semrel-api/api/gen/v1"
)

// Condition implements semrelv1.CIConditionPluginServer for GitHub Actions.
type Condition struct {
	semrelv1.UnimplementedCIConditionPluginServer

	// env is the OS environment lookup function; override in tests.
	env func(string) string
}

// New returns a new Condition that reads real OS environment variables.
func New() *Condition {
	return &Condition{env: os.Getenv}
}

// NewWithEnv returns a Condition backed by the supplied env function (for testing).
func NewWithEnv(env func(string) string) *Condition {
	return &Condition{env: env}
}

// VerifyConditions checks:
//  1. GITHUB_ACTIONS == "true"
//  2. GITHUB_TOKEN or GH_TOKEN is non-empty
//  3. If ctx.Branch is set, it must match GITHUB_REF_NAME
func (c *Condition) VerifyConditions(_ context.Context, req *semrelv1.VerifyConditionsRequest) (*semrelv1.VerifyConditionsResponse, error) {
	var errs []string

	if c.env("GITHUB_ACTIONS") != "true" {
		errs = append(errs, "GITHUB_ACTIONS is not set to \"true\"; this plugin requires a GitHub Actions environment")
	}

	if c.env("GITHUB_TOKEN") == "" && c.env("GH_TOKEN") == "" {
		errs = append(errs, "neither GITHUB_TOKEN nor GH_TOKEN is set; a GitHub token is required to create releases")
	}

	if req.GetCtx() != nil && req.GetCtx().GetBranch() != "" {
		wantBranch := req.GetCtx().GetBranch()
		gotBranch := c.env("GITHUB_REF_NAME")
		if gotBranch == "" {
			gotBranch = strings.TrimPrefix(c.env("GITHUB_REF"), "refs/heads/")
		}

		if gotBranch != wantBranch {
			errs = append(errs, fmt.Sprintf(
				"branch mismatch: configured branch is %q but GITHUB_REF_NAME is %q",
				wantBranch, gotBranch,
			))
		}
	}

	if len(errs) > 0 {
		return &semrelv1.VerifyConditionsResponse{
			ErrorMessage: strings.Join(errs, "; "),
			Details:      "Ensure this workflow runs on the correct branch with a GITHUB_TOKEN or GH_TOKEN secret.",
		}, nil
	}

	return &semrelv1.VerifyConditionsResponse{}, nil
}
