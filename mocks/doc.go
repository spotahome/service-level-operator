/*
Package mocks will have all the mocks of the application, we'll try to use mocking using blackbox
testing and integration tests whenever is possible.
*/
package mocks // import "github.com/slok/service-level-operator/mocks"

// Service mocks.
//go:generate mockery -output ./service/sli -outpkg sli -dir ../pkg/service/sli -name Retriever
//go:generate mockery -output ./service/output -outpkg slo -dir ../pkg/service/output -name Output

// Third party
//go:generate mockery -output ./github.com/prometheus/client_golang/api/prometheus/v1 -outpkg v1 -dir . -name API
