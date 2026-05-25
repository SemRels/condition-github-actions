// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The go-semrel Authors

// Plugin binary: condition-github-actions
//
// This binary is launched by the go-semrel host via hashicorp/go-plugin.
// stdout is reserved for the go-plugin handshake; all logging goes to stderr.
package main

import (
	"log"
	"os"

	hcplugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	semrelv1 "github.com/SemRels/condition-github-actions/api/gen/v1"
	grpcserver "github.com/SemRels/condition-github-actions/internal/grpc"
	"github.com/SemRels/condition-github-actions/internal/plugin"
)

// HandshakeConfig is the magic-cookie configuration shared between the host and
// this plugin binary.  Both sides must present the same values.
var HandshakeConfig = hcplugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "SEMREL_PLUGIN",
	MagicCookieValue: "semrel",
}

// ciConditionGRPCPlugin bridges hashicorp/go-plugin with the generated gRPC service.
type ciConditionGRPCPlugin struct {
	hcplugin.Plugin
	impl *grpcserver.CIConditionServer
}

func (p *ciConditionGRPCPlugin) GRPCServer(_ *hcplugin.GRPCBroker, s *grpc.Server) error {
	semrelv1.RegisterCIConditionPluginServer(s, p.impl)
	return nil
}

func (p *ciConditionGRPCPlugin) GRPCClient(_ interface{}, _ *hcplugin.GRPCBroker, _ *grpc.ClientConn) (interface{}, error) {
	// The host side constructs the client; the plugin binary never needs this.
	return nil, nil
}

func main() {
	logger := log.New(os.Stderr, "[condition-github-actions] ", log.LstdFlags)

	checker := plugin.NewConditionChecker(&plugin.OSEnvProvider{})
	srv := grpcserver.NewCIConditionServer(checker)

	logger.Println("starting gRPC plugin server")

	hcplugin.Serve(&hcplugin.ServeConfig{
		HandshakeConfig: HandshakeConfig,
		Plugins: map[string]hcplugin.Plugin{
			"ci-condition": &ciConditionGRPCPlugin{impl: srv},
		},
		GRPCServer: hcplugin.DefaultGRPCServer,
		Logger:     nil, // go-plugin uses hclog; we use stderr via log package above
	})
}
