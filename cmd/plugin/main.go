// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The SemRels Authors

package main

import (
	"fmt"
	"io"
	"os"

	plugin "github.com/SemRels/condition-github-actions/internal/plugin"
)

func run(getenv func(string) string, stderr io.Writer) int {
	c := plugin.NewWithEnv(getenv)
	if err := c.Check(); err != nil {
		fmt.Fprintln(stderr, "condition-github-actions:", err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run(os.Getenv, os.Stderr))
}
