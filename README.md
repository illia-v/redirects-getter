# Redirects Getter
A highly concurrent tool for getting chains of redirects from Chrome.

It is solution for a task of Round 1 of
[DEV Challenge 12](https://devchallenge.it/) (standard level).

## Prerequisite
The project uses Docker and Docker Compose.
You should have them installed.

## Usage
Run `docker-compose up --build`. The command will take care of running
a browser instance, building and running the project.

#### CLI arguments
```
  --chrome-remote-debugging-url string
        URL of a Chrome instance for debugging.
        (default "http://chrome:9222")
  --max-time-to-redirect duration
        Maximum time that you allow the tool to wait for a redirect.
        Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
        (default 5s)
  --resource-url string
        URL of an HTTP resource where links are stored.
        (default "http://docker-host1.cli.bz:8888/test-links.json")
```
You should use the `REDIRECTS_GETTER_ARGS` environment variable to
provide the arguments to a Docker container. Example:
```bash
REDIRECTS_GETTER_ARGS="--resource-url=http://example.com --max-time-to-redirect=60s" \
    docker-compose up --build
```

#### Example of a response of the HTTP resourse where links are stored:
```json
{
    "links": [
        "http://example.com",
        "http://example.net"
    ]
}
```

#### Example of output of the tool:
```json
[
     {
         "link": "http://example.com",
         "redirects": [
             "https://example.com",
             "https://example.com/first-redirect",
             "https://example.com/second-redirect",
             "https://example.com/third-redirect"
         ]
     },
     {
         "link": "http://example.net",
         "redirects": [
             "http://example.net/first-redirect",
             "http://example.net/last-redirect"
         ]
     }
]
```

## Usage as a Package
You can use `redirects` as a separate Go package.
See [redirects/README.md](redirects/README.md) for a basic
documentation.
