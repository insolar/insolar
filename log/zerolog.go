package log

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func test() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Debug().Msg("This message appears only when log level set to Debug")

	log.Ctx(context.Background()).Info()
}
