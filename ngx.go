package ngx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/exp/slices"
)

const (
	// APIVersion is the default version of NGINX Plus API supported by the client.
	DefaultAPIVersion = 8

	pathNotFoundCode  = "PathNotFound"
	streamContext     = true
	httpContext       = false
	defaultServerPort = "80"
)

var (
	supportedAPIVersions = []int{4, 5, 6, 7, 8}

	// Default values for servers in Upstreams.
	defaultMaxConns    = 0
	defaultMaxFails    = 1
	defaultFailTimeout = "10s"
	defaultSlowStart   = "0s"
	defaultBackup      = false
	defaultDown        = false
	defaultWeight      = 1
)

// UpstreamServer lets you configure HTTP upstreams.
type UpstreamServer struct {
	ID          int    `json:"id,omitempty"`
	Server      string `json:"server"`
	MaxConns    *int   `json:"max_conns,omitempty"`
	MaxFails    *int   `json:"max_fails,omitempty"`
	FailTimeout string `json:"fail_timeout,omitempty"`
	SlowStart   string `json:"slow_start,omitempty"`
	Route       string `json:"route,omitempty"`
	Backup      *bool  `json:"backup,omitempty"`
	Down        *bool  `json:"down,omitempty"`
	Drain       bool   `json:"drain,omitempty"`
	Weight      *int   `json:"weight,omitempty"`
	Service     string `json:"service,omitempty"`
}

// StreamUpstreamServer lets you configure Stream upstreams.
type StreamUpstreamServer struct {
	ID          int    `json:"id,omitempty"`
	Server      string `json:"server"`
	MaxConns    *int   `json:"max_conns,omitempty"`
	MaxFails    *int   `json:"max_fails,omitempty"`
	FailTimeout string `json:"fail_timeout,omitempty"`
	SlowStart   string `json:"slow_start,omitempty"`
	Backup      *bool  `json:"backup,omitempty"`
	Down        *bool  `json:"down,omitempty"`
	Weight      *int   `json:"weight,omitempty"`
	Service     string `json:"service,omitempty"`
}

type apiErrorResponse struct {
	Error     apiError
	RequestID string `json:"request_id"`
	Href      string
}

func (resp *apiErrorResponse) toString() string {
	return fmt.Sprintf("error.status=%v; error.text=%v; error.code=%v; request_id=%v; href=%v",
		resp.Error.Status, resp.Error.Text, resp.Error.Code, resp.RequestID, resp.Href)
}

type apiError struct {
	Status int
	Text   string
	Code   string
}

type internalError struct {
	apiError
	err string
}

// Error allows internalError to match the Error interface.
func (internalError *internalError) Error() string {
	return internalError.err
}

// Wrap is a way of including current context while preserving previous error information,
// similar to `return fmt.Errorf("error doing foo, err: %v", err)` but for our internalError type.
func (internalError *internalError) Wrap(err string) *internalError {
	internalError.err = fmt.Sprintf("%v. %v", err, internalError.err)
	return internalError
}

// Stats represents NGINX Plus stats fetched from the NGINX Plus API.
//
// https://nginx.org/en/docs/http/ngx_http_api_module.html
type Stats struct {
	NginxInfo              NginxInfo
	Caches                 Caches
	Processes              Processes
	Connections            Connections
	Slabs                  Slabs
	HTTPRequests           HTTPRequests
	SSL                    SSL
	ServerZones            ServerZones
	Upstreams              Upstreams
	StreamServerZones      StreamServerZones
	StreamUpstreams        StreamUpstreams
	StreamZoneSync         StreamZoneSync
	LocationZones          LocationZones
	Resolvers              Resolvers
	HTTPLimitRequests      HTTPLimitRequests
	HTTPLimitConnections   HTTPLimitConnections
	StreamLimitConnections StreamLimitConnections
}

// NginxInfo contains general information about NGINX Plus.
type NginxInfo struct {
	Version         string
	Build           string
	Address         string
	Generation      uint64
	LoadTimestamp   string `json:"load_timestamp"`
	Timestamp       string
	ProcessID       uint64 `json:"pid"`
	ParentProcessID uint64 `json:"ppid"`
}

// Caches is a map of cache stats by cache zone
type Caches = map[string]HTTPCache

// HTTPCache represents a zone's HTTP Cache
type HTTPCache struct {
	Size        uint64
	MaxSize     uint64 `json:"max_size"`
	Cold        bool
	Hit         CacheStats
	Stale       CacheStats
	Updating    CacheStats
	Revalidated CacheStats
	Miss        CacheStats
	Expired     ExtendedCacheStats
	Bypass      ExtendedCacheStats
}

// CacheStats are basic cache stats.
type CacheStats struct {
	Responses uint64
	Bytes     uint64
}

// ExtendedCacheStats are extended cache stats.
type ExtendedCacheStats struct {
	CacheStats
	ResponsesWritten uint64 `json:"responses_written"`
	BytesWritten     uint64 `json:"bytes_written"`
}

// Connections represents connection related stats.
type Connections struct {
	Accepted uint64
	Dropped  uint64
	Active   uint64
	Idle     uint64
}

// Slabs is map of slab stats by zone name.
type Slabs map[string]Slab

// Slab represents slab related stats.
type Slab struct {
	Pages Pages
	Slots Slots
}

// Pages represents the slab memory usage stats.
type Pages struct {
	Used uint64
	Free uint64
}

// Slots is a map of slots by slot size
type Slots map[string]Slot

// Slot represents slot related stats.
type Slot struct {
	Used  uint64
	Free  uint64
	Reqs  uint64
	Fails uint64
}

// HTTPRequests represents HTTP request related stats.
type HTTPRequests struct {
	Total   uint64
	Current uint64
}

// SSL represents SSL related stats.
type SSL struct {
	Handshakes       uint64
	HandshakesFailed uint64 `json:"handshakes_failed"`
	SessionReuses    uint64 `json:"session_reuses"`
}

// ServerZones is map of server zone stats by zone name
type ServerZones map[string]ServerZone

// ServerZone represents server zone related stats.
type ServerZone struct {
	Processing uint64
	Requests   uint64
	Responses  Responses
	Discarded  uint64
	Received   uint64
	Sent       uint64
	SSL        SSL
}

// StreamServerZones is map of stream server zone stats by zone name.
type StreamServerZones map[string]StreamServerZone

// StreamServerZone represents stream server zone related stats.
type StreamServerZone struct {
	Processing  uint64
	Connections uint64
	Sessions    Sessions
	Discarded   uint64
	Received    uint64
	Sent        uint64
	SSL         SSL
}

// StreamZoneSync represents the sync information per each shared memory zone and the sync information per node in a cluster
type StreamZoneSync struct {
	Zones  map[string]SyncZone
	Status StreamZoneSyncStatus
}

// SyncZone represents the synchronization status of a shared memory zone
type SyncZone struct {
	RecordsPending uint64 `json:"records_pending"`
	RecordsTotal   uint64 `json:"records_total"`
}

// StreamZoneSyncStatus represents the status of a shared memory zone
type StreamZoneSyncStatus struct {
	BytesIn     uint64 `json:"bytes_in"`
	MsgsIn      uint64 `json:"msgs_in"`
	MsgsOut     uint64 `json:"msgs_out"`
	BytesOut    uint64 `json:"bytes_out"`
	NodesOnline uint64 `json:"nodes_online"`
}

// Responses represents HTTP response related stats.
type Responses struct {
	Codes        HTTPCodes
	Responses1xx uint64 `json:"1xx"`
	Responses2xx uint64 `json:"2xx"`
	Responses3xx uint64 `json:"3xx"`
	Responses4xx uint64 `json:"4xx"`
	Responses5xx uint64 `json:"5xx"`
	Total        uint64
}

// HTTPCodes represents HTTP response codes
type HTTPCodes struct {
	HTTPContinue              uint64 `json:"100,omitempty"`
	HTTPSwitchingProtocols    uint64 `json:"101,omitempty"`
	HTTPProcessing            uint64 `json:"102,omitempty"`
	HTTPOk                    uint64 `json:"200,omitempty"`
	HTTPCreated               uint64 `json:"201,omitempty"`
	HTTPAccepted              uint64 `json:"202,omitempty"`
	HTTPNoContent             uint64 `json:"204,omitempty"`
	HTTPPartialContent        uint64 `json:"206,omitempty"`
	HTTPSpecialResponse       uint64 `json:"300,omitempty"`
	HTTPMovedPermanently      uint64 `json:"301,omitempty"`
	HTTPMovedTemporarily      uint64 `json:"302,omitempty"`
	HTTPSeeOther              uint64 `json:"303,omitempty"`
	HTTPNotModified           uint64 `json:"304,omitempty"`
	HTTPTemporaryRedirect     uint64 `json:"307,omitempty"`
	HTTPBadRequest            uint64 `json:"400,omitempty"`
	HTTPUnauthorized          uint64 `json:"401,omitempty"`
	HTTPForbidden             uint64 `json:"403,omitempty"`
	HTTPNotFound              uint64 `json:"404,omitempty"`
	HTTPNotAllowed            uint64 `json:"405,omitempty"`
	HTTPRequestTimeOut        uint64 `json:"408,omitempty"`
	HTTPConflict              uint64 `json:"409,omitempty"`
	HTTPLengthRequired        uint64 `json:"411,omitempty"`
	HTTPPreconditionFailed    uint64 `json:"412,omitempty"`
	HTTPRequestEntityTooLarge uint64 `json:"413,omitempty"`
	HTTPRequestURITooLarge    uint64 `json:"414,omitempty"`
	HTTPUnsupportedMediaType  uint64 `json:"415,omitempty"`
	HTTPRangeNotSatisfiable   uint64 `json:"416,omitempty"`
	HTTPTooManyRequests       uint64 `json:"429,omitempty"`
	HTTPClose                 uint64 `json:"444,omitempty"`
	HTTPRequestHeaderTooLarge uint64 `json:"494,omitempty"`
	HTTPSCertError            uint64 `json:"495,omitempty"`
	HTTPSNoCert               uint64 `json:"496,omitempty"`
	HTTPToHTTPS               uint64 `json:"497,omitempty"`
	HTTPClientClosedRequest   uint64 `json:"499,omitempty"`
	HTTPInternalServerError   uint64 `json:"500,omitempty"`
	HTTPNotImplemented        uint64 `json:"501,omitempty"`
	HTTPBadGateway            uint64 `json:"502,omitempty"`
	HTTPServiceUnavailable    uint64 `json:"503,omitempty"`
	HTTPGatewayTimeOut        uint64 `json:"504,omitempty"`
	HTTPInsufficientStorage   uint64 `json:"507,omitempty"`
}

// Sessions represents stream session related stats.
type Sessions struct {
	Sessions2xx uint64 `json:"2xx"`
	Sessions4xx uint64 `json:"4xx"`
	Sessions5xx uint64 `json:"5xx"`
	Total       uint64
}

// Upstreams is a map of upstream stats by upstream name.
type Upstreams map[string]Upstream

// Upstream represents upstream related stats.
type Upstream struct {
	Peers      []Peer
	Keepalives int
	Zombies    int
	Zone       string
	Queue      Queue
}

// StreamUpstreams is a map of stream upstream stats by upstream name.
type StreamUpstreams map[string]StreamUpstream

// StreamUpstream represents stream upstream related stats.
type StreamUpstream struct {
	Peers   []StreamPeer
	Zombies int
	Zone    string
}

// Queue represents queue related stats for an upstream.
type Queue struct {
	Size      int
	MaxSize   int `json:"max_size"`
	Overflows uint64
}

// Peer represents peer (upstream server) related stats.
type Peer struct {
	ID           int
	Server       string
	Service      string
	Name         string
	Backup       bool
	Weight       int
	State        string
	Active       uint64
	SSL          SSL
	MaxConns     int `json:"max_conns"`
	Requests     uint64
	Responses    Responses
	Sent         uint64
	Received     uint64
	Fails        uint64
	Unavail      uint64
	HealthChecks HealthChecks `json:"health_checks"`
	Downtime     uint64
	Downstart    string
	Selected     string
	HeaderTime   uint64 `json:"header_time"`
	ResponseTime uint64 `json:"response_time"`
}

// StreamPeer represents peer (stream upstream server) related stats.
type StreamPeer struct {
	ID            int
	Server        string
	Service       string
	Name          string
	Backup        bool
	Weight        int
	State         string
	Active        uint64
	SSL           SSL
	MaxConns      int `json:"max_conns"`
	Connections   uint64
	ConnectTime   int    `json:"connect_time"`
	FirstByteTime int    `json:"first_byte_time"`
	ResponseTime  uint64 `json:"response_time"`
	Sent          uint64
	Received      uint64
	Fails         uint64
	Unavail       uint64
	HealthChecks  HealthChecks `json:"health_checks"`
	Downtime      uint64
	Downstart     string
	Selected      string
}

// HealthChecks represents health check related stats for a peer.
type HealthChecks struct {
	Checks     uint64
	Fails      uint64
	Unhealthy  uint64
	LastPassed bool `json:"last_passed"`
}

// LocationZones represents location_zones related stats
type LocationZones map[string]LocationZone

// Resolvers represents resolvers related stats
type Resolvers map[string]Resolver

// LocationZone represents location_zones related stats
type LocationZone struct {
	Requests  int64
	Responses Responses
	Discarded int64
	Received  int64
	Sent      int64
}

// Resolver represents resolvers related stats
type Resolver struct {
	Requests  ResolverRequests  `json:"requests"`
	Responses ResolverResponses `json:"responses"`
}

// ResolverRequests represents resolver requests
type ResolverRequests struct {
	Name int64
	Srv  int64
	Addr int64
}

// ResolverResponses represents resolver responses
type ResolverResponses struct {
	Noerror  int64
	Formerr  int64
	Servfail int64
	Nxdomain int64
	Notimp   int64
	Refused  int64
	Timedout int64
	Unknown  int64
}

// Processes represents processes related stats
type Processes struct {
	Respawned int64
}

// HTTPLimitRequest represents HTTP Requests Rate Limiting
type HTTPLimitRequest struct {
	Passed         uint64
	Delayed        uint64
	Rejected       uint64
	DelayedDryRun  uint64 `json:"delayed_dry_run"`
	RejectedDryRun uint64 `json:"rejected_dry_run"`
}

// HTTPLimitRequests represents limit requests related stats
type HTTPLimitRequests map[string]HTTPLimitRequest

// LimitConnection represents Connections Limiting
type LimitConnection struct {
	Passed         uint64
	Rejected       uint64
	RejectedDryRun uint64 `json:"rejected_dry_run"`
}

// HTTPLimitConnections represents limit connections related stats
type HTTPLimitConnections map[string]LimitConnection

// StreamLimitConnections represents limit connections related stats
type StreamLimitConnections map[string]LimitConnection

// Option helps to configure the Client with user specified parameters.
type Option func(*Client) error

// WithHttpClient is a func option that configures NGINX Client
// to use a custom HTTP Client.
func WithHttpClient(h *http.Client) Option {
	return func(c *Client) error {
		if h == nil {
			return errors.New("nil http client")
		}
		c.HttpClient = h
		return nil
	}
}

// WithBaseURL is a func option that configures Client's base URL.
func WithBaseURL(s string) Option {
	return func(c *Client) error {
		if s == "" {
			return errors.New("empty base URL")
		}
		c.BaseURL = s
		return nil
	}
}

// WithVersion is a func option that configures version of the NGINX API
// the Client talks to. It is user's responsibility to provide valid
// version of the NGINX Plus that the Client talks to.
// Valid versions are 4,5,6,7,8. The Client's default version is 8.
func WithVersion(v int) Option {
	return func(c *Client) error {
		if !slices.Contains(supportedAPIVersions, v) {
			return errors.New("unsupported NGINX API version")
		}
		c.Version = v
		return nil
	}
}

// NginxClient lets you access NGINX Plus API.
type Client struct {
	Version    int
	BaseURL    string
	HttpClient *http.Client
}

// NewClient takes NGIXN base URL and constructs a new default client.
// The client can be customized by passing functional options that
// allow to configure client version and HTTP Client used by NGINX client.
func NewClient(baseURL string, opts ...Option) (*Client, error) {
	c := Client{
		Version: DefaultAPIVersion,
		BaseURL: baseURL,
		HttpClient: &http.Client{
			Timeout: time.Second * 5,
		},
	}
	for _, opt := range opts {
		if err := opt(&c); err != nil {
			return nil, err
		}
	}
	return &c, nil
}

// CheckIfUpstreamExists checks if the upstream exists in NGINX.
// If the upstream doesn't exist, it returns the error.
func (c *Client) CheckIfUpstreamExists(upstream string) error {
	_, err := c.GetHTTPServers(upstream)
	return err
}

// GetHTTPServers returns the servers of the upstream from NGINX.
func (c *Client) GetHTTPServers(upstream string) ([]UpstreamServer, error) {
	path := fmt.Sprintf("http/upstreams/%v/servers", upstream)
	var servers []UpstreamServer
	err := c.get(path, &servers)
	if err != nil {
		return nil, fmt.Errorf("retriving HTTP servers of upstream %v: %w", upstream, err)
	}
	return servers, nil
}

// AddHTTPServer adds the server to the upstream.
func (c *Client) AddHTTPServer(upstream string, server UpstreamServer) error {
	id, err := c.getIDOfHTTPServer(upstream, server.Server)
	if err != nil {
		return fmt.Errorf("adding %v server to %v upstream: %w", server.Server, upstream, err)
	}
	if id != -1 {
		return fmt.Errorf("adding %v server to %v upstream: server already exists", server.Server, upstream)
	}
	path := fmt.Sprintf("http/upstreams/%v/servers/", upstream)
	err = c.post(path, server)
	if err != nil {
		return fmt.Errorf("adding %v server to %v upstream: %w", server.Server, upstream, err)
	}
	return nil
}

// DeleteHTTPServer the server from the upstream.
func (c *Client) DeleteHTTPServer(upstream string, server string) error {
	id, err := c.getIDOfHTTPServer(upstream, server)
	if err != nil {
		return fmt.Errorf("removing %v server from  %v upstream: %w", server, upstream, err)
	}
	if id == -1 {
		return fmt.Errorf("removing %v server from %v upstream: server doesn't exist", server, upstream)
	}
	path := fmt.Sprintf("http/upstreams/%v/servers/%v", upstream, id)
	err = c.delete(path, http.StatusOK)
	if err != nil {
		return fmt.Errorf("removing %v server from %v upstream: %w", server, upstream, err)
	}
	return nil
}

// UpdateHTTPServers updates the servers of the upstream.
// Servers that are in the slice, but don't exist in NGINX will be added to NGINX.
// Servers that aren't in the slice, but exist in NGINX, will be removed from NGINX.
// Servers that are in the slice and exist in NGINX, but have different parameters, will be updated.
func (c *Client) UpdateHTTPServers(upstream string, servers []UpstreamServer) ([]UpstreamServer, []UpstreamServer, []UpstreamServer, error) {
	serversInNginx, err := c.GetHTTPServers(upstream)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("updating servers of %v upstream: %w", upstream, err)
	}

	// We assume port 80 if no port is set for servers.
	var formattedServers []UpstreamServer
	for _, server := range servers {
		server.Server = addPortToServer(server.Server)
		formattedServers = append(formattedServers, server)
	}

	toAdd, toDelete, toUpdate := determineUpdates(formattedServers, serversInNginx)

	for _, server := range toAdd {
		err := c.AddHTTPServer(upstream, server)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("updating servers of %v upstream: %w", upstream, err)
		}
	}

	for _, server := range toDelete {
		err := c.DeleteHTTPServer(upstream, server.Server)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("updating servers of %v upstream: %w", upstream, err)
		}
	}

	for _, server := range toUpdate {
		err := c.UpdateHTTPServer(upstream, server)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("updating servers of %v upstream: %w", upstream, err)
		}
	}

	return toAdd, toDelete, toUpdate, nil
}

// haveSameParameters checks if a given server has the same parameters as a server already present in NGINX. Order matters
func haveSameParameters(newServer UpstreamServer, serverNGX UpstreamServer) bool {
	newServer.ID = serverNGX.ID

	if serverNGX.MaxConns != nil && newServer.MaxConns == nil {
		newServer.MaxConns = &defaultMaxConns
	}

	if serverNGX.MaxFails != nil && newServer.MaxFails == nil {
		newServer.MaxFails = &defaultMaxFails
	}

	if serverNGX.FailTimeout != "" && newServer.FailTimeout == "" {
		newServer.FailTimeout = defaultFailTimeout
	}

	if serverNGX.SlowStart != "" && newServer.SlowStart == "" {
		newServer.SlowStart = defaultSlowStart
	}

	if serverNGX.Backup != nil && newServer.Backup == nil {
		newServer.Backup = &defaultBackup
	}

	if serverNGX.Down != nil && newServer.Down == nil {
		newServer.Down = &defaultDown
	}

	if serverNGX.Weight != nil && newServer.Weight == nil {
		newServer.Weight = &defaultWeight
	}

	return cmp.Equal(newServer, serverNGX)
}

func (c *Client) getIDOfHTTPServer(upstream string, name string) (int, error) {
	servers, err := c.GetHTTPServers(upstream)
	if err != nil {
		return -1, fmt.Errorf("getting id of server %v of upstream %v: %w", name, upstream, err)
	}
	for _, s := range servers {
		if s.Server == name {
			return s.ID, nil
		}
	}
	return -1, nil
}

func (c *Client) get(path string, data interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url := fmt.Sprintf("%v/%v/%v", c.BaseURL, c.Version, path)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending request, path: %s, %w", url, err)
	}
	if resp.StatusCode != http.StatusOK {
		return createResponseMismatchError(resp.Body).Wrap(fmt.Sprintf(
			"expected %v response, got %v",
			http.StatusOK, resp.StatusCode))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		return fmt.Errorf("unmarshaling response: %w", err)
	}
	return nil
}

func (c *Client) post(path string, payload interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url := fmt.Sprintf("%v/%v/%v", c.BaseURL, c.Version, path)

	jsonInput, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshalling input: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonInput))
	if err != nil {
		return fmt.Errorf("creating POST request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending POST request %v: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return createResponseMismatchError(resp.Body).Wrap(fmt.Sprintf(
			"expected %v response, got %v",
			http.StatusCreated, resp.StatusCode))
	}

	return nil
}

func (c *Client) delete(path string, expectedStatusCode int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	path = fmt.Sprintf("%v/%v/%v/", c.BaseURL, c.Version, path)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("creating DELETE request: %w", err)
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending DELETE request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatusCode {
		return createResponseMismatchError(resp.Body).Wrap(fmt.Sprintf(
			"failed to complete delete request: expected %v response, got %v",
			expectedStatusCode, resp.StatusCode))
	}
	return nil
}

func (c *Client) patch(path string, input interface{}, expectedStatusCode int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	path = fmt.Sprintf("%v/%v/%v/", c.BaseURL, c.Version, path)

	jsonInput, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("marshalling input: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, path, bytes.NewBuffer(jsonInput))
	if err != nil {
		return fmt.Errorf("creating PATCH request: %w", err)
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending PATCH request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatusCode {
		return createResponseMismatchError(resp.Body).Wrap(fmt.Sprintf(
			"failed to complete patch request: expected %v response, got %v",
			expectedStatusCode, resp.StatusCode))
	}
	return nil
}

// CheckIfStreamUpstreamExists checks if the stream upstream exists in NGINX. If the upstream doesn't exist, it returns the error.
func (c *Client) CheckIfStreamUpstreamExists(upstream string) error {
	_, err := c.GetStreamServers(upstream)
	return err
}

// GetStreamServers returns the stream servers of the upstream from NGINX.
func (client *Client) GetStreamServers(upstream string) ([]StreamUpstreamServer, error) {
	path := fmt.Sprintf("stream/upstreams/%v/servers", upstream)

	var servers []StreamUpstreamServer
	err := client.get(path, &servers)
	if err != nil {
		return nil, fmt.Errorf("getting stream servers of upstream server %v: %w", upstream, err)
	}
	return servers, nil
}

// AddStreamServer adds the stream server to the upstream.
func (c *Client) AddStreamServer(upstream string, server StreamUpstreamServer) error {
	id, err := c.getIDOfStreamServer(upstream, server.Server)
	if err != nil {
		return fmt.Errorf("adding %v stream server to %v upstream: %w", server.Server, upstream, err)
	}
	if id != -1 {
		return fmt.Errorf("adding %v stream server to %v upstream: server already exists", server.Server, upstream)
	}
	path := fmt.Sprintf("stream/upstreams/%v/servers/", upstream)
	err = c.post(path, &server)
	if err != nil {
		return fmt.Errorf("adding %v stream server to %v upstream: %w", server.Server, upstream, err)
	}
	return nil
}

// DeleteStreamServer the server from the upstream.
func (c *Client) DeleteStreamServer(upstream string, server string) error {
	id, err := c.getIDOfStreamServer(upstream, server)
	if err != nil {
		return fmt.Errorf("removing %v stream server from  %v upstream: %w", server, upstream, err)
	}
	if id == -1 {
		return fmt.Errorf("removing %v stream server from %v upstream: server doesn't exist", server, upstream)
	}

	path := fmt.Sprintf("stream/upstreams/%v/servers/%v", upstream, id)
	err = c.delete(path, http.StatusOK)
	if err != nil {
		return fmt.Errorf("removing %v stream server from %v upstream: %w", server, upstream, err)
	}
	return nil
}

// UpdateStreamServers updates the servers of the upstream.
// Servers that are in the slice, but don't exist in NGINX will be added to NGINX.
// Servers that aren't in the slice, but exist in NGINX, will be removed from NGINX.
// Servers that are in the slice and exist in NGINX, but have different parameters, will be updated.
func (c *Client) UpdateStreamServers(upstream string, servers []StreamUpstreamServer) ([]StreamUpstreamServer, []StreamUpstreamServer, []StreamUpstreamServer, error) {
	serversInNginx, err := c.GetStreamServers(upstream)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("updating stream servers of %v upstream: %w", upstream, err)
	}

	var formattedServers []StreamUpstreamServer
	for _, server := range servers {
		server.Server = addPortToServer(server.Server)
		formattedServers = append(formattedServers, server)
	}

	toAdd, toDelete, toUpdate := determineStreamUpdates(formattedServers, serversInNginx)

	for _, server := range toAdd {
		err := c.AddStreamServer(upstream, server)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("updating stream servers of %v upstream: %w", upstream, err)
		}
	}

	for _, server := range toDelete {
		err := c.DeleteStreamServer(upstream, server.Server)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("updating stream servers of %v upstream: %w", upstream, err)
		}
	}

	for _, server := range toUpdate {
		err := c.UpdateStreamServer(upstream, server)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("updating stream servers of %v upstream: %w", upstream, err)
		}
	}

	return toAdd, toDelete, toUpdate, nil
}

func (c *Client) getIDOfStreamServer(upstream string, name string) (int, error) {
	servers, err := c.GetStreamServers(upstream)
	if err != nil {
		return -1, fmt.Errorf("getting id of stream server %v of upstream %v: %w", name, upstream, err)
	}
	for _, s := range servers {
		if s.Server == name {
			return s.ID, nil
		}
	}
	return -1, nil
}

// GetStats gets process, slab, connection, request, ssl, zone, stream zone, upstream and stream upstream related stats from the NGINX Plus API.
func (c *Client) GetStats() (_ Stats, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("getting stats: %w", err)
		}
	}()

	info, err := c.GetNginxInfo()
	if err != nil {
		return Stats{}, err
	}

	caches, err := c.GetCaches()
	if err != nil {
		return Stats{}, err
	}

	processes, err := c.GetProcesses()
	if err != nil {
		return Stats{}, err
	}

	slabs, err := c.GetSlabs()
	if err != nil {
		return Stats{}, err
	}

	cons, err := c.GetConnections()
	if err != nil {
		return Stats{}, err
	}

	requests, err := c.GetHTTPRequests()
	if err != nil {
		return Stats{}, err
	}

	ssl, err := c.GetSSL()
	if err != nil {
		return Stats{}, err
	}

	zones, err := c.GetServerZones()
	if err != nil {
		return Stats{}, err
	}

	upstreams, err := c.GetUpstreams()
	if err != nil {
		return Stats{}, err
	}

	streamZones, err := c.GetStreamServerZones()
	if err != nil {
		return Stats{}, err
	}

	streamUpstreams, err := c.GetStreamUpstreams()
	if err != nil {
		return Stats{}, err
	}

	streamZoneSync, err := c.GetStreamZoneSync()
	if err != nil {
		return Stats{}, err
	}

	locationZones, err := c.GetLocationZones()
	if err != nil {
		return Stats{}, err
	}

	resolvers, err := c.GetResolvers()
	if err != nil {
		return Stats{}, err
	}

	limitReqs, err := c.GetHTTPLimitReqs()
	if err != nil {
		return Stats{}, err
	}

	limitConnsHTTP, err := c.GetHTTPConnectionsLimit()
	if err != nil {
		return Stats{}, err
	}

	limitConnsStream, err := c.GetStreamConnectionsLimit()
	if err != nil {
		return Stats{}, err
	}

	return Stats{
		NginxInfo:              info,
		Caches:                 caches,
		Processes:              processes,
		Slabs:                  slabs,
		Connections:            cons,
		HTTPRequests:           requests,
		SSL:                    ssl,
		ServerZones:            zones,
		StreamServerZones:      streamZones,
		Upstreams:              upstreams,
		StreamUpstreams:        streamUpstreams,
		StreamZoneSync:         streamZoneSync,
		LocationZones:          locationZones,
		Resolvers:              resolvers,
		HTTPLimitRequests:      limitReqs,
		HTTPLimitConnections:   limitConnsHTTP,
		StreamLimitConnections: limitConnsStream,
	}, nil
}

// GetNginxInfo returns Nginx stats.
func (c *Client) GetNginxInfo() (NginxInfo, error) {
	var info NginxInfo
	err := c.get("nginx", &info)
	if err != nil {
		return NginxInfo{}, fmt.Errorf("getting NGINX info: %w", err)
	}
	return info, nil
}

// GetCaches returns Cache stats
func (c *Client) GetCaches() (Caches, error) {
	var caches Caches
	err := c.get("http/caches", &caches)
	if err != nil {
		return nil, fmt.Errorf("getting caches: %w", err)
	}
	return caches, nil
}

// GetSlabs returns Slabs stats.
func (c *Client) GetSlabs() (Slabs, error) {
	var slabs Slabs
	err := c.get("slabs", &slabs)
	if err != nil {
		return nil, fmt.Errorf("getting slabs: %w", err)
	}
	return slabs, nil
}

// GetConnections returns Connections stats.
func (c *Client) GetConnections() (Connections, error) {
	var cons Connections
	err := c.get("connections", &cons)
	if err != nil {
		return Connections{}, fmt.Errorf("failed to get connections: %w", err)
	}
	return cons, nil
}

// GetHTTPRequests returns http/requests stats.
func (c *Client) GetHTTPRequests() (HTTPRequests, error) {
	var requests HTTPRequests
	if err := c.get("http/requests", &requests); err != nil {
		return HTTPRequests{}, fmt.Errorf("getting http requests: %w", err)
	}
	return requests, nil
}

// GetSSL returns SSL stats.
func (c *Client) GetSSL() (SSL, error) {
	var ssl SSL
	if err := c.get("ssl", &ssl); err != nil {
		return SSL{}, fmt.Errorf("getting ssl: %w", err)
	}
	return ssl, nil
}

// GetServerZones returns http/server_zones stats.
func (c *Client) GetServerZones() (ServerZones, error) {
	var zones ServerZones
	if err := c.get("http/server_zones", &zones); err != nil {
		return nil, fmt.Errorf("getting server zones: %w", err)
	}
	return zones, nil
}

// GetStreamServerZones returns stream/server_zones stats.
func (c *Client) GetStreamServerZones() (StreamServerZones, error) {
	var zones StreamServerZones
	err := c.get("stream/server_zones", &zones)
	if err != nil {
		var ie *internalError
		if errors.As(err, &ie) {
			if ie.Code == pathNotFoundCode {
				return zones, nil
			}
		}
		return nil, fmt.Errorf("getting stream server zones: %w", err)
	}
	return zones, err
}

// GetUpstreams returns http/upstreams stats.
func (c *Client) GetUpstreams() (Upstreams, error) {
	var upstreams Upstreams
	err := c.get("http/upstreams", &upstreams)
	if err != nil {
		return nil, fmt.Errorf("getting upstreams: %w", err)
	}
	return upstreams, nil
}

// GetStreamUpstreams returns stream/upstreams stats.
func (c *Client) GetStreamUpstreams() (StreamUpstreams, error) {
	var upstreams StreamUpstreams
	err := c.get("stream/upstreams", &upstreams)
	if err != nil {
		var ie *internalError
		if errors.As(err, &ie) {
			if ie.Code == pathNotFoundCode {
				return upstreams, nil
			}
		}
		return nil, fmt.Errorf("getting stream upstreams: %w", err)
	}
	return upstreams, nil
}

// GetStreamZoneSync returns stream/zone_sync stats.
func (c *Client) GetStreamZoneSync() (StreamZoneSync, error) {
	var streamZoneSync StreamZoneSync
	err := c.get("stream/zone_sync", &streamZoneSync)
	if err != nil {
		var ie *internalError
		if errors.As(err, &ie) {
			if ie.Code == pathNotFoundCode {
				return StreamZoneSync{}, nil
			}
		}
		return StreamZoneSync{}, fmt.Errorf("getting stream zone sync: %w", err)
	}
	return streamZoneSync, nil
}

// GetLocationZones returns http/location_zones stats.
func (c *Client) GetLocationZones() (LocationZones, error) {
	var locationZones LocationZones
	if c.Version < 5 {
		return LocationZones{}, nil
	}
	if err := c.get("http/location_zones", &locationZones); err != nil {
		return nil, fmt.Errorf("gettign location zones: %w", err)
	}
	return locationZones, nil
}

// GetResolvers returns Resolvers stats.
func (c *Client) GetResolvers() (Resolvers, error) {
	var resolvers Resolvers
	if c.Version < 5 {
		return Resolvers{}, nil
	}
	if err := c.get("resolvers", &resolvers); err != nil {
		return nil, fmt.Errorf("getting resolvers: %w", err)
	}
	return resolvers, nil
}

// GetProcesses returns Processes stats.
func (c *Client) GetProcesses() (Processes, error) {
	var processes Processes
	err := c.get("processes", &processes)
	if err != nil {
		return Processes{}, fmt.Errorf("getting processes: %w", err)
	}
	return processes, nil
}

// KeyValPairs are the key-value pairs stored in a zone.
type KeyValPairs map[string]string

// KeyValPairsByZone are the KeyValPairs for all zones, by zone name.
type KeyValPairsByZone map[string]KeyValPairs

// GetKeyValPairs fetches key/value pairs for a given HTTP zone.
func (c *Client) GetKeyValPairs(zone string) (KeyValPairs, error) {
	return c.getKeyValPairs(zone, httpContext)
}

// GetStreamKeyValPairs fetches key/value pairs for a given Stream zone.
func (c *Client) GetStreamKeyValPairs(zone string) (KeyValPairs, error) {
	return c.getKeyValPairs(zone, streamContext)
}

func (c *Client) getKeyValPairs(zone string, stream bool) (KeyValPairs, error) {
	base := "http"
	if stream {
		base = "stream"
	}
	if zone == "" {
		return nil, fmt.Errorf("zone required")
	}
	path := fmt.Sprintf("%v/keyvals/%v", base, zone)
	var keyValPairs KeyValPairs
	if err := c.get(path, &keyValPairs); err != nil {
		return nil, fmt.Errorf("getting keyvals for %v/%v zone: %w", base, zone, err)
	}
	return keyValPairs, nil
}

// GetAllKeyValPairs fetches all key/value pairs for all HTTP zones.
func (c *Client) GetAllKeyValPairs() (KeyValPairsByZone, error) {
	return c.getAllKeyValPairs(httpContext)
}

// GetAllStreamKeyValPairs fetches all key/value pairs for all Stream zones.
func (c *Client) GetAllStreamKeyValPairs() (KeyValPairsByZone, error) {
	return c.getAllKeyValPairs(streamContext)
}

func (c *Client) getAllKeyValPairs(stream bool) (KeyValPairsByZone, error) {
	base := "http"
	if stream {
		base = "stream"
	}
	path := fmt.Sprintf("%v/keyvals", base)

	var keyValPairsByZone KeyValPairsByZone
	if err := c.get(path, &keyValPairsByZone); err != nil {
		return nil, fmt.Errorf("getting keyvals for all %v zones: %w", base, err)
	}
	return keyValPairsByZone, nil
}

// AddKeyValPair adds a new key/value pair to a given HTTP zone.
func (c *Client) AddKeyValPair(zone string, key string, val string) error {
	return c.addKeyValPair(zone, key, val, httpContext)
}

// AddStreamKeyValPair adds a new key/value pair to a given Stream zone.
func (c *Client) AddStreamKeyValPair(zone string, key string, val string) error {
	return c.addKeyValPair(zone, key, val, streamContext)
}

func (c *Client) addKeyValPair(zone string, key string, val string, stream bool) error {
	base := "http"
	if stream {
		base = "stream"
	}
	if zone == "" {
		return fmt.Errorf("zone required")
	}
	path := fmt.Sprintf("%v/keyvals/%v", base, zone)
	input := KeyValPairs{key: val}
	if err := c.post(path, &input); err != nil {
		return fmt.Errorf("adding key value pair for %v/%v zone: %w", base, zone, err)
	}
	return nil
}

// ModifyKeyValPair modifies the value of an existing key in a given HTTP zone.
func (c *Client) ModifyKeyValPair(zone string, key string, val string) error {
	return c.modifyKeyValPair(zone, key, val, httpContext)
}

// Modify10KeyValPair modifies the value of an existing key in a given Stream zone.
func (c *Client) ModifyStreamKeyValPair(zone string, key string, val string) error {
	return c.modifyKeyValPair(zone, key, val, streamContext)
}

func (c *Client) modifyKeyValPair(zone string, key string, val string, stream bool) error {
	base := "http"
	if stream {
		base = "stream"
	}
	if zone == "" {
		return fmt.Errorf("zone required")
	}

	path := fmt.Sprintf("%v/keyvals/%v", base, zone)
	input := KeyValPairs{key: val}
	if err := c.patch(path, &input, http.StatusNoContent); err != nil {
		return fmt.Errorf("updating key value pair for %v/%v zone: %w", base, zone, err)
	}
	return nil
}

// DeleteKeyValuePair deletes the key/value pair for a key in a given HTTP zone.
func (c *Client) DeleteKeyValuePair(zone string, key string) error {
	return c.deleteKeyValuePair(zone, key, httpContext)
}

// DeleteStreamKeyValuePair deletes the key/value pair for a key in a given Stream zone.
func (c *Client) DeleteStreamKeyValuePair(zone string, key string) error {
	return c.deleteKeyValuePair(zone, key, streamContext)
}

// To delete a key/value pair you set the value to null via the API,
// then NGINX+ will delete the key.
func (c *Client) deleteKeyValuePair(zone string, key string, stream bool) error {
	base := "http"
	if stream {
		base = "stream"
	}
	if zone == "" {
		return fmt.Errorf("zone required")
	}

	// map[string]string can't have a nil value so we use a different type here.
	keyval := make(map[string]interface{})
	keyval[key] = nil

	path := fmt.Sprintf("%v/keyvals/%v", base, zone)
	if err := c.patch(path, &keyval, http.StatusNoContent); err != nil {
		return fmt.Errorf("removing key values pair for %v/%v zone: %w", base, zone, err)
	}
	return nil
}

// DeleteKeyValPairs deletes all the key-value pairs in a given HTTP zone.
func (c *Client) DeleteKeyValPairs(zone string) error {
	return c.deleteKeyValPairs(zone, httpContext)
}

// DeleteStreamKeyValPairs deletes all the key-value pairs in a given Stream zone.
func (c *Client) DeleteStreamKeyValPairs(zone string) error {
	return c.deleteKeyValPairs(zone, streamContext)
}

func (c *Client) deleteKeyValPairs(zone string, stream bool) error {
	base := "http"
	if stream {
		base = "stream"
	}
	if zone == "" {
		return fmt.Errorf("zone required")
	}
	path := fmt.Sprintf("%v/keyvals/%v", base, zone)
	if err := c.delete(path, http.StatusNoContent); err != nil {
		return fmt.Errorf("removing all key value pairs for %v/%v zone: %w", base, zone, err)
	}
	return nil
}

// UpdateHTTPServer updates the server of the upstream.
func (c *Client) UpdateHTTPServer(upstream string, server UpstreamServer) error {
	path := fmt.Sprintf("http/upstreams/%v/servers/%v", upstream, server.ID)
	server.ID = 0
	if err := c.patch(path, &server, http.StatusOK); err != nil {
		return fmt.Errorf("updating %v server to %v upstream: %w", server.Server, upstream, err)
	}
	return nil
}

// UpdateStreamServer updates the stream server of the upstream.
func (c *Client) UpdateStreamServer(upstream string, server StreamUpstreamServer) error {
	path := fmt.Sprintf("stream/upstreams/%v/servers/%v", upstream, server.ID)
	server.ID = 0
	if err := c.patch(path, &server, http.StatusOK); err != nil {
		return fmt.Errorf("updating %v stream server to %v upstream: %w", server.Server, upstream, err)
	}
	return nil
}

// GetHTTPLimitReqs returns http/limit_reqs stats.
func (c *Client) GetHTTPLimitReqs() (HTTPLimitRequests, error) {
	var limitReqs HTTPLimitRequests
	if c.Version < 6 {
		return HTTPLimitRequests{}, nil
	}
	if err := c.get("http/limit_reqs", &limitReqs); err != nil {
		return nil, fmt.Errorf("getting http limit requests: %w", err)
	}
	return limitReqs, nil
}

// GetHTTPConnectionsLimit returns http/limit_conns stats.
func (c *Client) GetHTTPConnectionsLimit() (HTTPLimitConnections, error) {
	var limitConns HTTPLimitConnections
	if c.Version < 6 {
		return HTTPLimitConnections{}, nil
	}
	if err := c.get("http/limit_conns", &limitConns); err != nil {
		return nil, fmt.Errorf("getting http connections limit: %w", err)
	}
	return limitConns, nil
}

// GetStreamConnectionsLimit returns stream/limit_conns stats.
func (c *Client) GetStreamConnectionsLimit() (StreamLimitConnections, error) {
	var limitConns StreamLimitConnections
	if c.Version < 6 {
		return StreamLimitConnections{}, nil
	}
	err := c.get("stream/limit_conns", &limitConns)
	if err != nil {
		var ie *internalError
		if errors.As(err, &ie) {
			if ie.Code == pathNotFoundCode {
				return limitConns, nil
			}
		}
		return nil, fmt.Errorf("failed to get stream connections limit: %w", err)
	}
	return limitConns, nil
}

func addPortToServer(server string) string {
	if len(strings.Split(server, ":")) == 2 {
		return server
	}
	if len(strings.Split(server, "]:")) == 2 {
		return server
	}
	if strings.HasPrefix(server, "unix:") {
		return server
	}
	return fmt.Sprintf("%v:%v", server, defaultServerPort)
}

func createResponseMismatchError(respBody io.ReadCloser) *internalError {
	apiErrResp, err := readAPIErrorResponse(respBody)
	if err != nil {
		return &internalError{
			err: fmt.Sprintf("reading response body: %v", err),
		}
	}

	return &internalError{
		err:      apiErrResp.toString(),
		apiError: apiErrResp.Error,
	}
}

func readAPIErrorResponse(respBody io.ReadCloser) (apiErrorResponse, error) {
	body, err := io.ReadAll(respBody)
	if err != nil {
		return apiErrorResponse{}, fmt.Errorf("reading the response body: %w", err)
	}

	var apiErr apiErrorResponse
	if err = json.Unmarshal(body, &apiErr); err != nil {
		return apiErrorResponse{}, fmt.Errorf("unmarshalling apiErrorResponse: got %q response: %w", string(body), err)
	}
	return apiErr, nil
}

func getAPIVersions(httpClient *http.Client, endpoint string) ([]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating GET request: %w", err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request to %v: %w", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got %v requesting API version", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var vers []int
	err = json.Unmarshal(body, &vers)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response: %w", err)
	}
	return vers, nil
}

func determineUpdates(updatedServers []UpstreamServer, nginxServers []UpstreamServer) ([]UpstreamServer, []UpstreamServer, []UpstreamServer) {
	var toAdd, toRemove, toUpdate []UpstreamServer

	for _, server := range updatedServers {
		updateFound := false
		for _, serverNGX := range nginxServers {
			if server.Server == serverNGX.Server && !haveSameParameters(server, serverNGX) {
				server.ID = serverNGX.ID
				updateFound = true
				break
			}
		}
		if updateFound {
			toUpdate = append(toUpdate, server)
		}
	}

	for _, server := range updatedServers {
		found := false
		for _, serverNGX := range nginxServers {
			if server.Server == serverNGX.Server {
				found = true
				break
			}
		}
		if !found {
			toAdd = append(toAdd, server)
		}
	}

	for _, serverNGX := range nginxServers {
		found := false
		for _, server := range updatedServers {
			if serverNGX.Server == server.Server {
				found = true
				break
			}
		}
		if !found {
			toRemove = append(toRemove, serverNGX)
		}
	}

	return toAdd, toRemove, toUpdate
}

func determineStreamUpdates(updatedServers []StreamUpstreamServer, nginxServers []StreamUpstreamServer) ([]StreamUpstreamServer, []StreamUpstreamServer, []StreamUpstreamServer) {
	var toAdd, toRemove, toUpdate []StreamUpstreamServer

	for _, server := range updatedServers {
		updateFound := false
		for _, serverNGX := range nginxServers {
			if server.Server == serverNGX.Server && !haveSameParametersForStream(server, serverNGX) {
				server.ID = serverNGX.ID
				updateFound = true
				break
			}
		}
		if updateFound {
			toUpdate = append(toUpdate, server)
		}
	}

	for _, server := range updatedServers {
		found := false
		for _, serverNGX := range nginxServers {
			if server.Server == serverNGX.Server {
				found = true
				break
			}
		}
		if !found {
			toAdd = append(toAdd, server)
		}
	}

	for _, serverNGX := range nginxServers {
		found := false
		for _, server := range updatedServers {
			if serverNGX.Server == server.Server {
				found = true
				break
			}
		}
		if !found {
			toRemove = append(toRemove, serverNGX)
		}
	}

	return toAdd, toRemove, toUpdate
}

// haveSameParametersForStream checks if a given server has the same parameters as a server already present in NGINX. Order matters
func haveSameParametersForStream(newServer StreamUpstreamServer, serverNGX StreamUpstreamServer) bool {
	newServer.ID = serverNGX.ID
	if serverNGX.MaxConns != nil && newServer.MaxConns == nil {
		newServer.MaxConns = &defaultMaxConns
	}

	if serverNGX.MaxFails != nil && newServer.MaxFails == nil {
		newServer.MaxFails = &defaultMaxFails
	}

	if serverNGX.FailTimeout != "" && newServer.FailTimeout == "" {
		newServer.FailTimeout = defaultFailTimeout
	}

	if serverNGX.SlowStart != "" && newServer.SlowStart == "" {
		newServer.SlowStart = defaultSlowStart
	}

	if serverNGX.Backup != nil && newServer.Backup == nil {
		newServer.Backup = &defaultBackup
	}

	if serverNGX.Down != nil && newServer.Down == nil {
		newServer.Down = &defaultDown
	}

	if serverNGX.Weight != nil && newServer.Weight == nil {
		newServer.Weight = &defaultWeight
	}

	return reflect.DeepEqual(newServer, serverNGX)
}
