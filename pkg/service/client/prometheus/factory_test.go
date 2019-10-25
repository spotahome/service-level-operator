package prometheus_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/spotahome/service-level-operator/pkg/service/client/prometheus"
)

func TestBaseFactoryV1Client(t *testing.T) {
	tests := map[string]struct {
		cli     func() *prometheus.BaseFactory
		address string
		expErr  bool
	}{
		"A regular client address should be returned without error.": {
			cli: func() *prometheus.BaseFactory {
				return prometheus.NewBaseFactory()
			},
			address: "http://127.0.0.1:9090",
		},

		"Getting a missing address client should error.": {
			cli: func() *prometheus.BaseFactory {
				return prometheus.NewBaseFactory()
			},
			address: "",
			expErr:  true,
		},

		"Getting a missing address client with a default client it should not error.": {
			cli: func() *prometheus.BaseFactory {
				f := prometheus.NewBaseFactory()
				f.WithDefaultV1APIClient("http://127.0.0.1:9090")
				return f
			},
			address: "",
			expErr:  false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			f := test.cli()
			_, err := f.GetV1APIClient(test.address)

			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}

}
