![Go](https://github.com/qba73/ngx/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/qba73/ngx)](https://goreportcard.com/report/github.com/qba73/ngx)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/qba73/ngx?logo=go)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# ngx

```ngx``` is a Go client library for NGINX Plus API. The project was initially based on the fork of the open source NGINX Plus client API.

The library works against versions 4 to 8 of the NGINX Plus API. The table below shows the version of NGINX Plus where the API was first introduced.

| API version | NGINX Plus version |
|-------------|--------------------|
| 4 | R18 |
| 5 | R19 |
| 6 | R20 |
| 7 | R25 |
| 8 | R27 |

## Using the Go library

Import the library using:

```go
import "github.com/qba73/ngx"
```

## Creating a client

Create a new ```Client``` object by calling ```ngx.NewClient(baseURL)```
```go
client, err := ngx.NewClient("http://localhost:8080/api")
if err != nil {
    // handle error
}
```
Or create a client with customized http Client:
```go
customHTTPClient := &http.Client{}

client, err := ngx.NewClient(
    "http://localhost:8080/api",
    ngx.WithHTTPClient(customHTTPClient),
)
if err != nil {
    // handle error
}
```
Or create a client to work with specific version of NGINX instance:
```go
client, err := ngx.NewClient(
    "http://localhost:8080/api",
    ngx.WithVersion(7),
)
if err != nil {
    // handle error
}

```

## Testing

Install [gotestdox](https://github.com/bitfield/gotestdox)
```bash
go install github.com/bitfield/gotestdox/cmd/gotestdox@latest
```
Run tests
```bash
$ gotestdox
github.com/qba73/ngx:
 ✔ Calculates server updates on valid input (0.00s)
 ✔ Calculates stream server updates on valid input (0.00s)
 ✔ Determines upstream servers configuration equality (0.00s)
 ✔ Determines upstream stream servers configuration equality (0.00s)
 ✔ Builds address on valid input with host and port (0.00s)
 ✔ Builds address on valid input with IPV4 address and without port (0.00s)
 ✔ Builds address on valid input with address and without port (0.00s)
 ✔ Builds address on valid input with unix socket (0.00s)
 ✔ Builds address on valid input with IPV6 and port (0.00s)
 ✔ Builds address on valid input with IPV4 and port (0.00s)
 ✔ Builds address on valid input with IPV6 address and without port (0.00s)
 ✔ Client retrives info about running NGINX instance (0.00s)
 ✔ Client retrives NGINX status on valid parameters (0.00s)
```

Or use [**earthly**](https://docs.earthly.dev)

```bash
$ earthly +checks
```


## Contributing

If you have any suggestions or experience issues with the NGINX Plus Go Client, please create an issue or send a pull request on GitHub.
