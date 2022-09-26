![Go](https://github.com/qba73/ngx/workflows/Go/badge.svg)
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

```
Or 



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

Or use ```earthly``` 

```bash
$ earthly +checks
```
Example output"

<details>
  <summary>Click to see output</summary>

➜  ngx git:(main) ✗ earthly +checks
 1. Init 🚀
————————————————————————————————————————————————————————————————————————————————

           buildkitd | Found buildkit daemon as docker container (earthly-buildkitd)

 2. Build 🔧
————————————————————————————————————————————————————————————————————————————————

golang:1.19-bullseye | --> Load metadata linux/amd64
            +go-base | --> FROM golang:1.19-bullseye
             context | --> local context .
             context | --> local context .
            +go-base | [          ]   0% resolve docker.io/library/golang:1.19-bullseye@sha256:d92ddd8ad9d960c67dc34cffc2ed7b0ef399be2549505bf5ef94a7f4ca016a05    [██████████] 100% resolve docker.io/library/golang:1.19-bullseye@sha256:d92ddd8ad9d960c67dc34cffc2ed7b0ef399be2549505bf5ef94a7f4ca016a05
            +go-base | [██████████] 100% sha256:8709771bd9da550643f5f4e3b49e92bb3f90543507ff36b5a998dd461fb8dd28
            +go-base | [██████████] 100% sha256:a42821cd14fb31c4aa253203e7f8e34fc3b15d69ce370f1223fbbe4252a64202
            +go-base | [          ]   0% transferring .:
             context | transferred 3 file(s) for context . (97 kB, 3 file/dir stats)
             context | [██████████] 100% transferring .:
             context | [██████████] 100% transferring .:
            +go-base | [██████████] 100% sha256:326f452ade5c33097eba4ba88a24bd77a93a3d994d4dc39b936482655e664857
            +go-base | [██████████] 100% sha256:b7aa120dd02d9fa75bb50861103f7837514028ea6f06e3e821b8c47c2c10d386
            +go-base | [██████████] 100% sha256:8471b75885efc7790a16be5328e3b368567b76a60fc3feabd6869c15e175ee05
            +go-base | [██████████] 100% sha256:23858da423a6737f0467fab0014e5b53009ea7405d575636af0c3f100bbb2f10
            +go-base | [██████████] 100% extracting sha256:23858da423a6737f0467fab0014e5b53009ea7405d575636af0c3f100bbb2f10
            +go-base | [██████████] 100% extracting sha256:326f452ade5c33097eba4ba88a24bd77a93a3d994d4dc39b936482655e664857
            +go-base | [██████████] 100% extracting sha256:a42821cd14fb31c4aa253203e7f8e34fc3b15d69ce370f1223fbbe4252a64202
            +go-base | [██████████] 100% extracting sha256:8471b75885efc7790a16be5328e3b368567b76a60fc3feabd6869c15e175ee05
            +go-base | [██████████] 100% sha256:292167c9d1ff24956858651ef9906e9a987b65f7362854e13c28b98d9ae4e09f
            +go-base | [██████████] 100% extracting sha256:b7aa120dd02d9fa75bb50861103f7837514028ea6f06e3e821b8c47c2c10d386
            +go-base | [          ]   0% extracting sha256:292167c9d1ff24956858651ef9906e9a987b65f7362854e13c28b98d9ae4e09f
             ongoing | +go-base (5 seconds ago)
            +go-base | [██████████] 100% extracting sha256:292167c9d1ff24956858651ef9906e9a987b65f7362854e13c28b98d9ae4e09f
            +go-base | [██████████] 100% extracting sha256:8709771bd9da550643f5f4e3b49e92bb3f90543507ff36b5a998dd461fb8dd28
            +go-base | --> WORKDIR /ngx
            +go-base | --> COPY ngx.go ngx_test.go ngx_internal_test.go .
            +go-base | --> COPY go.mod go.sum .
            +go-base | --> RUN go mod download
            +go-test | --> RUN go install github.com/mfridman/tparse@latest
            +go-test | go: downloading github.com/mfridman/tparse v0.11.1
            +go-test | go: downloading github.com/charmbracelet/lipgloss v0.4.0
            +go-test | go: downloading github.com/muesli/termenv v0.11.0
            +go-test | go: downloading github.com/olekukonko/tablewriter v0.0.5
            +go-test | go: downloading github.com/lucasb-eyer/go-colorful v1.2.0
            +go-test | go: downloading github.com/mattn/go-runewidth v0.0.13
            +go-test | go: downloading github.com/muesli/reflow v0.3.0
            +go-test | go: downloading github.com/mattn/go-isatty v0.0.14
            +go-test | go: downloading golang.org/x/sys v0.0.0-20220513210249-45d2b4557a2a
            +go-test | go: downloading github.com/rivo/uniseg v0.2.0
            +go-test | --> RUN go test -count=1 -shuffle=on -trimpath -race -cover -covermode=atomic -json ./... | tparse -all
             ongoing | +go-test (14 seconds ago)
            +go-test | ┌───────────────────────────────────────────────────────────────────────────────────────────┐
            +go-test | │  STATUS │ ELAPSED │                            TEST                            │ PACKAGE  │
            +go-test | │─────────┼─────────┼────────────────────────────────────────────────────────────┼──────────│
            +go-test | │  PASS   │    0.00 │ TestClientRetrivesInfoAboutRunningNGINXInstance            │ ngx      │
            +go-test | │  PASS   │    0.00 │ TestBuildsAddressOnValidInputWithHostAndPort               │ ngx      │
            +go-test | │  PASS   │    0.00 │ TestBuildsAddressOnValidInputWithUnixSocket                │ ngx      │
            +go-test | │  PASS   │    0.00 │ TestBuildsAddressOnValidInputWithIPV6AddressAndWithoutPort │ ngx      │
            +go-test | │  PASS   │    0.00 │ TestCalculatesStreamServerUpdatesOnValidInput              │ ngx      │
            +go-test | │  PASS   │    0.00 │ TestDeterminesUpstreamServersConfigurationEquality         │ ngx      │
            +go-test | │  PASS   │    0.00 │ TestClientRetrivesNGINXStatusOnValidParameters             │ ngx      │
            +go-test | │  PASS   │    0.00 │ TestDeterminesUpstreamStreamServersConfigurationEquality   │ ngx      │
            +go-test | │  PASS   │    0.00 │ TestBuildsAddressOnValidInputWithIPV6AndPort               │ ngx      │
            +go-test | │  PASS   │    0.00 │ TestBuildsAddressOnValidInputWithIPV4AddressAndWithoutPort │ ngx      │
            +go-test | │  PASS   │    0.00 │ TestCalculatesServerUpdatesOnValidInput                    │ ngx      │
            +go-test | │  PASS   │    0.00 │ TestBuildsAddressOnValidInputWithIPV4AndPort               │ ngx      │
            +go-test | │  PASS   │    0.00 │ TestBuildsAddressOnValidInputWithAddressAndWithoutPort     │ ngx      │
            +go-test | └───────────────────────────────────────────────────────────────────────────────────────────┘
            +go-test | ┌────────────────────────────────────────────────────────────────────────┐
            +go-test | │  STATUS │ ELAPSED │       PACKAGE        │ COVER │ PASS │ FAIL │ SKIP  │
            +go-test | │─────────┼─────────┼──────────────────────┼───────┼──────┼──────┼───────│
            +go-test | │  PASS   │  0.04s  │ github.com/qba73/ngx │ 21.4% │  13  │  0   │  0    │
            +go-test | └────────────────────────────────────────────────────────────────────────┘
              output | --> exporting outputs

 3. Push ⏫ (disabled)
————————————————————————————————————————————————————————————————————————————————

To enable pushing use earthly --push

 4. Local Output 🎁
————————————————————————————————————————————————————————————————————————————————



========================== 🌍 Earthly Build  ✅ SUCCESS ==========================

</details>


## Contributing

If you have any suggestions or experience issues with the NGINX Plus Go Client, please create an issue or send a pull request on GitHub.
