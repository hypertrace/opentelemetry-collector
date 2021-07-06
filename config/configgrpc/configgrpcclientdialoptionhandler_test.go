package configgrpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/obsreport/obsreporttest"
)

type mockHost struct {
	component.Host
	ext map[config.ComponentID]component.Extension
}

func (nh *mockHost) GetExtensions() map[config.ComponentID]component.Extension {
	return nh.ext
}

func TestRegisterClientDialOptionHandler(t *testing.T) {
	tt, err := obsreporttest.SetupTelemetry()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, tt.Shutdown(context.Background())) })

	gcs := &GRPCClientSettings{}
	opts, err := gcs.ToDialOptions(
		&mockHost{ext: map[config.ComponentID]component.Extension{}},
		tt.TelemetrySettings,
	)
	require.NoError(t, err)

	defaultOptsLen := len(opts)

	RegisterClientDialOptionHandlers(func() grpc.DialOption {
		return grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		})
	})
	gcs = &GRPCClientSettings{}
	opts, err = gcs.ToDialOptions(
		&mockHost{ext: map[config.ComponentID]component.Extension{}},
		tt.TelemetrySettings,
	)
	assert.NoError(t, err)
	assert.Len(t, opts, defaultOptsLen+1)
}
