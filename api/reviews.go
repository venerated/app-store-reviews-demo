package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Label struct {
	Label string `json:"label"`
}

type Author struct {
	Name Label `json:"name"`
}

type Entry struct {
	ID      Label  `json:"id"`
	Author  Author `json:"author"`
	Content Label  `json:"content"`
	Rating  Label  `json:"im:rating"`
	Updated Label  `json:"updated"`
}

type Feed struct {
	Entry   []Entry `json:"entry"`
	Updated Label   `json:"updated"`
}

type RssResponse struct {
	Feed Feed `json:"feed"`
}

type GetCachedReviews func(filePath string, readFromCacheFunc ReadRssFromCache) (RssResponse, error)

// getCachedReviews reads the cache file and unmarshals it into an RssResponse.
// Returns an error if the file can't be read or decoded.
func getCachedReviews(filePath string, readFromCacheFunc ReadRssFromCache) (RssResponse, error) {
	cacheFile, cacheFileErr := readFromCacheFunc(filePath)
	if cacheFileErr != nil {
		return RssResponse{}, fmt.Errorf("error in getCachedReviews while getting cache file: %w", cacheFileErr)
	}

	return cacheFile, nil
}

type GetUpdatedData func(cacheFilePath string, resp *http.Response, saveToCacheFunc SaveToCache) (RssResponse, error)

// getUpdatedData reads and unmarshals the HTTP response into an RssResponse.
// Also writes the raw response body to the cache file, logging cache errors silently.
func getUpdatedData(cacheFilePath string, resp *http.Response, saveToCacheFunc SaveToCache) (RssResponse, error) {
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return RssResponse{}, fmt.Errorf("error in getUpdatedData while reading response body: %w", err)
	}

	// Format data
	var data RssResponse
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return RssResponse{}, fmt.Errorf("error in getUpdatedData while decoding JSON: %w", err)
	}

	// Cache raw data
	if err := saveToCacheFunc(cacheFilePath, bodyBytes); err != nil {
		fmt.Printf("non-critical error in getUpdatedData while saving cache file: %v", err)
	}

	return data, nil
}

type GetAppStoreRssDeps struct {
	AppId                string
	CacheFilePath        string
	ClientFunc           CustomClient
	GetCachedReviewsFunc GetCachedReviews
	GetUpdatedDataFunc   GetUpdatedData
	ReadFromCacheFunc    ReadRssFromCache
	SaveToCacheFunc      SaveToCache
}

type GetAppStoreRss func(lastUpdated string, deps GetAppStoreRssDeps) (RssResponse, error)

// getAppStoreRss fetches the app store RSS feed, using conditional headers
// if a cache file exists/is readable.
// Uses cache file if data hasn't changed.
// Fetches new data if cache file doesn't exist/is unreadable.
func getAppStoreRss(lastUpdated string, deps GetAppStoreRssDeps) (RssResponse, error) {
	var rssUrl = "https://itunes.apple.com/us/rss/customerreviews/id=" + deps.AppId + "/sortBy=mostRecent/page=1/json"

	var customHeaders CustomHeaders = CustomHeaders{}

	// Send lastUpdated so we use cache file if data hasn't changed
	if lastUpdated != "" {
		customHeaders["If-Modified-Since"] = lastUpdated
	}

	// Create custom client so headers can be sent
	customClientResp, customClientErr := deps.ClientFunc(rssUrl, customHeaders)
	if customClientErr != nil {
		return RssResponse{}, fmt.Errorf("error in getAppStoreRss while using customClient: %w", customClientErr)
	}

	if customClientResp.StatusCode == 304 {
		// No change, use cache
		cachedData, cachedDataErr := deps.GetCachedReviewsFunc(deps.CacheFilePath, deps.ReadFromCacheFunc)
		if cachedDataErr != nil {
			// Cached data not available, request fresh data
			return deps.GetUpdatedDataFunc(deps.CacheFilePath, customClientResp, deps.SaveToCacheFunc)
		}
		return cachedData, nil
	} else {
		return deps.GetUpdatedDataFunc(deps.CacheFilePath, customClientResp, deps.SaveToCacheFunc)
	}
}

type GetDataDeps struct {
	AppId                string
	CacheFilePath        string
	GetAppStoreRssFunc   GetAppStoreRss
	ClientFunc           CustomClient
	GetCachedReviewsFunc GetCachedReviews
	GetUpdatedDataFunc   GetUpdatedData
	ReadFromCacheFunc    ReadRssFromCache
	SaveToCacheFunc      SaveToCache
}

// getData coordinates retrieval of RSS data.
func getData(deps GetDataDeps) (RssResponse, error) {

	cacheData, err := deps.ReadFromCacheFunc(deps.CacheFilePath)

	getAppStoreRssDeps := GetAppStoreRssDeps{
		AppId:                deps.AppId,
		CacheFilePath:        deps.CacheFilePath,
		ClientFunc:           deps.ClientFunc,
		GetCachedReviewsFunc: deps.GetCachedReviewsFunc,
		GetUpdatedDataFunc:   deps.GetUpdatedDataFunc,
		ReadFromCacheFunc:    deps.ReadFromCacheFunc,
		SaveToCacheFunc:      deps.SaveToCacheFunc,
	}

	// Cache file missing or unreadable, get fresh data
	if err != nil {
		return deps.GetAppStoreRssFunc("", getAppStoreRssDeps)
	}

	// Cache is usable, read the `updated` field from the file
	lastUpdated := convertTimestamp(cacheData.Feed.Updated.Label, time.RFC3339, http.TimeFormat)

	return deps.GetAppStoreRssFunc(lastUpdated, getAppStoreRssDeps)
}

type Review struct {
	ID      string `json:"id"`
	Author  string `json:"author"`
	Content string `json:"content"`
	Rating  string `json:"rating"`
	Updated string `json:"updated"`
}

// formatReviews transforms the raw RSS entries into a flattened, client-friendly format.
func formatReviews(data []Entry) []Review {
	var reviews []Review
	for i := range data {
		review := Review{
			ID:      data[i].ID.Label,
			Author:  data[i].Author.Name.Label,
			Content: data[i].Content.Label,
			Rating:  data[i].Rating.Label,
			Updated: data[i].Updated.Label,
		}
		reviews = append(reviews, review)
	}

	return reviews
}

// handleReviewsRequest is the Gin handler for the /reviews endpoint.
// It fetches RSS data, formats the reviews, and returns them as JSON.
func handleReviewsRequest(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	var cacheFilePath string = "cache/reviews_" + id + ".json"

	getDataDeps := GetDataDeps{
		AppId:                id,
		CacheFilePath:        cacheFilePath,
		GetAppStoreRssFunc:   getAppStoreRss,
		ClientFunc:           customHttpClient,
		GetCachedReviewsFunc: getCachedReviews,
		GetUpdatedDataFunc:   getUpdatedData,
		ReadFromCacheFunc:    readRssFromCache,
		SaveToCacheFunc:      saveToCache,
	}

	rss, rssErr := getData(getDataDeps)

	if rssErr != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": rssErr})
		return
	}

	reviews := formatReviews(rss.Feed.Entry)

	c.IndentedJSON(http.StatusOK, reviews)
}
