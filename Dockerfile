
FROM golang:1.16-alpine AS build

COPY . /go/src/
WORKDIR /go/src/
RUN make build

FROM alpine:3.15.0

COPY --from=build /go/src/build/bin/tpl /usr/bin/
