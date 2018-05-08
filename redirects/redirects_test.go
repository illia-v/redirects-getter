package redirects

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestGetRedirects(t *testing.T) {
	// Mock HTTP request.
	defer gock.Off()
	gock.New("http://example.com").
		Get("/").
		Reply(200).
		BodyString(`{"links": ["http://example.com", "http://example.net"]}`)

	_, err := GetRedirects("http://example.com")
	assert.NoError(t, err)
}

func TestInvalidURL(t *testing.T) {
	_, err := GetRedirects("not-url")
	assert.Error(t, err)
}

func TestInvalidResourceNotJson(t *testing.T) {
	// Mock HTTP request.
	defer gock.Off()
	gock.New("http://example.com").
		Get("/").
		Reply(200).
		BodyString("{")

	_, err := GetRedirects("http://example.com")
	assert.Error(t, err)
}
