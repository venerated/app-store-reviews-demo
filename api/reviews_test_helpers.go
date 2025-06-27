package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var testAppId = "595068606"

var testCacheFilePath = "cache/reviews_" + testAppId + ".json"

var testEntry = Entry{
	ID: Label{
		Label: "12802724742",
	},
	Author: Author{
		Label{
			Label: "ZenZen2025",
		},
	},
	Content: Label{
		Label: "We had multiple dinners with ten people and 5 payers and this app worked great.",
	},
	Rating: Label{
		Label: "5",
	},
	Updated: Label{
		Label: "2025-06-21T21:09:29-07:00",
	},
}

var testReview = Review{
	ID:      "12802724742",
	Author:  "ZenZen2025",
	Content: "We had multiple dinners with ten people and 5 payers and this app worked great.",
	Rating:  "5",
	Updated: "2025-06-21T21:09:29-07:00",
}

var testRssJson = `
	{
		"feed": {
			"entry": [
				{
					"id": { "label": "12802724742" },
					"author": {
						"name": { "label": "ZenZen2025" }
					},
					"content": { "label": "We had multiple dinners with ten people and 5 payers and this app worked great." },
					"im:rating": { "label": "5" },
					"updated": { "label": "2025-06-21T21:09:29-07:00" }
				}
			],
			"updated": { "label": "2025-06-26T10:06:25-07:00" }
		}
	}`

var testRssStructData = RssResponse{
	Feed: Feed{
		Entry: []Entry{
			{
				ID: Label{
					Label: "12802724742",
				},
				Author: Author{
					Label{
						Label: "ZenZen2025",
					},
				},
				Content: Label{
					Label: "We had multiple dinners with ten people and 5 payers and this app worked great.",
				},
				Rating: Label{
					Label: "5",
				},
				Updated: Label{
					Label: "2025-06-21T21:09:29-07:00",
				},
			},
		},
		Updated: Label{
			Label: "2025-06-26T10:06:25-07:00",
		},
	},
}

func mockCustomClient(url string, customAddedHeaders CustomHeaders) (*http.Response, error) {
	return mockHTTPResponse(testRssJson, 200), nil
}

func mockHTTPResponse(body string, status int) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

func mockGetAppStoreRssSuccess(lastUpdated string, deps GetAppStoreRssDeps) (RssResponse, error) {
	return testRssStructData, nil
}

func mockGetCachedReviewsError(filePath string, readFromCacheFunc ReadRssFromCache) (RssResponse, error) {
	return RssResponse{}, errors.New("error in getCachedReviews")
}

func mockGetCachedReviewsSuccess(filePath string, readFromCacheFunc ReadRssFromCache) (RssResponse, error) {
	return testRssStructData, nil
}

func mockGetUpdatedData(cacheFilePath string, resp *http.Response, saveToCacheFunc SaveToCache) (RssResponse, error) {
	if resp.StatusCode >= 400 {
		return RssResponse{}, fmt.Errorf("external API returned status %d", resp.StatusCode)
	}
	return testRssStructData, nil
}

func mockReadFromCacheFileDoesntExistOrUnreadable(fileName string) (RssResponse, error) {
	return RssResponse{}, errors.New("file doesn't exist or is unreadable")
}

func mockReadFromCacheValidStruct(fileName string) (RssResponse, error) {
	return RssResponse{
		Feed: Feed{
			Entry:   []Entry{testEntry},
			Updated: Label{Label: "2025-06-26T10:06:25-07:00"},
		},
	}, nil
}

func mockSaveToCacheError(fileName string, data []byte) error {
	return errors.New("error in saveToCache")
}

func mockSaveToCacheSuccess(fileName string, data []byte) error {
	return nil
}
