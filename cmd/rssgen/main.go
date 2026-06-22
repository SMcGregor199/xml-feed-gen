package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"shaynemcgregor.dev/xml-feed-gen/feed"
)

func main() {
	data, err := fetchBlogData()
	if err != nil {
		log.Fatal(err)
	}

	var posts []feed.Post
	if err := json.Unmarshal(data, &posts); err != nil {
		log.Fatal(err)
	}

	outputPath := "rss.xml"
	if configured := os.Getenv("RSS_OUTPUT_PATH"); configured != "" {
		outputPath = configured
	}

	config := feed.DefaultConfig()
	config.BuildTime = time.Now().UTC()
	if configured := os.Getenv("RSS_SITE_BASE_URL"); configured != "" {
		config.SiteBaseURL = configured
	}
	if configured := os.Getenv("RSS_PUBLIC_URL"); configured != "" {
		config.FeedURL = configured
	}

	if err := feed.WriteRSSXMLFile(outputPath, posts, config); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("RSS generated at %s\n", outputPath)
}

func fetchBlogData() ([]byte, error) {
	return feed.FetchBlogDataBytes()
}
