[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/qba73/ngx?logo=go)

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

## Using the Client

```go
import "github.com/qba73/ngx"
```

## Testing

Prerequisites:
* Docker
* Go
* Make
* NGINX Plus license - put `nginx-repo.crt` and `nginx-repo.key` into the `docker` folder.

Run Tests:

```
$ make docker-build && make test
```

This will build and run two NGINX Plus containers and create one docker network of type bridge, execute the client tests against both NGINX Plus APIs, and then clean up. If it fails and you want to clean up (i.e. stop the running containers and remove the docker network), please use `$ make clean`


## Contributing

If you have any suggestions or experience issues with the NGINX Plus Go Client, please create an issue or send a pull request on GitHub.