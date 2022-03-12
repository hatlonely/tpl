package rpcx

var opsHelmValuesAdapterYaml = `
name: "${NAME}"
replicaCount: "${REPLICA_COUNT}"

image:
  repository: "${REGISTRY_ENDPOINT}/${REGISTRY_NAMESPACE}/${NAME}"
  tag: "${VERSION}"
  pullPolicy: Always
  pullSecret: "${PULL_SECRET_NAME}"

ingress:
  enable: true
  host: "${INGRESS_HOST}"
  secretName: "${SECRET_NAME}"

config:
  base.json: |
    {
      "decoder": {
        "type": "Json"
      },
      "provider": {
        "type": "Local",
        "options": {
          "filename": "config/app.json"
        }
      }
    }
  app.json: |
    {
      "grpcGateway": {
        "httpPort": 80,
        "grpcPort": 6080,
        "exitTimeout": "20s",
        "validators": [
          "Default"
        ],
        "usePascalNameLogKey": false,
        "usePascalNameErrKey": false,
        "marshalUseProtoNames": true,
        "marshalEmitUnpopulated": false,
        "unmarshalDiscardUnknown": true,
        "enablePing": true,
        "enableMetric": false,
        "enablePprof": false,{{ if .Ops.EnableTrace }}
        "enableTrace": true,
        "jaeger": {
          "serviceName": "{{ .Name }}",
          "sampler": {
            "type": "const",
            "param": 1,
            "samplingServerURL": "${JAEGER_SAMPLING_SERVER_URL}"
          },
          "reporter": {
            "logSpans": false,
            "localAgentHostPort": "${JAEGER_REPORTER_LOCAL_AGENT_HOST_PORT}"
          }
        },{{ end }}{{ if .Ops.EnableCors }}
        "enableCors": true,
        "cors": {
          "allowAll": true,
          "allowMethod": ["GET, HEAD, POST, PUT, DELETE"],
        }{{ end }}
      },
      "service": {
      },
      "logger": {
        "grpc": {
          "level": "Info",
          "writers": [{
            "type": "RotateFile",
            "options": {
              "filename": "log/app.rpc",
              "maxAge": "24h",
              "formatter": {
                "type": "Json",
                "options": {
                  "flatMap": true,
                  "pascalNameKey": true
                }
              }
            }{{ if .Ops.EnableEsLog }}
          }, {
            "type": "ElasticSearch",
            "options": {
              "index": "{{ .Name }}-rpc",
              "idField": "requestID",
              "timeout": "200ms",
              "msgChanLen": 200,
              "workerNum": 2,
              "es": {
                "es": {
                  "uri": "${ELASTICSEARCH_ENDPOINT}",
                  "username": "elastic",
                  "password": "${ELASTICSEARCH_PASSWORD}"
                },
                "retry": {
                  "attempt": 3,
                  "delay": "1s",
                  "lastErrorOnly": true,
                  "delayType": "BackOff"
                }
              }
            }{{ end }}
          }]
        },
        "info": {
          "level": "Info",
          "writers": [{
            "type": "RotateFile",
            "options": {
              "filename": "log/app.log",
              "maxAge": "24h",
              "formatter": {
                "type": "Json",
                "options": {
                  "pascalNameKey": true
                }
              }
            }{{ if .Ops.EnableEsLog }}
          }, {
            "type": "ElasticSearch",
            "options": {
              "index": "{{ .Name }}-log",
              "idField": "requestID",
              "timeout": "200ms",
              "msgChanLen": 200,
              "workerNum": 2,
              "es": {
                "es": {
                  "uri": "${ELASTICSEARCH_ENDPOINT}",
                  "username": "elastic",
                  "password": "${ELASTICSEARCH_PASSWORD}"
                },
                "retry": {
                  "attempt": 3,
                  "delay": "1s",
                  "lastErrorOnly": true,
                  "delayType": "BackOff"
                }
              }
            }{{ end }}
          }]
        }
      }
    }
`
