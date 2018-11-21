// Copyright (c) 2018 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package yarpctchannelfx

import (
	"strings"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/config"
	"go.uber.org/fx/fxtest"
	"go.uber.org/yarpc/v2"
	"go.uber.org/yarpc/v2/yarpcchooser"
	"go.uber.org/yarpc/v2/yarpctchannel"
	"go.uber.org/yarpc/v2/yarpctest"
	"go.uber.org/zap"
)

func newChooserProvider(t *testing.T) yarpc.ChooserProvider {
	p, err := yarpcchooser.NewProvider(yarpctest.NewFakePeerChooser("roundrobin"))
	require.NoError(t, err)
	return p
}

func TestNewInboundConfig(t *testing.T) {
	cfg := strings.NewReader("yarpc: {tchannel: {inbounds: {address: 127.0.0.1:0}}}")
	provider, err := config.NewYAML(config.Source(cfg))
	require.NoError(t, err)

	res, err := NewInboundConfig(InboundConfigParams{
		Provider: provider,
	})
	require.NoError(t, err)
	assert.Equal(t, InboundConfig{Address: "127.0.0.1:0"}, res.Config)
}

func TestStartInbounds(t *testing.T) {
	assert.NoError(t, StartInbounds(StartInboundsParams{
		Lifecycle: fxtest.NewLifecycle(t),
		Router:    yarpctest.NewFakeRouter(nil),
		Config:    InboundConfig{Address: "127.0.0.1:0"},
	}))
}

func TestNewOutboundsConfig(t *testing.T) {
	cfg := strings.NewReader("yarpc: {tchannel: {outbounds: {bar: {address: 127.0.0.1:0}}}}")
	provider, err := config.NewYAML(config.Source(cfg))
	require.NoError(t, err)

	res, err := NewOutboundsConfig(OutboundsConfigParams{
		Provider: provider,
	})
	require.NoError(t, err)
	assert.Equal(t,
		OutboundsConfig{
			Outbounds: map[string]OutboundConfig{
				"bar": {Address: "127.0.0.1:0"},
			},
		},
		res.Config,
	)
}

func TestNewClients(t *testing.T) {
	tests := []struct {
		desc        string
		giveCfg     OutboundConfig
		wantCaller  string
		wantName    string
		wantService string
		wantErr     string
	}{
		{
			desc:        "chooser successfully configured",
			giveCfg:     OutboundConfig{Chooser: "roundrobin"},
			wantCaller:  "foo",
			wantName:    "bar",
			wantService: "bar",
		},
		{
			desc:    "chooser does not exist",
			giveCfg: OutboundConfig{Chooser: "dne"},
			wantErr: `failed to resolve outbound peer list chooser: "dne"`,
		},
		{
			desc:        "address successfully configured",
			giveCfg:     OutboundConfig{Address: "127.0.0.1:0"},
			wantCaller:  "foo",
			wantName:    "bar",
			wantService: "bar",
		},
		{
			desc:        "with configured name",
			giveCfg:     OutboundConfig{Address: "127.0.0.1:0", Service: "baz"},
			wantCaller:  "foo",
			wantName:    "bar",
			wantService: "baz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			res, err := NewClients(ClientParams{
				Lifecycle: fxtest.NewLifecycle(t),
				Config: OutboundsConfig{
					Outbounds: map[string]OutboundConfig{
						"bar": tt.giveCfg,
					},
				},
				Dialer:          &yarpctchannel.Dialer{},
				ChooserProvider: newChooserProvider(t),
			})
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Len(t, res.Clients, 1)

			client := res.Clients[0]
			assert.Equal(t, client.Caller, tt.wantCaller)
			assert.Equal(t, client.Name, tt.wantName)
			assert.Equal(t, client.Service, tt.wantService)
		})
	}
}

func TestNewDialer(t *testing.T) {
	result, err := NewDialer(DialerParams{
		Lifecycle: fxtest.NewLifecycle(t),
		Logger:    zap.NewNop(),
		Tracer:    opentracing.NoopTracer{},
	})
	require.NoError(t, err)

	assert.NotNil(t, result.Dialer)
	assert.NotNil(t, result.TChannelDialer)
}