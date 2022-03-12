package rpcx

var dockerfile = `{{ if eq .ImageType "centos" }}
FROM golang:1.16 AS build

ARG version

COPY . /go/src/
WORKDIR /go/src/
RUN mkdir -p build/bin && mkdir -p build/config && cp config/* build/config && \
    go build -ldflags "-X 'main.Version=$version'" -o build/bin/{{ .Name }} cmd/main.go

FROM centos:centos7
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" >> /etc/timezone

COPY --from=build /go/src/build /work/{{ .Name }}
WORKDIR /work/{{ .Name }}
CMD [ "bin/{{ .Name }}", "-c", "config/base.json" ]

{{ else if eq .ImageType "alpine" }}

FROM golang:1.16-alpine AS build

ARG version

COPY . /go/src/
WORKDIR /go/src/
RUN mkdir -p build/bin && mkdir -p build/config && cp config/* build/config && \
    go build -ldflags "-X 'main.Version=$version'" -o build/bin/{{ .Name }} cmd/main.go

FROM alpine:3.15.0

COPY --from=build /go/src/build /work/{{ .Name }}
WORKDIR /work/{{ .Name }}
CMD [ "bin/{{ .Name }}", "-c", "config/base.json" ]
{{ end }}
`
