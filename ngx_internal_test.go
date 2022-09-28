package ngx

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCheckServerUpdatesIsValidOnValidInput(t *testing.T) {
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
		toAdd, toDelete, toUpdate := determineServerUpdates(test.updated, test.nginx)

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

func TestCheckStreamServerUpdatesIsValidOnValidInput(t *testing.T) {
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
		if !cmp.Equal(toAdd, test.expectedToAdd) {
			t.Error(cmp.Diff(toAdd, test.expectedToAdd))
		}
		if !cmp.Equal(toDelete, test.expectedToDelete) {
			t.Error(cmp.Diff(toDelete, test.expectedToDelete))
		}
		if !cmp.Equal(toUpdate, test.expectedToUpdate) {
			t.Errorf(cmp.Diff(toDelete, test.expectedToDelete))
		}
	}
}

func TestServerAddressIsValidOnValidInputWithHostAndPort(t *testing.T) {
	t.Parallel()
	input := "example.com:8080"
	want := "example.com:8080"
	got := addPortToServer(input)
	if want != got {
		t.Errorf("want: %s, got %s", want, got)
	}
}

func TestServerAddressIsValidOnValidInputWithIPV4AndPort(t *testing.T) {
	t.Parallel()
	input := "127.0.0.1:8080"
	want := "127.0.0.1:8080"
	got := addPortToServer(input)
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestServerAddressIsValidOnValidInputWithIPV6AndPort(t *testing.T) {
	t.Parallel()
	input := "[::]:8080"
	want := "[::]:8080"
	got := addPortToServer(input)
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestServerAddressIsValidOnValidInputWithUnixSocket(t *testing.T) {
	t.Parallel()
	input := "unix:/path/to/socket"
	want := "unix:/path/to/socket"
	got := addPortToServer(input)
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestServerAddressIsValidOnValidInputWithAddressAndWithoutPort(t *testing.T) {
	t.Parallel()
	input := "example.com"
	want := "example.com:80"
	got := addPortToServer(input)
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestServerAddressIsValidOnValidInputWithIPV4AddressAndWithoutPort(t *testing.T) {
	t.Parallel()
	input := "127.0.0.1"
	want := "127.0.0.1:80"
	got := addPortToServer(input)
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestServerAddressIsValidOnValidInputWithIPV6AddressAndWithoutPort(t *testing.T) {
	t.Parallel()
	input := "[::]"
	want := "[::]:80"
	got := addPortToServer(input)
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestUpstreamServersConfigIsValidOnValidInput(t *testing.T) {
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

func TestUpstreamStreamServersConfigurationIsValidOnValidInput(t *testing.T) {
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

func TestNGINXServerStatusIsValidOnValidInputRequestParams(t *testing.T) {
	t.Parallel()

	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"version":"1.21.6"}`)
	}

	r := httptest.NewRequest(http.MethodGet, "/8/nginx?fields=version", nil)

	w := httptest.NewRecorder()
	handler(w, r)

	resp := w.Result()

	body, _ := io.ReadAll(resp.Body)

	want := `{"version":"1.21.6"}`
	got := string(body)

	if want != got {
		t.Error(cmp.Diff(want, got))
	}

	if want != got {
		t.Error(cmp.Diff(want, got))
	}

}

func TestRequestGetNGINXURLIsValidOnValidFields(t *testing.T) {
	t.Parallel()

	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"version":"1.21.6"}`)
	}

	r := httptest.NewRequest(http.MethodGet, "/8/nginx?fields=version", nil)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	w := httptest.NewRecorder()
	handler(w, r)

	resp := w.Result()

	t.Log(resp.Header)

	t.Logf("%+v", r.RequestURI)
	t.Log("request path:")

	//body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	// want := `{"version":"1.21.6"}`
	// got := string(body)

	// if want != got {
	// 	t.Error(cmp.Diff(want, got))
	// }

	// if want != got {
	// 	t.Error(cmp.Diff(want, got))
	// }

}
