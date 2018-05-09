package main

import (
	"flag"
	"encoding/json"
	"redirects_getter/redirects"
	"fmt"
	"time"
)

func main() {
	url := flag.String(
		"resource-url",
		"http://docker-host1.cli.bz:8888/test-links.json",
		"URL of an HTTP resource where links are stored.",
	)
	chromeRemoteDebuggingUrl := flag.String(
		"chrome-remote-debugging-url",
		"http://chrome:9222",
		"URL of a Chrome instance for debugging.",
	)
	maxTimeToRedirect := flag.Duration(
		"max-time-to-redirect",
		5 * time.Second,
		"Maximum time that you allow the tool to wait for a redirect. " +
		`Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".`,
	)
	flag.Parse()

	r, err := redirects.GetRedirects(*url, *chromeRemoteDebuggingUrl, *maxTimeToRedirect)
	if err != nil {
		panic(err)
	}
	rJson, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(rJson))
}
