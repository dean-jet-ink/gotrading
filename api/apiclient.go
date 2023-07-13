package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type APIClient struct {
	URL         *url.URL
	HTTPClient  *http.Client
	Key, Secret string
}

func (c *APIClient) NewRequest(ctx context.Context, urlPath, method string, queryParams map[string]string, header map[string]string, body []byte) (*http.Request, error) {
	ref, err := url.Parse(urlPath)
	if err != nil {
		return nil, err
	}
	endpoint := c.URL.ResolveReference(ref).String()

	req, err := http.NewRequestWithContext(ctx, method, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	for key, value := range header {
		req.Header.Add(key, value)
	}

	queries := req.URL.Query()
	for key, value := range queryParams {
		queries.Add(key, value)
	}
	req.URL.RawQuery = queries.Encode()

	return req, nil
}

func (c *APIClient) DecodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}
