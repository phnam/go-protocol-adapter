package main

import (
	"github.com/phnam/go-protocol-adapter/common"
	"github.com/phnam/go-protocol-adapter/request"
	"github.com/phnam/go-protocol-adapter/responder"
	"github.com/phnam/go-protocol-adapter/server"
)

type HTTPServerData struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func main33() {

	// different init
	server := server.NewServer(server.ServerConfig{
		Protocol: common.Protocol.HTTP,
	})

	// same body
	server.SetHandler(common.APIMethod.GET, "/", func(req request.APIRequest, res responder.APIResponder) error {
		return res.Respond(&common.APIResponse[any]{
			Status:  common.APIStatus.Ok,
			Message: "Hello world",
			Data:    []any{HTTPServerData{Message: "Hello world from HTTP Server", Code: 201}},
		})
	})
	server.Expose(80)
	server.Start(nil)
}
