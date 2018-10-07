package prometheus

import (
	"context"
	"net/http"
	"net/url"

	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

type fakeFactory struct {
}

// NewFakeFactory returns a new fake factory.
func NewFakeFactory() ClientFactory {
	return &fakeFactory{}
}

// GetV1APIClient satisfies ClientFactory interface.
func (f *fakeFactory) GetV1APIClient(_ string) (promv1.API, error) {
	cli := &fakeClient{}
	return promv1.NewAPI(cli), nil
}

// fakeClient is a faked http client.
// TODO
type fakeClient struct{}

func (c *fakeClient) URL(ep string, args map[string]string) *url.URL { return nil }
func (c *fakeClient) Do(ctx context.Context, req *http.Request) (*http.Response, []byte, error) {
	b := []byte{}
	resp := &http.Response{}
	return resp, b, nil
}
