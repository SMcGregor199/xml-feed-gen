package feed

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	defaultSiteBaseURL           = "https://shaynemcgregor.dev"
	defaultFeedURL               = "https://shaynemcgregor.dev/rss.xml"
	defaultChannelTitle          = "shaynemcgregor.dev"
	defaultChannelDescription    = "Posts from shaynemcgregor.dev"
	defaultLanguage              = "en-us"
	defaultTTL                   = 60
	defaultRSSOutputFilePerm     = 0644
	defaultRSSOutputTempFilePerm = 0600
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

type Config struct {
	SiteBaseURL        string
	FeedURL            string
	ChannelTitle       string
	ChannelDescription string
	Language           string
	TTL                int
	BuildTime          time.Time
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
	Description string     `xml:"description,omitempty"`
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

func DefaultConfig() Config {
	return Config{
		SiteBaseURL:        defaultSiteBaseURL,
		FeedURL:            defaultFeedURL,
		ChannelTitle:       defaultChannelTitle,
		ChannelDescription: defaultChannelDescription,
		Language:           defaultLanguage,
		TTL:                defaultTTL,
		BuildTime:          time.Now().UTC(),
	}
}

func WriteRSSXML(posts []Post) error {
	return WriteRSSXMLFile("rss.xml", posts, DefaultConfig())
}

func WriteRSSXMLFile(path string, posts []Post, config Config) error {
	rssXML, err := GenerateRSSXML(posts, config)
	if err != nil {
		return err
	}

	return writeFileAtomically(path, rssXML)
}

func GenerateRSSXML(posts []Post, config Config) ([]byte, error) {
	config = normalizeConfig(config)
	orderedPosts, err := orderPostsByDate(posts)
	if err != nil {
		return nil, err
	}

	items := make([]Item, 0, len(orderedPosts))
	for _, post := range orderedPosts {
		item, err := itemFromPost(post, config)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	rss := RSS{
		Version: "2.0",
		AtomNS:  "http://www.w3.org/2005/Atom",
		Channel: Channel{
			Title:       config.ChannelTitle,
			Link:        config.SiteBaseURL,
			Description: config.ChannelDescription,
			AtomLink: AtomLink{
				Href: config.FeedURL,
				Rel:  "self",
				Type: "application/rss+xml",
			},
			LastBuildDate: config.BuildTime.UTC().Format(time.RFC1123Z),
			Language:      config.Language,
			TTL:           config.TTL,
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

func itemFromPost(post Post, config Config) (Item, error) {
	id := strings.TrimSpace(post.ID)
	if id == "" {
		return Item{}, fmt.Errorf("rss post is missing stable id")
	}

	title := strings.TrimSpace(post.Title)
	if title == "" {
		return Item{}, fmt.Errorf("rss post %q is missing title", id)
	}

	link, err := postURL(config.SiteBaseURL, post.Link)
	if err != nil {
		return Item{}, fmt.Errorf("rss post %q has invalid link: %w", id, err)
	}

	published, err := parseRequiredRSSDate(post.PublishedDate)
	if err != nil {
		return Item{}, fmt.Errorf("rss post %q has invalid publishedDate: %w", id, err)
	}

	item := Item{
		Title:       title,
		Link:        link,
		Description: descriptionFromPost(post),
		GUID: GUID{
			IsPermaLink: "false",
			Value:       id,
		},
		PubDate: published,
	}

	if enclosure := enclosureFromPost(post); enclosure != nil {
		item.Enclosure = enclosure
	}

	return item, nil
}

func normalizeConfig(config Config) Config {
	defaults := DefaultConfig()

	if strings.TrimSpace(config.SiteBaseURL) == "" {
		config.SiteBaseURL = defaults.SiteBaseURL
	}
	config.SiteBaseURL = strings.TrimRight(strings.TrimSpace(config.SiteBaseURL), "/")

	if strings.TrimSpace(config.FeedURL) == "" {
		config.FeedURL = defaults.FeedURL
	}
	if strings.TrimSpace(config.ChannelTitle) == "" {
		config.ChannelTitle = defaults.ChannelTitle
	}
	if strings.TrimSpace(config.ChannelDescription) == "" {
		config.ChannelDescription = defaults.ChannelDescription
	}
	if strings.TrimSpace(config.Language) == "" {
		config.Language = defaults.Language
	}
	if config.TTL <= 0 {
		config.TTL = defaults.TTL
	}
	if config.BuildTime.IsZero() {
		config.BuildTime = defaults.BuildTime
	}

	return config
}

func postURL(base, slugOrURL string) (string, error) {
	trimmed := strings.TrimSpace(slugOrURL)
	if trimmed == "" {
		return "", fmt.Errorf("empty slug")
	}

	if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") {
		if _, err := url.ParseRequestURI(trimmed); err != nil {
			return "", err
		}
		return trimmed, nil
	}

	trimmed = strings.TrimLeft(trimmed, "/")
	if !strings.HasPrefix(trimmed, "blog/") {
		trimmed = "blog/" + trimmed
	}

	return strings.TrimRight(base, "/") + "/" + trimmed, nil
}

func parseRequiredRSSDate(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", fmt.Errorf("date is required")
	}

	parsed, err := time.Parse(time.RFC3339, trimmed)
	if err != nil {
		return "", err
	}

	return parsed.UTC().Format(time.RFC1123Z), nil
}

func descriptionFromPost(post Post) string {
	if summary := strings.TrimSpace(post.Summary); summary != "" {
		return summary
	}

	for _, section := range post.Body {
		for _, para := range section.Paras {
			if text := strings.TrimSpace(para); text != "" {
				return text
			}
		}
	}

	return ""
}

func enclosureFromPost(post Post) *Enclosure {
	thumbnail := strings.TrimSpace(post.Thumbnail)
	if thumbnail == "" || !isHTTPURL(thumbnail) {
		return nil
	}

	return &Enclosure{
		URL:  thumbnail,
		Type: imageContentType(thumbnail),
	}
}

func isHTTPURL(value string) bool {
	parsed, err := url.ParseRequestURI(value)
	if err != nil {
		return false
	}
	return parsed.Scheme == "http" || parsed.Scheme == "https"
}

func imageContentType(imageURL string) string {
	parsed, err := url.Parse(imageURL)
	if err != nil {
		return "image/webp"
	}

	switch strings.ToLower(filepath.Ext(parsed.Path)) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	default:
		return "image/webp"
	}
}

func orderPostsByDate(posts []Post) ([]Post, error) {
	type postWithTime struct {
		post Post
		t    time.Time
	}

	parsed := make([]postWithTime, len(posts))
	for i, post := range posts {
		parsedTime, err := time.Parse(time.RFC3339, strings.TrimSpace(post.PublishedDate))
		if err != nil {
			return nil, fmt.Errorf("rss post %q has invalid publishedDate: %w", strings.TrimSpace(post.ID), err)
		}
		parsed[i] = postWithTime{post: post, t: parsedTime}
	}

	sort.SliceStable(parsed, func(i, j int) bool {
		if parsed[i].t.Equal(parsed[j].t) {
			return strings.TrimSpace(parsed[i].post.ID) < strings.TrimSpace(parsed[j].post.ID)
		}
		return parsed[i].t.After(parsed[j].t)
	})

	ordered := make([]Post, len(parsed))
	for i, item := range parsed {
		ordered[i] = item.post
	}

	return ordered, nil
}

func writeFileAtomically(path string, data []byte) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("rss output path is required")
	}

	dir := filepath.Dir(path)
	temp, err := os.CreateTemp(dir, ".rss-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp rss file: %w", err)
	}
	tempPath := temp.Name()
	defer func() {
		_ = os.Remove(tempPath)
	}()

	if _, err := temp.Write(data); err != nil {
		_ = temp.Close()
		return fmt.Errorf("write temp rss file: %w", err)
	}
	if err := temp.Chmod(defaultRSSOutputTempFilePerm); err != nil {
		_ = temp.Close()
		return fmt.Errorf("chmod temp rss file: %w", err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("close temp rss file: %w", err)
	}
	if err := os.Chmod(tempPath, defaultRSSOutputFilePerm); err != nil {
		return fmt.Errorf("chmod rss file: %w", err)
	}
	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("replace rss file: %w", err)
	}

	return nil
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
