package feed

import (
	"io"
	"net/http"
)

const blogDataEndPoint = "https://shaynemcgregordev-be.netlify.app/.netlify/functions/blog-posts-json"

func FetchBlogDataBytes() ([]byte, error) {
	resp, err := http.Get(blogDataEndPoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
