VERSION 0.6
FROM golang:1.19-bullseye
ARG GO_IMAGE=golang:1.19-bullseye
ARG WORKDIR=/ngx

configure:
  LOCALLY
  RUN git config pull.rebase true \
    && git config remote.origin.prune true \
    && git config branch.main.mergeoptions "--ff-only"

checks:
  #BUILD +go-linter
  BUILD +go-test

go-base:
  FROM $GO_IMAGE
  WORKDIR $WORKDIR
  COPY ngx.go ngx_test.go ngx_internal_test.go .
  COPY go.mod go.sum .
  RUN go mod download

go-test:
  FROM +go-base
  RUN go install github.com/mfridman/tparse@latest
  RUN go test -count=1 -shuffle=on -trimpath -race -cover -covermode=atomic -json ./... | tparse -all

go-linter:
  FROM +go-base
  WORKDIR $WORKDIR
  RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  COPY ngx.go ngx_test.go ngx_internal_test.go .
  COPY go.mod go.sum .
  RUN golangci-lint run 

