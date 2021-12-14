// +build integration

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/quantonganh/shortener"
)

const (
	longURL = "https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies"
)

func TestShortener(t *testing.T) {
	var shortURL string
	t.Run("create short url", func(t *testing.T) {
		shortURL = testCreateShortURL(t)
	})

	t.Run("redirect", func(t *testing.T) {
		testRedirect(t, shortURL)
	})
}

func testCreateShortURL(t *testing.T) string {
	r := shortener.URLCreationRequest{
		LongURL: longURL,
	}
	body, err := json.Marshal(r)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/create-short-url", bytes.NewBuffer(body))
	require.NoError(t, err)

	client := http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	return response["short_url"]
}

func testRedirect(t *testing.T, shortURL string) {
	req, err := http.NewRequest(http.MethodGet, shortURL, nil)
	require.NoError(t, err)

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, longURL, resp.Header.Get("Location"))
}
