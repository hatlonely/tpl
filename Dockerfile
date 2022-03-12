
FROM golang:1.16-alpine AS build

ARG version

COPY . /go/src/
WORKDIR /go/src/
RUN go build -ldflags "-X 'main.Version=$version'" -o build/bin/tpl cmd/main.go

FROM alpine:3.15.0

COPY --from=build /go/src/build/bin/tpl /usr/bin/
