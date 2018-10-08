/*
Package mocks will have all the mocks of the application, we'll try to use mocking using blackbox
testing and integration tests whenever is possible.
*/
package mocks // import "github.com/slok/service-level-operator/mocks"

// Service mocks.
//go:generate mockery -output ./service/sli -outpkg sli -dir ../pkg/service/sli -name Retriever
//go:generate mockery -output ./service/slo -outpkg slo -dir ../pkg/service/slo -name Output
