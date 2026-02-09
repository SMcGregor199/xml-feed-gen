package main

import (
	"encoding/json"
	"fmt"
	"log"

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

	if err := feed.WriteRSSXML(posts); err != nil {
		log.Fatal(err)
	}

	fmt.Println("RSS generated at rss.xml")
}

func fetchBlogData() ([]byte, error) {
	return feed.FetchBlogDataBytes()
}
