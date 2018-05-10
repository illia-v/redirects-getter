package redirects

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/network"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"
)

type Redirects struct {
	Link      string   `json:"link"`
	Redirects []string `json:"redirects"`
}

func GetRedirects(url, chromeRemoteDebuggingUrl string, maxTimeToRedirect time.Duration) ([]Redirects, error) {
	links, err := getLinks(url)
	if err != nil {
		return nil, err
	}

	// Set up a wait group to match the number of links.
	var wg sync.WaitGroup
	wg.Add(len(links))

	errChan := make(chan error, 1) // A channel for handling errors.
	var redirects []Redirects

	// Get redirects of a link concurrently.
	for _, link := range links {
		go func(link string) {
			defer wg.Done()
			linkRedirects, err := GetLinkRedirects(link, chromeRemoteDebuggingUrl, maxTimeToRedirect)
			if err != nil {
				// Add an error to the error channel.
				errChan <- err
				return
			}
			redirects = append(redirects, linkRedirects)
		}(link)
	}

	// Finish when all links are processed.
	finishedChan := make(chan bool, 1)
	go func() {
		wg.Wait()
		finishedChan <- true
	}()

	// Wait until finish or an error.
	select {
	case <-finishedChan:
	case err := <-errChan:
		if err != nil {
			return nil, err
		}
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

func GetLinkRedirects(link, chromeRemoteDebuggingUrl string, maxTimeToRedirect time.Duration) (Redirects, error) {
	redirects := Redirects{
		Link: link,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Use the DevTools HTTP/JSON API to manage targets (e.g. pages, webworkers).
	devTools := devtool.New(chromeRemoteDebuggingUrl)
	// Create a new target (tab).
	pt, err := devTools.Create(ctx)
	if err != nil {
		return redirects, err
	}
	defer devTools.Close(ctx, pt)

	// Initiate a new RPC connection to the Chrome Debugging Protocol target.
	conn, err := rpcc.DialContext(ctx, pt.WebSocketDebuggerURL)
	if err != nil {
		return redirects, err
	}
	defer conn.Close()

	c := cdp.NewClient(conn)

	// Open a RequestWillBeSent client to buffer this event.
	requestWillBeSent, err := c.Network.RequestWillBeSent(ctx)
	if err != nil {
		return redirects, err
	}
	defer requestWillBeSent.Close()

	// Enable network events.
	if err = c.Network.Enable(ctx, network.NewEnableArgs()); err != nil {
		return redirects, err
	}

	// Navigate to a given page.
	_, err = c.Page.Navigate(ctx, page.NewNavigateArgs(link))
	if err != nil {
		return redirects, err
	}

	// Wait for the first request.
	requestWillBeSent.Recv()

	// Wait at most `maxTimeToRedirect` for a redirect.
	redirectDeadline := time.After(maxTimeToRedirect)
	go func() {
		<-redirectDeadline
		cancel()
	}()

	// Get redirects.
	for {
		// Wait until a `RequestWillBeSent` event is received.
		requestWillBeSentReply, err := requestWillBeSent.Recv()
		if err != nil {
			// Context can be cancelled if `maxTimeToRedirect` exceeds.
			if ctx.Err() == context.Canceled {
				break
			}
			return redirects, err
		}
		// Process only requests with the "Document" type.
		if *requestWillBeSentReply.Type == "Document" {
			redirectDeadline = time.After(maxTimeToRedirect) // Reset time to wait for next redirect.
			redirects.Redirects = append(redirects.Redirects, requestWillBeSentReply.DocumentURL)
		}
	}

	return redirects, nil
}
