# redirects
Package redirects provides the ability to get chains of redirects from Chrome
using Chrome Debugging Protocol.

## Usage

#### type Redirects

```go
type Redirects struct {
        Link      string   `json:"link"`
        Redirects []string `json:"redirects"`
}
```


#### func  GetLinkRedirects

```go
func GetLinkRedirects(link, chromeRemoteDebuggingUrl string, maxTimeToRedirect time.Duration) (Redirects, error)
```
GetLinkRedirects returns redirects of one given page.

#### func  GetRedirects

```go
func GetRedirects(url, chromeRemoteDebuggingUrl string, maxTimeToRedirect time.Duration) ([]Redirects, error)
```
GetRedirects returns redirects of a couple of pages. It gets links from a given
HTTP resource.
