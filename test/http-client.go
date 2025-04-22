package main

import (
	"fmt"
	"time"

	"github.com/phnam/go-protocol-adapter/client"
	"github.com/phnam/go-protocol-adapter/common"
	"github.com/phnam/go-protocol-adapter/request"
)

type HTTPClientData struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func main() {
	cli := client.NewAPIClient[HTTPClientData](&client.APIClientConfiguration{
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

	fmt.Println(resp)
	if resp.Status == "OK" {
		fmt.Println(resp.Data[0].Code)
		fmt.Println(resp.Data[0].Message)
	}

}
