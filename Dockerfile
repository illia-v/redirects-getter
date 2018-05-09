FROM golang:1.10.2-alpine3.7

# Install dep (Go dependency management tool).
RUN apk add --no-cache git \
    && wget -O /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 \
    && chmod +x /usr/local/bin/dep

RUN mkdir -p /go/src/redirects_getter
WORKDIR /go/src/redirects_getter

# Install dependencies.
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only

# Build project.
COPY . .
RUN go build -o get-redirects
