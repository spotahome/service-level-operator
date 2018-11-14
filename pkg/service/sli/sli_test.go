package sli_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/spotahome/service-level-operator/pkg/service/sli"
)

func TestSLIResult(t *testing.T) {
	tests := []struct {
		name              string
		errorQ            float64
		totalQ            float64
		expAvailability   float64
		expAvailabiityErr bool
		expError          float64
		expErrorErr       bool
	}{
		{
			name:            "Not having a total quantity should return everything ok.",
			expAvailability: 1,
			expError:        0,
		},
		{
			name:              "Having more errors than total should be impossible.",
			errorQ:            600,
			totalQ:            300,
			expErrorErr:       true,
			expAvailabiityErr: true,
		},
		{
			name:            "If half of the total are errors then the ratio of availability and error should be 0.5.",
			errorQ:          300,
			totalQ:          600,
			expAvailability: 0.5,
			expError:        0.5,
		},
		{
			name:            "If a 33% of errors then the ratios should be 0.33 and 0.66.",
			errorQ:          33,
			totalQ:          100,
			expAvailability: 0.6699999999999999,
			expError:        0.33,
		},
		{
			name:            "In small quantities the ratios should be correctly calculated.",
			errorQ:          4,
			totalQ:          10,
			expAvailability: 0.6,
			expError:        0.4,
		},
		{
			name:            "In big quantities the ratios should be correctly calculated.",
			errorQ:          240,
			totalQ:          10000000,
			expAvailability: 0.999976,
			expError:        0.000024,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			res := sli.Result{
				TotalQ: test.totalQ,
				ErrorQ: test.errorQ,
			}

			av, err := res.AvailabilityRatio()
			if test.expAvailabiityErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expAvailability, av)
			}

			dw, err := res.ErrorRatio()
			if test.expErrorErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expError, dw)
			}
		})
	}
}
