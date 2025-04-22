package main

import (
	"github.com/phnam/go-protocol-adapter/common"
	"github.com/phnam/go-protocol-adapter/request"
	"github.com/phnam/go-protocol-adapter/responder"
	"github.com/phnam/go-protocol-adapter/server"
)

type ServerData struct {
	Message string "json:'message'"
	Code    int    "json:'code'"
}

func main333() {

	// different init
	server := server.NewServer(server.ServerConfig{
		Protocol: common.Protocol.THRIFT,
	})

	// same body
	server.SetHandler(common.APIMethod.GET, "/", func(req request.APIRequest, res responder.APIResponder) error {
		return res.Respond(&common.APIResponse[any]{
			Status:  common.APIStatus.Ok,
			Message: "This is a message",
			Data: []any{
				ServerData{
					Message: "Hello world 1234",
					Code:    200,
				},
			},
		})
	})
	server.Expose(8080)
	server.Start(nil)
}
