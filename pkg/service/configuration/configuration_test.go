package configuration_test

import (
	"context"
	"strings"
	"testing"

	"github.com/spotahome/service-level-operator/pkg/service/configuration"
	"github.com/stretchr/testify/assert"
)

func TestJSONLoaderLoadDefaultSLISource(t *testing.T) {
	tests := map[string]struct {
		jsonConfig string
		expConfig  *configuration.DefaultSLISource
		expErr     bool
	}{
		"Correct JSON configuration should be loaded without error.": {
			jsonConfig: `{"prometheus": {"address": "http://test:9090"}}`,
			expConfig: &configuration.DefaultSLISource{
				Prometheus: configuration.PrometheusSLISource{
					Address: "http://test:9090",
				},
			},
		},

		"A malformed JSON should error.": {
			jsonConfig: `{"prometheus":`,
			expErr:     true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			r := strings.NewReader(test.jsonConfig)
			gotConfig, err := configuration.JSONLoader{}.LoadDefaultSLISource(context.TODO(), r)
			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expConfig, gotConfig)
			}
		})
	}
}
