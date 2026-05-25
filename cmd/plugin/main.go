// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The SemRels Authors

package main

import (
	semrelapi "github.com/SemRels/semrel-api/plugin"
	conditionplugin "github.com/SemRels/condition-github-actions/internal/plugin"
	"github.com/hashicorp/go-plugin"
)

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: semrelapi.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"condition": &semrelapi.CIConditionGRPCPlugin{
				Impl: conditionplugin.New(),
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
