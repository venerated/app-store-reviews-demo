package main

import (
	"fmt"
	"net/http"
)

type CustomClient func(url string, customAddedHeaders CustomHeaders) (*http.Response, error)

type CustomHeaders map[string]string

func customHttpClient(url string, customAddedHeaders CustomHeaders) (*http.Response, error) {
	client := &http.Client{}
	req, reqErr := http.NewRequest("GET", url, nil)
	if reqErr != nil {
		return nil, fmt.Errorf("error in getAppStoreRss, request creation failed: %w", reqErr)
	}
	for key, value := range customAddedHeaders {
		req.Header.Add(key, value)
	}
	resp, respError := client.Do(req)
	if respError != nil {
		return nil, fmt.Errorf("error in getAppStoreRss, response failed: %w", respError)
	}

	// Caller is responsible for closing the response body.
	return resp, nil
}
