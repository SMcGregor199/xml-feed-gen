package feed

import (
	"bytes"
	"encoding/xml"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	siteBaseURL               = "https://shaynemcgregor.dev"
	feedURL                   = "https://shaynemcgregor.dev/rss.xml"
	defaultChannelTitle       = "shaynemcgregor.dev"
	defaultChannelDescription = "Posts from shaynemcgregor.dev"
)

type Post struct {
	ID            string `json:"id"`
	Tag           string `json:"tag"`
	Title         string `json:"title"`
	Summary       string `json:"summary"`
	Link          string `json:"link"`
	Thumbnail     string `json:"thumbnail"`
	PublishedDate string `json:"publishedDate"`
	UpdatedDate   string `json:"updatedDate"`
	Body          []struct {
		Heading string   `json:"heading"`
		Paras   []string `json:"paras"`
	} `json:"body"`
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	AtomNS  string   `xml:"xmlns:atom,attr,omitempty"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title         string   `xml:"title"`
	Link          string   `xml:"link"`
	Description   string   `xml:"description"`
	AtomLink      AtomLink `xml:"atom:link"`
	LastBuildDate string   `xml:"lastBuildDate"`
	Language      string   `xml:"language"`
	TTL           int      `xml:"ttl"`
	Items         []Item   `xml:"item"`
}

type Item struct {
	Title       string     `xml:"title"`
	Link        string     `xml:"link"`
	Description string     `xml:"description"`
	GUID        GUID       `xml:"guid"`
	PubDate     string     `xml:"pubDate"`
	Enclosure   *Enclosure `xml:"enclosure,omitempty"`
}

type GUID struct {
	IsPermaLink string `xml:"isPermaLink,attr"`
	Value       string `xml:",chardata"`
}

type Enclosure struct {
	URL  string `xml:"url,attr"`
	Type string `xml:"type,attr"`
}

type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

func (a AtomLink) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = []xml.Attr{
		{Name: xml.Name{Local: "href"}, Value: a.Href},
		{Name: xml.Name{Local: "rel"}, Value: a.Rel},
		{Name: xml.Name{Local: "type"}, Value: a.Type},
	}

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	return e.EncodeToken(start.End())
}

func (ecl Enclosure) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = []xml.Attr{
		{Name: xml.Name{Local: "url"}, Value: ecl.URL},
		{Name: xml.Name{Local: "type"}, Value: ecl.Type},
	}

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	return e.EncodeToken(start.End())
}

func WriteRSSXML(posts []Post) error {
	rssXML, err := buildRSSXML(posts)
	if err != nil {
		return err
	}

	return os.WriteFile("rss.xml", rssXML, 0644)
}

func buildRSSXML(posts []Post) ([]byte, error) {
	orderedPosts := orderPostsByDate(posts)
	items := make([]Item, 0, len(orderedPosts))

	for _, post := range orderedPosts {
		item := Item{
			Title:       post.Title,
			Link:        absoluteURL(siteBaseURL, post.Link),
			Description: post.Summary,
			GUID: GUID{
				IsPermaLink: "false",
				Value:       post.ID,
			},
			PubDate: rssDate(post.PublishedDate),
		}

		if strings.TrimSpace(post.Thumbnail) != "" {
			item.Enclosure = &Enclosure{
				URL:  post.Thumbnail,
				Type: "image/webp",
			}
		}

		items = append(items, item)
	}

	rss := RSS{
		Version: "2.0",
		AtomNS:  "http://www.w3.org/2005/Atom",
		Channel: Channel{
			Title:       defaultChannelTitle,
			Link:        siteBaseURL,
			Description: defaultChannelDescription,
			AtomLink: AtomLink{
				Href: feedURL,
				Rel:  "self",
				Type: "application/rss+xml",
			},
			LastBuildDate: time.Now().UTC().Format(time.RFC1123Z),
			Language:      "en-us",
			TTL:           60,
			Items:         items,
		},
	}

	output, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		return nil, err
	}

	withHeader := append([]byte(xml.Header), output...)
	return selfCloseEmptyElements(withHeader, []string{"atom:link", "enclosure"}), nil
}

func absoluteURL(base, slugOrURL string) string {
	trimmed := strings.TrimSpace(slugOrURL)
	if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") {
		return trimmed
	}
	if trimmed == "" {
		return base
	}
	return strings.TrimRight(base, "/") + "/" + strings.TrimLeft(trimmed, "/")
}

func rssDate(iso string) string {
	parsed, err := time.Parse(time.RFC3339, iso)
	if err != nil {
		return time.Now().UTC().Format(time.RFC1123Z)
	}

	return parsed.UTC().Format(time.RFC1123Z)
}

func orderPostsByDate(posts []Post) []Post {
	if len(posts) < 2 {
		return append([]Post(nil), posts...)
	}

	type postWithTime struct {
		post Post
		t    time.Time
	}

	parsed := make([]postWithTime, len(posts))
	for i, post := range posts {
		parsedTime, err := time.Parse(time.RFC3339, post.PublishedDate)
		if err != nil {
			return append([]Post(nil), posts...)
		}
		parsed[i] = postWithTime{post: post, t: parsedTime}
	}

	sort.SliceStable(parsed, func(i, j int) bool {
		return parsed[i].t.After(parsed[j].t)
	})

	ordered := make([]Post, len(parsed))
	for i, item := range parsed {
		ordered[i] = item.post
	}

	return ordered
}

func selfCloseEmptyElements(input []byte, tags []string) []byte {
	output := input
	for _, tag := range tags {
		closeForm := []byte("></" + tag + ">")
		shortForm := []byte(" />")
		output = bytes.ReplaceAll(output, closeForm, shortForm)
	}

	return output
}
