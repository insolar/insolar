// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package heavy

import (
	"context"
	"testing"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/ledger/heavy/exporter"
)

func TestValidateVersionHeavyVersion(t *testing.T) {

	t.Run("success validate version", func(t *testing.T) {
		// prepare configuration
		server, err := newGRPCServer(configuration.Exporter{
			CheckVersion: true,
		}, grpc_prometheus.NewServerMetrics())
		require.NoError(t, err)
		require.NotNil(t, server)

		// prepare initial MD
		data := map[string]string{
			exporter.KeyClientType:         exporter.ValidateHeavyVersion.String(),
			exporter.KeyClientVersionHeavy: "2",
		}

		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		// test
		err = validateClientVersion(ctx)
		require.NoError(t, err)
	})

	t.Run("failed without type client", func(t *testing.T) {
		// prepare configuration
		server, err := newGRPCServer(configuration.Exporter{
			CheckVersion: true,
		}, grpc_prometheus.NewServerMetrics())
		require.NoError(t, err)
		require.NotNil(t, server)

		// prepare initial MD
		data := map[string]string{}

		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		// test
		err = validateClientVersion(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown type client")
	})

	t.Run("failed unknown type client", func(t *testing.T) {
		// prepare configuration
		server, err := newGRPCServer(configuration.Exporter{
			CheckVersion: true,
		}, grpc_prometheus.NewServerMetrics())
		require.NoError(t, err)
		require.NotNil(t, server)

		// prepare initial MD
		data := map[string]string{
			exporter.KeyClientType: "unknown version",
		}

		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		// test
		err = validateClientVersion(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown type client")
	})

	t.Run("failed without heavy version", func(t *testing.T) {
		// prepare configuration
		server, err := newGRPCServer(configuration.Exporter{
			CheckVersion: true,
		}, grpc_prometheus.NewServerMetrics())
		require.NoError(t, err)
		require.NotNil(t, server)

		// prepare initial MD
		data := map[string]string{
			exporter.KeyClientType: exporter.ValidateHeavyVersion.String(),
		}

		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		// test
		err = validateClientVersion(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown heavy-version")
	})

	t.Run("failed incorrect format heavy version", func(t *testing.T) {
		// prepare configuration
		server, err := newGRPCServer(configuration.Exporter{
			CheckVersion: true,
		}, grpc_prometheus.NewServerMetrics())
		require.NoError(t, err)
		require.NotNil(t, server)

		// prepare initial MD
		data := map[string]string{
			exporter.KeyClientType:         exporter.ValidateHeavyVersion.String(),
			exporter.KeyClientVersionHeavy: "unknown version",
		}

		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		// test
		err = validateClientVersion(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "incorrect format of the heavy-version")
	})

	t.Run("failed deprecated heavy version", func(t *testing.T) {
		// prepare configuration
		server, err := newGRPCServer(configuration.Exporter{
			CheckVersion: true,
		}, grpc_prometheus.NewServerMetrics())
		require.NoError(t, err)
		require.NotNil(t, server)

		// prepare initial MD
		data := map[string]string{
			exporter.KeyClientType:         exporter.ValidateHeavyVersion.String(),
			exporter.KeyClientVersionHeavy: "1",
		}

		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		// test
		err = validateClientVersion(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), exporter.ErrDeprecatedClientVersion.Error())
	})

	t.Run("success new version observer", func(t *testing.T) {
		// prepare configuration
		server, err := newGRPCServer(configuration.Exporter{
			CheckVersion: true,
		}, grpc_prometheus.NewServerMetrics())
		require.NoError(t, err)
		require.NotNil(t, server)

		// prepare initial MD
		data := map[string]string{
			exporter.KeyClientType:         exporter.ValidateHeavyVersion.String(),
			exporter.KeyClientVersionHeavy: "3",
		}

		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		// test
		err = validateClientVersion(ctx)
		require.NoError(t, err)
	})

	t.Run("failed incorrect format heavy version less zero", func(t *testing.T) {
		// prepare configuration
		server, err := newGRPCServer(configuration.Exporter{

			CheckVersion: true,
		}, grpc_prometheus.NewServerMetrics())
		require.NoError(t, err)
		require.NotNil(t, server)

		// prepare initial MD
		data := map[string]string{

			exporter.KeyClientType:         exporter.ValidateHeavyVersion.String(),
			exporter.KeyClientVersionHeavy: "-1",
		}

		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		// test
		err = validateClientVersion(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "incorrect format of the heavy-version")
	})

}

func TestValidateVersionContractVersion(t *testing.T) {
	allowedVersionContract = 2

	t.Run("success validate version", func(t *testing.T) {
		// prepare configuration
		server, err := newGRPCServer(configuration.Exporter{
			CheckVersion: true,
		}, grpc_prometheus.NewServerMetrics())
		require.NoError(t, err)
		require.NotNil(t, server)

		// prepare initial MD
		data := map[string]string{
			exporter.KeyClientType:            exporter.ValidateContractVersion.String(),
			exporter.KeyClientVersionHeavy:    "2",
			exporter.KeyClientVersionContract: "2",
		}

		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		// test
		err = validateClientVersion(ctx)
		require.NoError(t, err)
	})

	t.Run("failed without contract version", func(t *testing.T) {
		// prepare configuration
		server, err := newGRPCServer(configuration.Exporter{
			CheckVersion: true,
		}, grpc_prometheus.NewServerMetrics())
		require.NoError(t, err)
		require.NotNil(t, server)

		// prepare initial MD
		data := map[string]string{
			exporter.KeyClientType:         exporter.ValidateContractVersion.String(),
			exporter.KeyClientVersionHeavy: "2",
		}

		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		// test
		err = validateClientVersion(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown contract-version")
	})

	t.Run("failed incorrect format heavy version", func(t *testing.T) {
		// prepare configuration
		server, err := newGRPCServer(configuration.Exporter{
			CheckVersion: true,
		}, grpc_prometheus.NewServerMetrics())
		require.NoError(t, err)
		require.NotNil(t, server)

		// prepare initial MD
		data := map[string]string{
			exporter.KeyClientType:            exporter.ValidateContractVersion.String(),
			exporter.KeyClientVersionHeavy:    "2",
			exporter.KeyClientVersionContract: "unknown version",
		}

		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		// test
		err = validateClientVersion(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "incorrect format of the contract-version")
	})

	t.Run("failed deprecated contract version", func(t *testing.T) {
		// prepare configuration
		server, err := newGRPCServer(configuration.Exporter{
			CheckVersion: true,
		}, grpc_prometheus.NewServerMetrics())
		require.NoError(t, err)
		require.NotNil(t, server)

		// prepare initial MD
		data := map[string]string{
			exporter.KeyClientType:            exporter.ValidateContractVersion.String(),
			exporter.KeyClientVersionHeavy:    "2",
			exporter.KeyClientVersionContract: "1",
		}

		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		// test
		err = validateClientVersion(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), exporter.ErrDeprecatedClientVersion.Error())
	})

	t.Run("success new version observer", func(t *testing.T) {
		// prepare configuration
		server, err := newGRPCServer(configuration.Exporter{
			CheckVersion: true,
		}, grpc_prometheus.NewServerMetrics())
		require.NoError(t, err)
		require.NotNil(t, server)

		// prepare initial MD
		data := map[string]string{
			exporter.KeyClientType:            exporter.ValidateContractVersion.String(),
			exporter.KeyClientVersionHeavy:    "2",
			exporter.KeyClientVersionContract: "3",
		}

		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		// test
		err = validateClientVersion(ctx)
		require.NoError(t, err)
	})

	t.Run("failed incorrect format heavy version less zero", func(t *testing.T) {
		// prepare configuration
		server, err := newGRPCServer(configuration.Exporter{
			CheckVersion: true,
		}, grpc_prometheus.NewServerMetrics())
		require.NoError(t, err)
		require.NotNil(t, server)

		// prepare initial MD
		data := map[string]string{
			exporter.KeyClientType:            exporter.ValidateContractVersion.String(),
			exporter.KeyClientVersionHeavy:    "2",
			exporter.KeyClientVersionContract: "-1",
		}

		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		// test
		err = validateClientVersion(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "incorrect format of the contract-version")
	})

}
