package rpcx

var ConfigBaseJson = `
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
`
