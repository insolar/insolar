package artifactmanager

import (
	"context"
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/stretchr/testify/require"
)

func TestMiddleware_waitForDrop(t *testing.T){
	t.Run("jetDropTimeout is nil", func(t *testing.T) {
		middleware := middleware{jetDropTimeoutProvider:jetDropTimeoutProvider{}}
		expectedParcel := message.Parcel{PulseNumber:8888}
		handler := func(context context.Context, parcel core.Parcel) (reply core.Reply, e error) {
				require.Equal(t, &expectedParcel, parcel)
				return nil, nil
		}

		internal := middleware.waitForDrop(handler)
		rep, err := internal(nil, &expectedParcel)

		require.Nil(t, rep)
		require.Nil(t, err)

	})

}