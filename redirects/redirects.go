package redirects

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/network"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"
)

func GetRedirects(url, chromeRemoteDebuggingUrl string, maxTimeToRedirect time.Duration) (map[string][]string, error) {
	links, err := getLinks(url)
	if err != nil {
		return nil, err
	}

	redirects := make(map[string][]string)
	for _, link := range links {
		linkRedirects, err := GetLinkRedirects(link, chromeRemoteDebuggingUrl, maxTimeToRedirect)
		if err != nil {
			return nil, err
		}
		redirects[link] = linkRedirects
	}

	return redirects, nil
}

func getLinks(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error getting the HTTP resource %s: %v", url, err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body of the HTTP resource %s: %v", url, err)
	}

	var links struct {
		Links []string `json:"links"`
	}
	if err := json.Unmarshal(respBody, &links); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON from the HTTP resource %s: %v", url, err)
	}

	return links.Links, nil
}

func GetLinkRedirects(link, chromeRemoteDebuggingUrl string, maxTimeToRedirect time.Duration) ([]string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Use the DevTools HTTP/JSON API to manage targets (e.g. pages, webworkers).
	devTools := devtool.New(chromeRemoteDebuggingUrl)
	pt, err := devTools.Get(ctx, devtool.Page)
	if err != nil {
		return nil, err
	}

	// Initiate a new RPC connection to the Chrome Debugging Protocol target.
	conn, err := rpcc.DialContext(ctx, pt.WebSocketDebuggerURL)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	c := cdp.NewClient(conn)

	// Open a RequestWillBeSent client to buffer this event.
	requestWillBeSent, err := c.Network.RequestWillBeSent(ctx)
	if err != nil {
		return nil, err
	}
	defer requestWillBeSent.Close()

	// Enable network events.
	if err = c.Network.Enable(ctx, network.NewEnableArgs()); err != nil {
		return nil, err
	}

	// Navigate to a given page.
	_, err = c.Page.Navigate(ctx, page.NewNavigateArgs(link))
	if err != nil {
		return nil, err
	}

	// Get redirects.
	var redirects []string
	for {
		redirected := make(chan bool, 1)
		// Wait for a redirect at most `maxTimeToRedirect`.
		go func(redirected chan bool) {
			select {
			case <-redirected:
			case <-time.After(maxTimeToRedirect):
				cancel()
			}
		}(redirected)

		// Wait until a `RequestWillBeSent` event is received.
		requestWillBeSentReply, err := requestWillBeSent.Recv()
		redirected <- true
		if ctx.Err() == context.Canceled {
			break
		}
		if err != nil {
			return nil, err
		}

		redirects = append(redirects, requestWillBeSentReply.DocumentURL)
	}

	return redirects, nil
}
