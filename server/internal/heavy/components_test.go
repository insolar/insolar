// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build slowtest

package heavy

import (
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"

	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/applicationbase/genesis"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/heavy/exporter"
)

func TestComponents(t *testing.T) {
	tests := []struct {
		name         string
		exporterAuth bool
		authSecret   string
		errorCheck   func(*testing.T, error)
	}{{"no exporter auth",
		false,
		"",
		func(t *testing.T, err error) { require.NoError(t, err) },
	}, {"exporter auth - 512-bit secret",
		true,
		"/B?E(H+MbQeThWmZq4t7w!z$C&F)J@NcRfUjXn2r5u8x/A?D*G-KaPdSgVkYp3s6",
		func(t *testing.T, err error) { require.NoError(t, err) },
	}, {"exporter auth - short secret",
		true,
		"E(G+KbPeShVmYq3t6w9z$C&F)J@McQfT",
		func(t *testing.T, err error) { require.Error(t, err) },
	}, {"exporter auth - long secret",
		true,
		"B?E(H+MbQeThVmYq3t6w9z$C&F)J@NcRfUjXnZr4u7x!A%D*G-KaPdSgVkYp3s5v8y/B?E(H+MbQeThWmZq4t7w9z$C&F)J@NcRfUjXn2r5u8x/A%D*G-KaPdSgVkYp3",
		func(t *testing.T, err error) { require.Error(t, err) },
	}, {"exporter auth - empty secret",
		true,
		"",
		func(t *testing.T, err error) { require.Error(t, err) },
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tmpdir, err := ioutil.TempDir("", "heavy-")
			defer os.RemoveAll(tmpdir)
			require.NoError(t, err)

			ctx := inslogger.UpdateLogger(context.Background(), func(logger insolar.Logger) (insolar.Logger, error) {
				return logger.Copy().WithBuffer(100, false).Build()
			})
			cfg := configuration.NewHeavyBadgerConfig()
			cfg.KeysPath = "testdata/bootstrap_keys.json"
			cfg.CertificatePath = "testdata/certificate.json"
			cfg.Metrics.ListenAddress = "0.0.0.0:0"
			cfg.APIRunner.Address = "0.0.0.0:0"
			cfg.AdminAPIRunner.Address = "0.0.0.0:0"
			cfg.APIRunner.SwaggerPath = "../../../api/testdata/api-exported.yaml"
			cfg.AdminAPIRunner.SwaggerPath = "../../../api/testdata/api-exported.yaml"
			cfg.Ledger.Storage.DataDirectory = tmpdir
			cfg.Exporter.Addr = ":0"
			cfg.Exporter.Auth.Required = test.exporterAuth
			cfg.Exporter.Auth.Secret = test.authSecret

			holder := &configuration.HeavyBadgerHolder{
				Configuration: &cfg,
			}

			_, err = newComponents(ctx, holder, genesis.HeavyConfig{Skip: true}, genesis.Options{}, false, api.Options{})
			test.errorCheck(t, err)
		})
	}
}

func TestAuthorize(t *testing.T) {
	t.Run("no metadata", func(t *testing.T) {
		newCtx, err := authorize(context.Background())
		require.Error(t, err)
		require.Contains(t, err.Error(), "code = InvalidArgument desc = failed to retrieve metadata")
		require.Nil(t, newCtx)
	})

	t.Run("no authorization token", func(t *testing.T) {
		data := map[string]string{}
		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		newCtx, err := authorize(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "code = InvalidArgument desc = auth data not supplied")
		require.Nil(t, newCtx)
	})

	t.Run("no key and issuer", func(t *testing.T) {
		data := map[string]string{
			"authorization": "Bearer ",
		}
		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		newCtx, err := authorize(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "required auth parameters are not configured")
		require.Nil(t, newCtx)
	})

	sub := "test-subject"
	issuer := "test-issuer"
	secret := "/B?E(H+MbQeThWmZq4t7w!z$C&F)J@NcRfUjXn2r5u8x/A?D*G-KaPdSgVkYp3s6"

	t.Run("empty token", func(t *testing.T) {
		// prepare configuration
		server, err := newGRPCServer(configuration.Exporter{
			Auth: configuration.Auth{
				Required: true,
				Issuer:   issuer,
				Secret:   secret,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, server)

		data := map[string]string{
			"authorization": "",
		}
		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		newCtx, err := authorize(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "code = Unauthenticated desc = jwt: malformed token")
		require.Nil(t, newCtx)
	})

	t.Run("signature verification failed", func(t *testing.T) {
		// prepare configuration
		server, err := newGRPCServer(configuration.Exporter{
			Auth: configuration.Auth{
				Required: true,
				Issuer:   issuer,
				Secret:   secret,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, server)

		// prepare JWT
		now := time.Now()
		pl := jwt.Payload{
			Issuer:         issuer,
			Subject:        sub,
			ExpirationTime: jwt.NumericDate(now.Add(time.Hour)),
			IssuedAt:       jwt.NumericDate(now),
		}
		hs := jwt.NewHS512([]byte("1111111111111111111111111111111111111111111111111111111111111111"))
		token, err := jwt.Sign(pl, hs)
		require.NoError(t, err)
		require.NotNil(t, token)

		// prepare initial MD
		data := map[string]string{
			"authorization": "Bearer " + string(token),
		}
		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		// test
		newCtx, err := authorize(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "code = Unauthenticated desc = jwt: HMAC verification failed")
		require.Nil(t, newCtx)
	})

	t.Run("unknown issuer", func(t *testing.T) {
		// prepare configuration
		server, err := newGRPCServer(configuration.Exporter{
			Auth: configuration.Auth{
				Required: true,
				Issuer:   issuer,
				Secret:   secret,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, server)

		// prepare JWT
		now := time.Now()
		pl := jwt.Payload{
			Issuer:         "unknown-issuer",
			Subject:        sub,
			ExpirationTime: jwt.NumericDate(now.Add(time.Hour)),
			IssuedAt:       jwt.NumericDate(now),
		}
		hs := jwt.NewHS512([]byte(secret))
		token, err := jwt.Sign(pl, hs)
		require.NoError(t, err)
		require.NotNil(t, token)

		// prepare initial MD
		data := map[string]string{
			"authorization": "Bearer " + string(token),
		}
		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		// test
		newCtx, err := authorize(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "code = Unauthenticated desc = unknown JWT issuer: unknown-issuer")
		require.Nil(t, newCtx)
	})

	t.Run("expired JWT", func(t *testing.T) {
		// prepare configuration
		server, err := newGRPCServer(configuration.Exporter{
			Auth: configuration.Auth{
				Required: true,
				Issuer:   issuer,
				Secret:   secret,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, server)

		// prepare JWT
		now := time.Now()
		pl := jwt.Payload{
			Issuer:         issuer,
			Subject:        sub,
			ExpirationTime: jwt.NumericDate(now.Add(-time.Second)),
			IssuedAt:       jwt.NumericDate(now.Add(-time.Minute)),
		}
		hs := jwt.NewHS512([]byte(secret))
		token, err := jwt.Sign(pl, hs)
		require.NoError(t, err)
		require.NotNil(t, token)

		// prepare initial MD
		data := map[string]string{
			"authorization": "Bearer " + string(token),
		}
		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		// test
		newCtx, err := authorize(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "code = Unauthenticated desc = JWT is expired")
		require.Nil(t, newCtx)
	})

	t.Run("success auth", func(t *testing.T) {
		// prepare configuration
		server, err := newGRPCServer(configuration.Exporter{
			Auth: configuration.Auth{
				Required: true,
				Issuer:   issuer,
				Secret:   secret,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, server)

		// prepare JWT
		now := time.Now()
		pl := jwt.Payload{
			Issuer:         issuer,
			Subject:        sub,
			ExpirationTime: jwt.NumericDate(now.Add(time.Hour)),
			NotBefore:      jwt.NumericDate(now),
			IssuedAt:       jwt.NumericDate(now),
		}
		hs := jwt.NewHS512([]byte(secret))
		token, err := jwt.Sign(pl, hs)
		require.NoError(t, err)
		require.NotNil(t, token)

		// prepare initial MD
		data := map[string]string{
			"authorization": "Bearer " + string(token),
		}
		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		// test
		newCtx, err := authorize(ctx)
		require.NoError(t, err)
		md, ok := metadata.FromIncomingContext(newCtx)
		require.True(t, ok)
		id := md.Get(exporter.ObsID)
		require.Len(t, id, 1, "there is no '%s' in the MD", exporter.ObsID)
		require.Equal(t, sub, id[0])
	})
}
