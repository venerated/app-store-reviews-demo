package main

import (
	"errors"
	"net/http"
	"reflect"
	"testing"
)

func TestGetCachedReviews(t *testing.T) {
	t.Run("err == nil and returns expected data", func(t *testing.T) {
		data, err := getCachedReviews(testCacheFilePath, mockReadFromCacheValidStruct)

		if err != nil {
			t.Errorf("err not nil %v", err)
		}

		expectedOutput := RssResponse{
			Feed: Feed{
				Entry: []Entry{testEntry},
				Updated: Label{
					Label: "2025-06-26T10:06:25-07:00",
				},
			},
		}

		got := data
		want := expectedOutput

		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v, got %v", want, got)
		}
	})

	t.Run("returns an error when file doesn't exist or is unreadable", func(t *testing.T) {
		_, err := getCachedReviews(testCacheFilePath, mockReadFromCacheFileDoesntExistOrUnreadable)

		got := err.Error()
		want := "error in getCachedReviews while getting cache file: file doesn't exist or is unreadable"

		if got != want {
			t.Errorf("expected %v, got %v", want, got)
		}
	})
}

func TestGetUpdatedData(t *testing.T) {
	t.Run("err == nil and returns expected data", func(t *testing.T) {
		resp := mockHTTPResponse(testRssJson, 200)
		data, err := getUpdatedData(testCacheFilePath, resp, mockSaveToCacheSuccess)

		if err != nil {
			t.Errorf("err not nil %v", err)
		}

		got := data
		want := testRssStructData

		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v, got %v", want, got)
		}
	})

	t.Run("err == nil and returns expected data when saving cache fails", func(t *testing.T) {
		resp := mockHTTPResponse(testRssJson, 200)
		data, err := getUpdatedData(testCacheFilePath, resp, mockSaveToCacheError)

		if err != nil {
			t.Errorf("err not nil %v", err)
		}

		got := data
		want := testRssStructData

		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v, got %v", want, got)
		}
	})

	t.Run("returns an error when data is malformed", func(t *testing.T) {
		resp := mockHTTPResponse("???", 200)
		_, err := getUpdatedData(testCacheFilePath, resp, mockSaveToCacheSuccess)

		got := err.Error()
		want := "error in getUpdatedData while decoding JSON: invalid character '?' looking for beginning of value"

		if got != want {
			t.Errorf("expected %v, got %v", want, got)
		}
	})
}

func TestGetAppStoreRss(t *testing.T) {
	testGetAppStoreRssDeps := GetAppStoreRssDeps{
		CacheFilePath:        testCacheFilePath,
		ClientFunc:           mockCustomClient,
		GetCachedReviewsFunc: mockGetCachedReviewsSuccess,
		GetUpdatedDataFunc:   mockGetUpdatedData,
		ReadFromCacheFunc:    mockReadFromCacheValidStruct,
		SaveToCacheFunc:      mockSaveToCacheSuccess,
	}

	t.Run("returns an error when error from external API", func(t *testing.T) {
		deps := testGetAppStoreRssDeps
		deps.ClientFunc = func(url string, customAddedHeaders CustomHeaders) (*http.Response, error) {
			return mockHTTPResponse("", 500), nil
		}

		_, err := getAppStoreRss("", deps)

		got := err.Error()
		want := "external API returned status 500"

		if got != want {
			t.Errorf("expected %v, got %v", want, got)
		}
	})

	t.Run("fetches fresh RSS data with no cache", func(t *testing.T) {
		data, err := getAppStoreRss("", testGetAppStoreRssDeps)

		if err != nil {
			t.Errorf("err not nil %v", err)
		}

		got := data
		want := testRssStructData

		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v, got %v", want, got)
		}
	})

	t.Run("uses cache timestamp and returns 304 response", func(t *testing.T) {
		deps := testGetAppStoreRssDeps
		deps.ClientFunc = func(url string, customAddedHeaders CustomHeaders) (*http.Response, error) {
			resp := mockHTTPResponse(testRssJson, 304)
			return resp, nil
		}

		data, err := getAppStoreRss("2025-06-26T10:06:25-07:00", deps)

		if err != nil {
			t.Errorf("err not nil %v", err)
		}

		got := data
		want := testRssStructData

		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v, got %v", want, got)
		}
	})
}

func TestGetData(t *testing.T) {
	testGetDataDeps := GetDataDeps{
		AppId:                testAppId,
		CacheFilePath:        testCacheFilePath,
		GetAppStoreRssFunc:   mockGetAppStoreRssSuccess,
		ClientFunc:           mockCustomClient,
		GetCachedReviewsFunc: mockGetCachedReviewsSuccess,
		GetUpdatedDataFunc:   mockGetUpdatedData,
		ReadFromCacheFunc:    mockReadFromCacheFileDoesntExistOrUnreadable,
		SaveToCacheFunc:      mockSaveToCacheSuccess,
	}

	t.Run("err == nil and returns expected data", func(t *testing.T) {
		data, err := getData(testGetDataDeps)

		if err != nil {
			t.Errorf("err not nil %v", err)
		}

		got := data
		want := testRssStructData

		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v, got %v", want, got)
		}
	})

	t.Run("falls back to fresh fetch if cache is unavailable", func(t *testing.T) {
		deps := testGetDataDeps
		deps.GetCachedReviewsFunc = mockGetCachedReviewsError

		data, err := getData(deps)

		if err != nil {
			t.Errorf("err not nil %v", err)
		}

		got := data
		want := testRssStructData

		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v, got %v", want, got)
		}
	})

	t.Run("errors if both cache and fetch fail", func(t *testing.T) {
		deps := testGetDataDeps
		deps.GetCachedReviewsFunc = mockGetCachedReviewsError
		deps.GetAppStoreRssFunc = func(lastUpdated string, deps GetAppStoreRssDeps) (RssResponse, error) {
			return RssResponse{}, errors.New("no data available")
		}

		_, err := getData(deps)

		got := err.Error()
		want := "no data available"

		if got != want {
			t.Errorf("expected %v, got %v", want, got)
		}
	})
}

func TestFormatReviews(t *testing.T) {

	t.Run("returns formatted review when given properly structured data", func(t *testing.T) {
		tests := []Entry{testEntry}

		expectedOutput := []Review{testReview}

		got := formatReviews(tests)
		want := expectedOutput

		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v, got %v", want, got)
		}
	})

	t.Run("returns partial formatted review when given partial data", func(t *testing.T) {
		tests := []Entry{
			{
				ID: Label{
					Label: "12802724742",
				},
				Author: Author{
					Label{
						Label: "ZenZen2025",
					},
				},
				Rating: Label{
					Label: "5",
				},
				Updated: Label{
					Label: "2025-06-21T21:09:29-07:00",
				},
			},
		}

		expectedOutput := []Review{
			{
				ID:      "12802724742",
				Author:  "ZenZen2025",
				Rating:  "5",
				Updated: "2025-06-21T21:09:29-07:00",
			},
		}

		got := formatReviews(tests)
		want := expectedOutput

		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v, got %v", want, got)
		}
	})

	t.Run("returns zero-valued review when given zero-valued entry", func(t *testing.T) {
		tests := []Entry{
			{},
		}

		expectedOutput := []Review{
			{},
		}

		got := formatReviews(tests)
		want := expectedOutput

		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v, got %v", want, got)
		}
	})

	t.Run("returns empty when given empty slice", func(t *testing.T) {
		got := formatReviews([]Entry{})

		if len(got) != 0 {
			t.Errorf("expected empty slice, got %v", got)
		}
	})
}
