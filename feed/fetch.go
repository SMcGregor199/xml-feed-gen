package feed

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const blogDataEndPoint = "https://shaynemcgregordev-be.netlify.app/.netlify/functions/blog-posts-json"

func FetchBlogDataBytes() ([]byte, error) {
	return FetchBlogDataBytesFromURL(blogDataEndPoint, http.DefaultClient)
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func FetchBlogDataBytesFromURL(endpoint string, client httpClient) ([]byte, error) {
	if strings.TrimSpace(endpoint) == "" {
		return nil, fmt.Errorf("blog data endpoint is required")
	}
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create blog data request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch blog data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil, fmt.Errorf("fetch blog data: unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read blog data response: %w", err)
	}

	return body, nil
}
