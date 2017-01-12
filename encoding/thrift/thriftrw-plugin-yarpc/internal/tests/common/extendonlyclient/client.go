// Code generated by thriftrw-plugin-yarpc
// @generated

package extendonlyclient

import (
	"go.uber.org/yarpc"
	"go.uber.org/yarpc/api/transport"
	"go.uber.org/yarpc/encoding/thrift"
	"go.uber.org/yarpc/encoding/thrift/thriftrw-plugin-yarpc/internal/tests/common/baseserviceclient"
)

// Interface is a client for the ExtendOnly service.
type Interface interface {
	baseserviceclient.Interface
}

// New builds a new client for the ExtendOnly service.
//
// 	client := extendonlyclient.New(dispatcher.ClientConfig("extendonly"))
func New(c transport.ClientConfig, opts ...thrift.ClientOption) Interface {
	return client{
		c: thrift.New(thrift.Config{
			Service:      "ExtendOnly",
			ClientConfig: c,
		}, opts...),
		Interface: baseserviceclient.New(c),
	}
}

func init() {
	yarpc.RegisterClientBuilder(func(c transport.ClientConfig) Interface {
		return New(c)
	})
}

type client struct {
	baseserviceclient.Interface

	c thrift.Client
}
