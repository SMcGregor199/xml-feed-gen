package main

import (
	"fmt"
	"log"

	"shaynemcgregor.dev/xml-feed-gen/feed"
)

func main() {
	data, err := feed.FetchBlogDataBytes()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data))
}
