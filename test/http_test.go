package main

import (
	"strconv"
	"testing"
	"time"

	"github.com/phnam/go-protocol-adapter/client"
	"github.com/phnam/go-protocol-adapter/common"
	"github.com/phnam/go-protocol-adapter/request"
	"github.com/phnam/go-protocol-adapter/responder"
	"github.com/phnam/go-protocol-adapter/server"
)

func TestHTTPServer(t *testing.T) {
	type HTTPServerData struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}

	// init HTTP server
	server := server.NewServer(server.ServerConfig{
		Protocol: common.Protocol.HTTP,
	})
	server.SetHandler(common.APIMethod.GET, "/", func(req request.APIRequest, res responder.APIResponder) error {
		return res.Respond(&common.APIResponse[any]{
			Status:  common.APIStatus.Ok,
			Message: "Hello world",
			Data:    []any{HTTPServerData{Message: "Hello world from HTTP Server", Code: 201}},
		})
	})
	server.Expose(80)
	go server.Start(nil)

	// wait for startup
	time.Sleep(1000 * time.Millisecond)

	// call API
	cli := client.NewAPIClient[HTTPServerData](&client.APIClientConfiguration{
		Address:              "localhost",
		Timeout:              100 * time.Millisecond,
		MaxRetry:             1,
		WaitToRetry:          100,
		MaxConnection:        10,
		KeepDataStringFormat: nil,
		Protocol:             common.Protocol.HTTP,
		ErrorLogOnly:         false,
	})

	resp := cli.MakeRequest(&request.OutboundAPIRequest{
		Method: "GET",
		Path:   "/",
	})

	if resp.Status != common.APIStatus.Ok {
		t.Error("HTTP Server test failed. Wrong status: " + resp.Status)
	}

	if resp.Data == nil || len(resp.Data) == 0 || resp.Data[0].Code != 201 {
		t.Error("HTTP Server test failed. Wrong data: " + strconv.Itoa(resp.Data[0].Code))
	}

}
