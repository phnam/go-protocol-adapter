package main

import (
	"fmt"

	"github.com/phnam/go-protocol-adapter/client"
	"github.com/phnam/go-protocol-adapter/common"
	"github.com/phnam/go-protocol-adapter/request"
)

type ClientData struct {
	Message string "json:'message'"
	Code    int    "json:'code'"
}

func main3333() {
	cli := client.NewThriftClient[ClientData](&client.APIClientConfiguration{
		Address:              "localhost:8080",
		Timeout:              1000,
		MaxRetry:             1,
		WaitToRetry:          100,
		MaxConnection:        10,
		KeepDataStringFormat: nil,
		Protocol:             common.Protocol.THRIFT,
		ErrorLogOnly:         false,
	})

	resp := cli.MakeRequest(&request.OutboundAPIRequest{
		Method: "GET",
		Path:   "/",
	})

	fmt.Println(resp)
	fmt.Println(resp.Data[0].Code)
	fmt.Println(resp.Data[0].Message)
}
