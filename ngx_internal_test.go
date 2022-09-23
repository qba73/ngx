package ngx

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDetermineUpdates(t *testing.T) {
	maxConns := 1
	tests := []struct {
		updated          []UpstreamServer
		nginx            []UpstreamServer
		expectedToAdd    []UpstreamServer
		expectedToDelete []UpstreamServer
		expectedToUpdate []UpstreamServer
	}{
		{
			updated: []UpstreamServer{
				{
					Server: "10.0.0.3:80",
				},
				{
					Server: "10.0.0.4:80",
				},
			},
			nginx: []UpstreamServer{
				{
					ID:     1,
					Server: "10.0.0.1:80",
				},
				{
					ID:     2,
					Server: "10.0.0.2:80",
				},
			},
			expectedToAdd: []UpstreamServer{
				{
					Server: "10.0.0.3:80",
				},
				{
					Server: "10.0.0.4:80",
				},
			},
			expectedToDelete: []UpstreamServer{
				{
					ID:     1,
					Server: "10.0.0.1:80",
				},
				{
					ID:     2,
					Server: "10.0.0.2:80",
				},
			},
		},
		{
			updated: []UpstreamServer{
				{
					Server: "10.0.0.2:80",
				},
				{
					Server: "10.0.0.3:80",
				},
				{
					Server: "10.0.0.4:80",
				},
			},
			nginx: []UpstreamServer{
				{
					ID:     1,
					Server: "10.0.0.1:80",
				},
				{
					ID:     2,
					Server: "10.0.0.2:80",
				},
				{
					ID:     3,
					Server: "10.0.0.3:80",
				},
			},
			expectedToAdd: []UpstreamServer{
				{
					Server: "10.0.0.4:80",
				},
			},
			expectedToDelete: []UpstreamServer{
				{
					ID:     1,
					Server: "10.0.0.1:80",
				},
			},
		},
		{
			updated: []UpstreamServer{
				{
					Server: "10.0.0.1:80",
				},
				{
					Server: "10.0.0.2:80",
				},
				{
					Server: "10.0.0.3:80",
				},
			},
			nginx: []UpstreamServer{
				{
					Server: "10.0.0.1:80",
				},
				{
					Server: "10.0.0.2:80",
				},
				{
					Server: "10.0.0.3:80",
				},
			},
		},
		{
			// empty values
		},
		{
			updated: []UpstreamServer{
				{
					Server:   "10.0.0.1:80",
					MaxConns: &maxConns,
				},
			},
			nginx: []UpstreamServer{
				{
					ID:     1,
					Server: "10.0.0.1:80",
				},
				{
					ID:     2,
					Server: "10.0.0.2:80",
				},
			},
			expectedToDelete: []UpstreamServer{
				{
					ID:     2,
					Server: "10.0.0.2:80",
				},
			},
			expectedToUpdate: []UpstreamServer{
				{
					ID:       1,
					Server:   "10.0.0.1:80",
					MaxConns: &maxConns,
				},
			},
		},
	}

	for _, test := range tests {
		toAdd, toDelete, toUpdate := determineUpdates(test.updated, test.nginx)

		if !cmp.Equal(toAdd, test.expectedToAdd) {
			t.Error(cmp.Diff(toAdd, test.expectedToAdd))
		}

		if !cmp.Equal(toDelete, test.expectedToDelete) {
			t.Error(cmp.Diff(toDelete, test.expectedToDelete))
		}

		if !cmp.Equal(toUpdate, test.expectedToUpdate) {
			t.Error(cmp.Diff(toUpdate, test.expectedToUpdate))
		}

	}
}

func TestStreamDetermineUpdates(t *testing.T) {
	maxConns := 1
	tests := []struct {
		updated          []StreamUpstreamServer
		nginx            []StreamUpstreamServer
		expectedToAdd    []StreamUpstreamServer
		expectedToDelete []StreamUpstreamServer
		expectedToUpdate []StreamUpstreamServer
	}{
		{
			updated: []StreamUpstreamServer{
				{
					Server: "10.0.0.3:80",
				},
				{
					Server: "10.0.0.4:80",
				},
			},
			nginx: []StreamUpstreamServer{
				{
					ID:     1,
					Server: "10.0.0.1:80",
				},
				{
					ID:     2,
					Server: "10.0.0.2:80",
				},
			},
			expectedToAdd: []StreamUpstreamServer{
				{
					Server: "10.0.0.3:80",
				},
				{
					Server: "10.0.0.4:80",
				},
			},
			expectedToDelete: []StreamUpstreamServer{
				{
					ID:     1,
					Server: "10.0.0.1:80",
				},
				{
					ID:     2,
					Server: "10.0.0.2:80",
				},
			},
		},
		{
			updated: []StreamUpstreamServer{
				{
					Server: "10.0.0.2:80",
				},
				{
					Server: "10.0.0.3:80",
				},
				{
					Server: "10.0.0.4:80",
				},
			},
			nginx: []StreamUpstreamServer{
				{
					ID:     1,
					Server: "10.0.0.1:80",
				},
				{
					ID:     2,
					Server: "10.0.0.2:80",
				},
				{
					ID:     3,
					Server: "10.0.0.3:80",
				},
			},
			expectedToAdd: []StreamUpstreamServer{
				{
					Server: "10.0.0.4:80",
				},
			},
			expectedToDelete: []StreamUpstreamServer{
				{
					ID:     1,
					Server: "10.0.0.1:80",
				},
			},
		},
		{
			updated: []StreamUpstreamServer{
				{
					Server: "10.0.0.1:80",
				},
				{
					Server: "10.0.0.2:80",
				},
				{
					Server: "10.0.0.3:80",
				},
			},
			nginx: []StreamUpstreamServer{
				{
					ID:     1,
					Server: "10.0.0.1:80",
				},
				{
					ID:     2,
					Server: "10.0.0.2:80",
				},
				{
					ID:     3,
					Server: "10.0.0.3:80",
				},
			},
		},
		{
			// empty values
		},
		{
			updated: []StreamUpstreamServer{
				{
					Server:   "10.0.0.1:80",
					MaxConns: &maxConns,
				},
			},
			nginx: []StreamUpstreamServer{
				{
					ID:     1,
					Server: "10.0.0.1:80",
				},
				{
					ID:     2,
					Server: "10.0.0.2:80",
				},
			},
			expectedToDelete: []StreamUpstreamServer{
				{
					ID:     2,
					Server: "10.0.0.2:80",
				},
			},
			expectedToUpdate: []StreamUpstreamServer{
				{
					ID:       1,
					Server:   "10.0.0.1:80",
					MaxConns: &maxConns,
				},
			},
		},
	}

	for _, test := range tests {
		toAdd, toDelete, toUpdate := determineStreamUpdates(test.updated, test.nginx)
		if !reflect.DeepEqual(toAdd, test.expectedToAdd) || !reflect.DeepEqual(toDelete, test.expectedToDelete) || !reflect.DeepEqual(toUpdate, test.expectedToUpdate) {
			t.Errorf("determiteUpdates(%v, %v) = (%v, %v, %v)", test.updated, test.nginx, toAdd, toDelete, toUpdate)
		}
	}
}

func TestAddPortToServer(t *testing.T) {
	// More info about addresses http://nginx.org/en/docs/http/ngx_http_upstream_module.html#server
	tests := []struct {
		address  string
		expected string
		msg      string
	}{
		{
			address:  "example.com:8080",
			expected: "example.com:8080",
			msg:      "host and port",
		},
		{
			address:  "127.0.0.1:8080",
			expected: "127.0.0.1:8080",
			msg:      "ipv4 and port",
		},
		{
			address:  "[::]:8080",
			expected: "[::]:8080",
			msg:      "ipv6 and port",
		},
		{
			address:  "unix:/path/to/socket",
			expected: "unix:/path/to/socket",
			msg:      "unix socket",
		},
		{
			address:  "example.com",
			expected: "example.com:80",
			msg:      "host without port",
		},
		{
			address:  "127.0.0.1",
			expected: "127.0.0.1:80",
			msg:      "ipv4 without port",
		},
		{
			address:  "[::]",
			expected: "[::]:80",
			msg:      "ipv6 without port",
		},
	}

	for _, test := range tests {
		result := addPortToServer(test.address)
		if result != test.expected {
			t.Errorf("addPortToServer(%v) returned %v but expected %v for %v", test.address, result, test.expected, test.msg)
		}
	}
}

func TestHaveSameParameters(t *testing.T) {
	tests := []struct {
		server    UpstreamServer
		serverNGX UpstreamServer
		expected  bool
	}{
		{
			server:    UpstreamServer{},
			serverNGX: UpstreamServer{},
			expected:  true,
		},
		{
			server:    UpstreamServer{ID: 2},
			serverNGX: UpstreamServer{ID: 3},
			expected:  true,
		},
		{
			server: UpstreamServer{},
			serverNGX: UpstreamServer{
				MaxConns:    &defaultMaxConns,
				MaxFails:    &defaultMaxFails,
				FailTimeout: defaultFailTimeout,
				SlowStart:   defaultSlowStart,
				Backup:      &defaultBackup,
				Weight:      &defaultWeight,
				Down:        &defaultDown,
			},
			expected: true,
		},
		{
			server: UpstreamServer{
				ID:          1,
				Server:      "127.0.0.1",
				MaxConns:    &defaultMaxConns,
				MaxFails:    &defaultMaxFails,
				FailTimeout: defaultFailTimeout,
				SlowStart:   defaultSlowStart,
				Backup:      &defaultBackup,
				Weight:      &defaultWeight,
				Down:        &defaultDown,
			},
			serverNGX: UpstreamServer{
				ID:          1,
				Server:      "127.0.0.1",
				MaxConns:    &defaultMaxConns,
				MaxFails:    &defaultMaxFails,
				FailTimeout: defaultFailTimeout,
				SlowStart:   defaultSlowStart,
				Backup:      &defaultBackup,
				Weight:      &defaultWeight,
				Down:        &defaultDown,
			},
			expected: true,
		},
		{
			server:    UpstreamServer{SlowStart: "10s"},
			serverNGX: UpstreamServer{},
			expected:  false,
		},
		{
			server:    UpstreamServer{},
			serverNGX: UpstreamServer{SlowStart: "10s"},
			expected:  false,
		},
		{
			server:    UpstreamServer{SlowStart: "20s"},
			serverNGX: UpstreamServer{SlowStart: "10s"},
			expected:  false,
		},
	}

	for _, test := range tests {
		result := haveSameParameters(test.server, test.serverNGX)
		if result != test.expected {
			t.Errorf("haveSameParameters(%v, %v) returned %v but expected %v", test.server, test.serverNGX, result, test.expected)
		}
	}
}

func TestHaveSameParametersForStream(t *testing.T) {
	tests := []struct {
		server    StreamUpstreamServer
		serverNGX StreamUpstreamServer
		expected  bool
	}{
		{
			server:    StreamUpstreamServer{},
			serverNGX: StreamUpstreamServer{},
			expected:  true,
		},
		{
			server:    StreamUpstreamServer{ID: 2},
			serverNGX: StreamUpstreamServer{ID: 3},
			expected:  true,
		},
		{
			server: StreamUpstreamServer{},
			serverNGX: StreamUpstreamServer{
				MaxConns:    &defaultMaxConns,
				MaxFails:    &defaultMaxFails,
				FailTimeout: defaultFailTimeout,
				SlowStart:   defaultSlowStart,
				Backup:      &defaultBackup,
				Weight:      &defaultWeight,
				Down:        &defaultDown,
			},
			expected: true,
		},
		{
			server: StreamUpstreamServer{
				ID:          1,
				Server:      "127.0.0.1",
				MaxConns:    &defaultMaxConns,
				MaxFails:    &defaultMaxFails,
				FailTimeout: defaultFailTimeout,
				SlowStart:   defaultSlowStart,
				Backup:      &defaultBackup,
				Weight:      &defaultWeight,
				Down:        &defaultDown,
			},
			serverNGX: StreamUpstreamServer{
				ID:          1,
				Server:      "127.0.0.1",
				MaxConns:    &defaultMaxConns,
				MaxFails:    &defaultMaxFails,
				FailTimeout: defaultFailTimeout,
				SlowStart:   defaultSlowStart,
				Backup:      &defaultBackup,
				Weight:      &defaultWeight,
				Down:        &defaultDown,
			},
			expected: true,
		},
		{
			server:    StreamUpstreamServer{},
			serverNGX: StreamUpstreamServer{SlowStart: "10s"},
			expected:  false,
		},
		{
			server:    StreamUpstreamServer{SlowStart: "20s"},
			serverNGX: StreamUpstreamServer{SlowStart: "10s"},
			expected:  false,
		},
	}

	for _, test := range tests {
		result := haveSameParametersForStream(test.server, test.serverNGX)
		if result != test.expected {
			t.Errorf("haveSameParametersForStream(%v, %v) returned %v but expected %v", test.server, test.serverNGX, result, test.expected)
		}
	}
}

const (
	cacheZone      = "http_cache"
	upstream       = "test"
	streamUpstream = "stream_test"
	streamZoneSync = "zone_test_sync"
	locationZone   = "location_test"
	resolverMetric = "resolver_test"
	reqZone        = "one"
	connZone       = "addr"
	streamConnZone = "addr_stream"
)

func compareUpstreamServers(x []UpstreamServer, y []UpstreamServer) bool {
	var xServers []string
	for _, us := range x {
		xServers = append(xServers, us.Server)
	}
	var yServers []string
	for _, us := range y {
		yServers = append(yServers, us.Server)
	}
	return cmp.Equal(xServers, yServers)
}

func compareStreamUpstreamServers(x []StreamUpstreamServer, y []StreamUpstreamServer) bool {
	var xServers []string
	for _, us := range x {
		xServers = append(xServers, us.Server)
	}
	var yServers []string
	for _, us := range y {
		yServers = append(yServers, us.Server)
	}
	return cmp.Equal(xServers, yServers)
}
