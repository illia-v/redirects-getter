package main

import (
	"flag"
	"redirects_getter/redirects"
)

func main() {
	url := flag.String(
		"resource-url",
		"http://docker-host1.cli.bz:8888/test-links.json",
		"URL of an HTTP resource where links are stored.",
	)
	flag.Parse()

	_, err := redirects.GetRedirects(*url)
	if err != nil {
		panic(err)
	}
}
