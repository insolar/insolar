package heavy

import (
	"context"
	"testing"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/ledger/heavy/exporter"
)

func TestValidateVersionHeavyVersion(t *testing.T) {

	sub := "test-subject"
	issuer := "test-issuer"
	secret := "/B?E(H+MbQeThWmZq4t7w!z$C&F)J@NcRfUjXn2r5u8x/A?D*G-KaPdSgVkYp3s6"

	t.Run("success validate version", func(t *testing.T) {
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
			"authorization":                "Bearer " + string(token),
			exporter.KeyClientType:         exporter.ValidateHeavyVersion.String(),
			exporter.KeyClientVersionHeavy: "2",
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

	t.Run("failed without type client", func(t *testing.T) {
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
		_, err = authorize(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown type client")
	})

	t.Run("failed unknown type client", func(t *testing.T) {
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
			"authorization":        "Bearer " + string(token),
			exporter.KeyClientType: "unknown version",
		}

		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		// test
		_, err = authorize(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown type client")
	})

	t.Run("failed without heavy version", func(t *testing.T) {
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
			"authorization":        "Bearer " + string(token),
			exporter.KeyClientType: exporter.ValidateHeavyVersion.String(),
		}

		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		// test
		_, err = authorize(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown heavy_version")
	})

	t.Run("failed incorrect format heavy version", func(t *testing.T) {
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
			"authorization":                "Bearer " + string(token),
			exporter.KeyClientType:         exporter.ValidateHeavyVersion.String(),
			exporter.KeyClientVersionHeavy: "unknown version",
		}

		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		// test
		_, err = authorize(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "incorrect format of the heavy_version")
	})

	t.Run("failed deprecated heavy version", func(t *testing.T) {
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
			"authorization":                "Bearer " + string(token),
			exporter.KeyClientType:         exporter.ValidateHeavyVersion.String(),
			exporter.KeyClientVersionHeavy: "1",
		}

		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		// test
		_, err = authorize(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "version of the observer is outdated. Please upgrade this client")
	})

	t.Run("success new version observer", func(t *testing.T) {
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
			"authorization":                "Bearer " + string(token),
			exporter.KeyClientType:         exporter.ValidateHeavyVersion.String(),
			exporter.KeyClientVersionHeavy: "3",
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

	t.Run("failed incorrect format heavy version less zero", func(t *testing.T) {
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
			"authorization":                "Bearer " + string(token),
			exporter.KeyClientType:         exporter.ValidateHeavyVersion.String(),
			exporter.KeyClientVersionHeavy: "-1",
		}

		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		// test
		_, err = authorize(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "incorrect format of the heavy_version")
	})

}

func TestValidateVersionContractVersion(t *testing.T) {

	sub := "test-subject"
	issuer := "test-issuer"
	secret := "/B?E(H+MbQeThWmZq4t7w!z$C&F)J@NcRfUjXn2r5u8x/A?D*G-KaPdSgVkYp3s6"
	allowedVersionContract = 2
	t.Run("success validate version", func(t *testing.T) {
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
			"authorization":                   "Bearer " + string(token),
			exporter.KeyClientType:            exporter.ValidateContractVersion.String(),
			exporter.KeyClientVersionHeavy:    "2",
			exporter.KeyClientVersionContract: "2",
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

	t.Run("failed without contract version", func(t *testing.T) {
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
			"authorization":                "Bearer " + string(token),
			exporter.KeyClientType:         exporter.ValidateContractVersion.String(),
			exporter.KeyClientVersionHeavy: "2",
		}

		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		// test
		_, err = authorize(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown contract_version")
	})

	t.Run("failed incorrect format heavy version", func(t *testing.T) {
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
			"authorization":                   "Bearer " + string(token),
			exporter.KeyClientType:            exporter.ValidateContractVersion.String(),
			exporter.KeyClientVersionHeavy:    "2",
			exporter.KeyClientVersionContract: "unknown version",
		}

		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		// test
		_, err = authorize(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "incorrect format of the contract_version")
	})

	t.Run("failed deprecated contract version", func(t *testing.T) {
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
			"authorization":                   "Bearer " + string(token),
			exporter.KeyClientType:            exporter.ValidateContractVersion.String(),
			exporter.KeyClientVersionHeavy:    "2",
			exporter.KeyClientVersionContract: "1",
		}

		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		// test
		_, err = authorize(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "version of the observer is outdated. Please upgrade this client")
	})

	t.Run("success new version observer", func(t *testing.T) {
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
			"authorization":                   "Bearer " + string(token),
			exporter.KeyClientType:            exporter.ValidateContractVersion.String(),
			exporter.KeyClientVersionHeavy:    "2",
			exporter.KeyClientVersionContract: "3",
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

	t.Run("failed incorrect format heavy version less zero", func(t *testing.T) {
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
			"authorization":                   "Bearer " + string(token),
			exporter.KeyClientType:            exporter.ValidateContractVersion.String(),
			exporter.KeyClientVersionHeavy:    "2",
			exporter.KeyClientVersionContract: "-1",
		}

		initialMD := metadata.New(data)
		ctx := metadata.NewIncomingContext(context.Background(), initialMD)

		// test
		_, err = authorize(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "incorrect format of the contract_version")
	})

}
