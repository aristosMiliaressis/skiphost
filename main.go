package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/aristosMiliaressis/httpc/pkg/httpc"
)

type Result struct {
	LooksGood bool   `json:"looks_good"`
	Reason    string `json:"reason,omitempty"`
}

func main() {
	var targetUrl string

	flag.StringVar(&targetUrl, "u", "", "Target url.")
	flag.Parse()

	client := httpc.NewHttpClient(httpc.DefaultOptions, context.Background())

	baseReq, _ := http.NewRequest("GET", targetUrl, nil)

	notFoundReq, _ := http.NewRequest("GET", httpc.ToAbsolute(targetUrl, "/jhgfsdfgjfgaskjfg"), nil)

	i := 0

	var baseResp *httpc.MessageDuplex
	for {
		baseResp = client.Send(baseReq)
		<-baseResp.Resolved
		if baseResp.Response != nil {
			break
		}
		i++
		if i == 3 {
			os.Exit(1)
		}
	}

	i = 0

	var notFoundResp *httpc.MessageDuplex
	for {
		notFoundResp = client.Send(notFoundReq)
		<-notFoundResp.Resolved
		if notFoundResp.Response != nil {
			break
		}
		i++
		if i == 3 {
			os.Exit(1)
		}
	}

	result, _ := json.Marshal(Result{LooksGood: true})

	if baseResp.Response.StatusCode == 502 && notFoundResp.Response.StatusCode == 502 {
		result, _ = json.Marshal(Result{LooksGood: false, Reason: "502"})
	} else if baseResp.Response.StatusCode == 503 && notFoundResp.Response.StatusCode == 503 {
		result, _ = json.Marshal(Result{LooksGood: false, Reason: "503"})
	} else if baseResp.Response.StatusCode == 504 && notFoundResp.Response.StatusCode == 504 {
		result, _ = json.Marshal(Result{LooksGood: false, Reason: "504"})
	} else if baseResp.Response.StatusCode >= 300 && baseResp.Response.StatusCode < 400 &&
		notFoundResp.Response.StatusCode >= 300 && notFoundResp.Response.StatusCode < 400 {
		redirectUrl := httpc.GetRedirectLocation(notFoundResp.Response)
		if httpc.IsCrossOrigin(targetUrl, redirectUrl) {
			result, _ = json.Marshal(Result{LooksGood: false, Reason: "cross-origin-redirect"})
		}
	}

	fmt.Println(string(result))
}
