package rpcx

var apiProto = `
syntax = "proto3";

package api;

option go_package = "{{ .Package }}/api/gen/go/api";

import "google/protobuf/empty.proto";
import "google/api/annotations.proto";
import "tagger/tagger.proto";

message EchoReq {
  string message = 1;
}

message EchoRes {
  string message = 2;
}

message AddReq {
  int32 i1 = 1 [(tagger.tags) = "rule:\"x >= 0 && x <= 100\""];
  int32 i2 = 2 [(tagger.tags) = "rule:\"x >= 0 && x <= 100\""];
}

message AddRes {
  int32 val = 1;
}

message PingRes {
  string Message = 1;
}

service ExampleService {
  rpc Echo(EchoReq) returns (EchoRes) {
    option (google.api.http) = {
      get: "/v1/echo"
    };
  }

  rpc Add(AddReq) returns (AddRes) {
    option (google.api.http) = {
      post: "/v1/add"
      body: "*"
    };
  }

  rpc Ping(google.protobuf.Empty) returns (PingRes) {
    option (google.api.http) = {
      get: "/ping"
    };
  }
}
`
