package mocks

import (
	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

// Third party resources to create mocks from their interfaces.

// API is the interface promv1.API.
type API interface{ promv1.API }
