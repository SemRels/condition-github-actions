// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The go-semrel Authors

// Package grpc wires the ConditionChecker into the generated gRPC service interface.
package grpc

import (
	"context"

	semrelv1 "github.com/SemRels/condition-github-actions/api/gen/v1"
	"github.com/SemRels/condition-github-actions/internal/plugin"
)

// CIConditionServer adapts a ConditionChecker to the generated
// CIConditionPluginServer gRPC interface.
type CIConditionServer struct {
	semrelv1.UnimplementedCIConditionPluginServer
	checker *plugin.ConditionChecker
}

// NewCIConditionServer returns a CIConditionServer backed by the supplied checker.
func NewCIConditionServer(checker *plugin.ConditionChecker) *CIConditionServer {
	return &CIConditionServer{checker: checker}
}

// VerifyConditions implements semrelv1.CIConditionPluginServer.
func (s *CIConditionServer) VerifyConditions(
	_ context.Context,
	req *semrelv1.VerifyConditionsRequest,
) (*semrelv1.VerifyConditionsResponse, error) {
	return s.checker.VerifyConditions(req.GetCtx()), nil
}
