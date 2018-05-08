package redirects

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetRedirects(url string) (interface{}, error) {
	_, err := getLinks(url)
	if err != nil {
		return nil, err
	}

	return nil, nil
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
