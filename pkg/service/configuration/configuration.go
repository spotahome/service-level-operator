package configuration

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
)

// DefaultSLISource is a configuration object with the default
// endpoints.
type DefaultSLISource struct {
	Prometheus PrometheusSLISource `json:"prometheus,omitempty"`
}

// PrometheusSLISource is the default prometheus source.
type PrometheusSLISource struct {
	Address string `json:"address,omitempty"`
}


// Loader knows how to load configuration based on different formats.
// At this moment configuration is not versioned, the configuration
// is so simple that if it grows we could refactor and add version,
// in this case not versioned configuration could be loaded as v1.
type Loader interface {
	// LoadDefaultSLISource will load the default sli source configuration .
	LoadDefaultSLISource(ctx context.Context, r io.Reader) (*DefaultSLISource, error)
}

// JSONLoader knows how to load application configuration.
type JSONLoader struct{}

// LoadDefaultSLISource satisfies Loader interface by loading in JSON format.
func (j JSONLoader) LoadDefaultSLISource(_ context.Context, r io.Reader) (*DefaultSLISource, error) {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	cfg := &DefaultSLISource{}
	err = json.Unmarshal(bs, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
