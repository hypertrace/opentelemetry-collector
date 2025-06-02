// Copyright The OpenTelemetry AuthorsAdd commentMore actions
// SPDX-License-Identifier: Apache-2.0

// This is a test added by us.
package configgrpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
)

func TestRegisterClientDialOptionHandler(t *testing.T) {
	tt := componenttest.NewNopTelemetrySettings()

	gcs := &ClientConfig{}
	opts, err := gcs.getGrpcDialOptions(
		context.Background(),
		&mockHost{ext: map[component.ID]component.Component{}},
		tt,
		[]ToClientConnOption{},
	)
	require.NoError(t, err)

	defaultOptsLen := len(opts)

	RegisterClientDialOptionHandlers(func() grpc.DialOption {
		return grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		})
	})
	gcs = &ClientConfig{}
	opts, err = gcs.getGrpcDialOptions(
		context.Background(),
		&mockHost{ext: map[component.ID]component.Component{}},
		tt,
		[]ToClientConnOption{},
	)
	assert.NoError(t, err)
	assert.Len(t, opts, defaultOptsLen+1)
}
