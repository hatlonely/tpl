package rpcx

var dockerfile = `
FROM golang:1.16 AS build

COPY . /go/src/
WORKDIR /go/src/
RUN make build

FROM centos:centos7
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" >> /etc/timezone

COPY --from=build /go/src/build /work/{{ .Name }}
WORKDIR /work/{{ .Name }}
CMD [ "bin/{{ .Name }}", "-c", "config/base.json" ]
`
