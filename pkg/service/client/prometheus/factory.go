package prometheus

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/api"
	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

// ClientFactory knows how to get prometheus API clients.
type ClientFactory interface {
	// GetV1APIClient returns a new prometheus v1 API client.
	// address is the address of the prometheus.
	GetV1APIClient(address string) (promv1.API, error)
}

// BaseFactory returns Prometheus clients based on the address.
// This factory implements a way of returning default Prometheus
// clients in case it was set.
type BaseFactory struct {
	v1Clis map[string]api.Client
	climu  sync.Mutex
}

// NewBaseFactory returns a new client factory.
// If default address is passed when an empty address
// client is requested it will return the client based
// on this address
func NewBaseFactory() *BaseFactory {
	return &BaseFactory{
		v1Clis: map[string]api.Client{},
	}
}

// GetV1APIClient satisfies ClientFactory interface.
func (f *BaseFactory) GetV1APIClient(address string) (promv1.API, error) {
	f.climu.Lock()
	defer f.climu.Unlock()

	var err error
	cli, ok := f.v1Clis[address]
	if !ok {
		cli, err = newClient(address)
		if err != nil {
			return nil, fmt.Errorf("error creating prometheus client: %s", err)
		}
		f.v1Clis[address] = cli
	}
	return promv1.NewAPI(cli), nil
}

// WithDefaultV1APIClient sets a default client for V1 api client.
func (f *BaseFactory) WithDefaultV1APIClient(address string) error {
	const defAddressKey = ""

	dc, err := newClient(address)
	if err != nil {
		return fmt.Errorf("error creating prometheus client: %s", err)
	}
	f.v1Clis[defAddressKey] = dc

	return nil
}

func newClient(address string) (api.Client, error) {
	if address == "" {
		return nil, fmt.Errorf("address can't be empty")
	}

	return api.NewClient(api.Config{Address: address})
}

// MockFactory returns a predefined prometheus v1 API client.
type MockFactory struct {
	Cli promv1.API
}

// GetV1APIClient satisfies ClientFactory interface.
func (m *MockFactory) GetV1APIClient(_ string) (promv1.API, error) {
	return m.Cli, nil
}
