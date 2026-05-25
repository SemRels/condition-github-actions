// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The SemRels Authors

package plugin

import (
	"fmt"
	"os"
	"strings"
)

type Condition struct {
	env func(string) string
}

func New() *Condition { return &Condition{env: os.Getenv} }

func NewWithEnv(env func(string) string) *Condition { return &Condition{env: env} }

func (c *Condition) Check() error {
	var errs []string

	if c.env("GITHUB_ACTIONS") != "true" {
		errs = append(errs, "GITHUB_ACTIONS is not set to \"true\"; this plugin requires a GitHub Actions environment")
	}

	if c.env("GITHUB_TOKEN") == "" && c.env("GH_TOKEN") == "" {
		errs = append(errs, "neither GITHUB_TOKEN nor GH_TOKEN is set")
	}

	if branch := c.env("SEMREL_PLUGIN_BRANCH"); branch != "" {
		gotBranch := c.env("GITHUB_REF_NAME")
		if gotBranch == "" {
			gotBranch = strings.TrimPrefix(c.env("GITHUB_REF"), "refs/heads/")
		}
		if gotBranch != branch {
			errs = append(errs, fmt.Sprintf("branch mismatch: want %q got %q", branch, gotBranch))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, "; "))
	}
	return nil
}
