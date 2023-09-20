package introspector

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/introspector/introproto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntrospector_Server(t *testing.T) {
	ctx := inslogger.TestContext(t)

	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err, "listener bind on random port without error")

	srv := NewServer(l.Addr().String(), NewPublisherServerMock(t))
	runErr := srv.run(ctx, l)
	require.NoError(t, runErr, "introspection server run without error")
	defer func() {
		if stopErr := srv.Stop(ctx); stopErr != nil {
			t.Fatal("server stop failed with error:", stopErr)
		}
	}()

	fetcher := defaultClient().POST().Tries(5).Server(l.Addr().String())

	rootResult := fetcher.Url("/").Do()
	assert.Equal(t, 404, rootResult.Code())

	swaggerResult := fetcher.Url("/swagger.json").Do()
	if assert.Equal(t, 200, swaggerResult.Code()) {
		assert.Contains(t, swaggerResult.Body(), `"swagger":`, "swagger json")
	}
}

func TestIntrospector_Server_FilterMessages(t *testing.T) {
	ctx := inslogger.TestContext(t)

	mockState := map[string]struct{}{}
	pubMock := NewPublisherServerMock(t)
	pubMock.SetMessagesFilterMock.Set(func(_ context.Context, in *introproto.MessageFilterByType) (*introproto.MessageFilterByType, error) {
		name, enable := in.Name, in.Enable
		if enable {
			mockState[name] = struct{}{}
		} else {
			delete(mockState, name)
		}
		return in, nil
	})
	pubMock.GetMessagesFiltersMock.Set(func(_ context.Context, _ *introproto.EmptyArgs) (*introproto.AllMessageFilterStats, error) {
		filters := []*introproto.MessageFilterWithStat{}
		for name := range mockState {
			filters = append(filters, &introproto.MessageFilterWithStat{
				Enable: true,
				Name:   name,
			})
		}
		return &introproto.AllMessageFilterStats{Filters: filters}, nil
	})

	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err, "listener bind on random port without error")

	srv := NewServer(l.Addr().String(), pubMock)
	runErr := srv.run(ctx, l)
	require.NoError(t, runErr, "introspection server run without error")
	defer func() {
		if stopErr := srv.Stop(ctx); stopErr != nil {
			t.Fatal("server stop failed with error:", stopErr)
		}
	}()

	name := "TypeCode"

	setUrl := "/setMessagesFilter"
	getUrl := "/getMessagesFilters"

	fetcher := defaultClient().POST().Tries(5).Server(l.Addr().String())

	getEmptyResult := fetcher.Url(getUrl).Do()
	if assert.Equalf(t, 200, getEmptyResult.Code(), "code of %v is fine", getUrl) {
		assert.Equalf(t, `{"Filters":[]}`, getEmptyResult.Body(), "body of %v is fine", getUrl)
	}

	enableJSON := fmt.Sprintf(`{"Name":"%v","Enable":true}`, name)
	setResult := fetcher.Url(setUrl).Body(enableJSON).Do()
	if assert.Equalf(t, 200, setResult.Code(), "code of %v is fine", setUrl) {
		assert.Equalf(t, enableJSON, setResult.Body(), "body of %v is fine", setUrl)
	}

	setResult2 := fetcher.Url(setUrl).Body(enableJSON).Do()
	if assert.Equalf(t, 200, setResult2.Code(), "code of %v is fine", setUrl) {
		assert.Equalf(t, enableJSON, setResult2.Body(), "body of %v is fine", setUrl)
	}

	getResult := fetcher.Url(getUrl).Do()
	if assert.Equalf(t, 200, getResult.Code(), "code of %v is fine", getUrl) {
		expectGetResultJSON := fmt.Sprintf(
			`{"Filters":[{"Name":"%v","Enable":true,"Filtered":"0"}]}`, name)
		assert.Equalf(t, expectGetResultJSON, getResult.Body(), "body of %v is fine", getUrl)
	}
}

type fetchResult struct {
	code int
	body []byte
	err  error
}

func (fr fetchResult) Body() string {
	return string(fr.body)
}

func (fr fetchResult) Code() int {
	return fr.code
}

type fetchClient struct {
	protocol string
	server   string
	url      string
	method   string
	headers  map[string]string
	body     []byte
	tries    int
}

func defaultClient() *fetchClient {
	return &fetchClient{
		protocol: "http",
		server:   "localhost",
		url:      "/",
		method:   "GET",
		headers: map[string]string{
			"Content-Type": "application/json",
		},
		tries: 1,
	}
}

func (tc *fetchClient) Server(srv string) *fetchClient {
	c := *tc
	c.server = srv
	return &c
}

func (tc *fetchClient) Url(u string) *fetchClient {
	c := *tc
	c.url = u
	return &c
}

func (tc *fetchClient) Tries(tries int) *fetchClient {
	c := *tc
	c.tries = tries
	return &c
}

func (tc *fetchClient) POST() *fetchClient {
	c := *tc
	c.method = "POST"
	return &c
}

func (tc *fetchClient) GET() *fetchClient {
	c := *tc
	c.method = "GET"
	return &c
}

func (tc *fetchClient) Body(s string) *fetchClient {
	c := *tc
	c.body = []byte(s)
	return &c
}

func (tc *fetchClient) Do() fetchResult {
	fetchUrl := fmt.Sprintf("%v://%v%v", tc.protocol, tc.server, tc.url)

	var r io.Reader
	if tc.body != nil {
		r = bytes.NewReader(tc.body)
	}
	req, err := http.NewRequest(tc.method, fetchUrl, r)
	if err != nil {
		return fetchResult{err: err}
	}
	for k, v := range tc.headers {
		if v != "" {
			req.Header.Set(k, v)
		}
	}

	var tRes fetchResult
	sleep := time.Millisecond * 200
	for i := 0; i < tc.tries; i++ {
		res, doErr := http.DefaultClient.Do(req)

		tRes.err = doErr
		if res != nil {
			tRes.code = res.StatusCode
			body, _ := ioutil.ReadAll(res.Body)
			tRes.body = body
			_ = res.Body.Close()
		}

		if doErr == nil {
			break
		}

		time.Sleep(sleep)
		if sleep < time.Second {
			sleep = sleep * 2
		}
	}
	return tRes
}
