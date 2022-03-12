package rpcx

var ConfigAppJson = `{
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
	"enableTrace": false,
	"enableMetric": false,
	"enablePprof": false,
	"enableCors": true,
	"cors": {
	  "allowAll": true,
	  "allowMethod": ["GET, HEAD, POST, PUT, DELETE"],
	}
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
			  "pascalNameKey": false
			}
		  }
		}
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
			  "pascalNameKey": false
			}
		  }
		}
	  }]
	}
  }
}
`
