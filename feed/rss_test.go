package feed

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGenerateRSSXML_ProducesValidDeterministicRSS(t *testing.T) {
	output := generateRSS(t, []Post{
		post("post_b", "Second & Later", "second-later", "2026-01-11T13:58:00.000Z"),
		post("post_a", "First <Earlier>", "first-earlier", "2026-01-10T13:58:00.000Z"),
	})

	var rss RSS
	if err := xml.Unmarshal(output, &rss); err != nil {
		t.Fatalf("rss should be valid XML: %v\n%s", err, output)
	}

	text := string(output)
	assertContains(t, text, `<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">`)
	assertContains(t, text, `<atom:link href="https://shaynemcgregor.dev/rss.xml" rel="self" type="application/rss+xml" />`)
	assertContains(t, text, `<lastBuildDate>Sun, 15 Feb 2026 12:00:00 +0000</lastBuildDate>`)
	assertContains(t, text, `<title>Second &amp; Later</title>`)
	assertContains(t, text, `<title>First &lt;Earlier&gt;</title>`)
	assertOrder(t, text, "post_b", "post_a")
}

func TestGenerateRSSXML_UsesBlogCanonicalLinksAndStableGUIDs(t *testing.T) {
	output := string(generateRSS(t, []Post{
		post("stable_id", "Canonical", "canonical-post", "2026-01-10T13:58:00.000Z"),
	}))

	assertContains(t, output, `<link>https://shaynemcgregor.dev/blog/canonical-post</link>`)
	assertContains(t, output, `<guid isPermaLink="false">stable_id</guid>`)
}

func TestGenerateRSSXML_SummaryBodyFallbackAndEnclosure(t *testing.T) {
	p := post("post_body", "Body Fallback", "body-fallback", "2026-01-10T13:58:00.000Z")
	p.Summary = ""
	p.Thumbnail = "https://example.test/image.png"
	p.Body = []struct {
		Heading string   `json:"heading"`
		Paras   []string `json:"paras"`
	}{
		{Heading: "Intro", Paras: []string{"", "Body paragraph fallback."}},
	}

	output := string(generateRSS(t, []Post{p}))

	assertContains(t, output, `<description>Body paragraph fallback.</description>`)
	assertSelfClosingElement(t, output, "enclosure", []string{
		`url="https://example.test/image.png"`,
		`type="image/png"`,
	})
}

func TestGenerateRSSXML_SkipsInvalidThumbnail(t *testing.T) {
	p := post("post_no_enclosure", "No Enclosure", "no-enclosure", "2026-01-10T13:58:00.000Z")
	p.Thumbnail = "/relative-image.png"

	output := string(generateRSS(t, []Post{p}))

	if strings.Contains(output, "<enclosure") {
		t.Fatalf("unexpected enclosure for relative thumbnail:\n%s", output)
	}
}

func TestGenerateRSSXML_RejectsMalformedPostsWithoutNowFallback(t *testing.T) {
	tests := []struct {
		name string
		post Post
	}{
		{name: "missing id", post: Post{Title: "Missing ID", Link: "missing-id", PublishedDate: "2026-01-10T13:58:00.000Z"}},
		{name: "missing title", post: Post{ID: "post_1", Link: "missing-title", PublishedDate: "2026-01-10T13:58:00.000Z"}},
		{name: "missing link", post: Post{ID: "post_1", Title: "Missing Link", PublishedDate: "2026-01-10T13:58:00.000Z"}},
		{name: "invalid date", post: Post{ID: "post_1", Title: "Bad Date", Link: "bad-date", PublishedDate: "not-a-date"}},
		{name: "missing date", post: Post{ID: "post_1", Title: "No Date", Link: "no-date"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := GenerateRSSXML([]Post{tt.post}, testConfig()); err == nil {
				t.Fatalf("expected malformed post error")
			}
		})
	}
}

func TestGenerateRSSXML_AllowsEmptyFeed(t *testing.T) {
	output := string(generateRSS(t, nil))

	assertContains(t, output, `<channel>`)
	if strings.Contains(output, "<item>") {
		t.Fatalf("empty feed should not include items:\n%s", output)
	}
}

func TestGenerateRSSXML_DeterministicTieOrdering(t *testing.T) {
	output := string(generateRSS(t, []Post{
		post("post_c", "C", "c", "2026-01-10T13:58:00.000Z"),
		post("post_a", "A", "a", "2026-01-10T13:58:00.000Z"),
		post("post_b", "B", "b", "2026-01-10T13:58:00.000Z"),
	}))

	assertOrder(t, output, "post_a", "post_b")
	assertOrder(t, output, "post_b", "post_c")
}

func TestWriteRSSXMLFile_DoesNotOverwriteExistingFileOnValidationError(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "rss.xml")
	original := []byte("last known good rss")
	if err := os.WriteFile(path, original, 0644); err != nil {
		t.Fatalf("seed rss file: %v", err)
	}

	err := WriteRSSXMLFile(path, []Post{{ID: "bad", Title: "Bad", Link: "bad", PublishedDate: "not-a-date"}}, testConfig())
	if err == nil {
		t.Fatalf("expected validation error")
	}

	data, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("read rss file: %v", readErr)
	}
	if string(data) != string(original) {
		t.Fatalf("existing rss file was overwritten: %q", data)
	}
}

func TestFetchBlogDataBytesFromURL_ValidatesStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		_, _ = w.Write([]byte(`{"error":"short and stout"}`))
	}))
	defer server.Close()

	if _, err := FetchBlogDataBytesFromURL(server.URL, server.Client()); err == nil {
		t.Fatalf("expected non-2xx status error")
	}
}

func TestFetchBlogDataBytesFromURL_ReturnsBodyForSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Fatalf("unexpected Accept header: %q", got)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"id":"post_1"}]`))
	}))
	defer server.Close()

	body, err := FetchBlogDataBytesFromURL(server.URL, server.Client())
	if err != nil {
		t.Fatalf("FetchBlogDataBytesFromURL: %v", err)
	}
	if string(body) != `[{"id":"post_1"}]` {
		t.Fatalf("unexpected body: %s", body)
	}
}

func generateRSS(t *testing.T, posts []Post) []byte {
	t.Helper()

	output, err := GenerateRSSXML(posts, testConfig())
	if err != nil {
		t.Fatalf("GenerateRSSXML: %v", err)
	}
	return output
}

func testConfig() Config {
	return Config{
		SiteBaseURL:        "https://shaynemcgregor.dev",
		FeedURL:            "https://shaynemcgregor.dev/rss.xml",
		ChannelTitle:       "shaynemcgregor.dev",
		ChannelDescription: "Posts from shaynemcgregor.dev",
		Language:           "en-us",
		TTL:                60,
		BuildTime:          time.Date(2026, 2, 15, 12, 0, 0, 0, time.UTC),
	}
}

func post(id, title, link, publishedDate string) Post {
	return Post{
		ID:            id,
		Tag:           "Story Time",
		Title:         title,
		Summary:       "Summary for " + title,
		Link:          link,
		Thumbnail:     "https://shaynemcgregordev-be.netlify.app/.netlify/functions/notion-image?blockId=abc",
		PublishedDate: publishedDate,
		UpdatedDate:   publishedDate,
	}
}

func assertContains(t *testing.T, output, snippet string) {
	t.Helper()

	if !strings.Contains(output, snippet) {
		t.Fatalf("missing %q in output:\n%s", snippet, output)
	}
}

func assertOrder(t *testing.T, output, first, second string) {
	t.Helper()

	firstIndex := strings.Index(output, first)
	secondIndex := strings.Index(output, second)
	if firstIndex == -1 || secondIndex == -1 {
		t.Fatalf("could not find %q and %q in output:\n%s", first, second, output)
	}
	if firstIndex >= secondIndex {
		t.Fatalf("expected %q before %q in output:\n%s", first, second, output)
	}
}

func assertSelfClosingElement(t *testing.T, output, tag string, attrs []string) {
	t.Helper()

	start := strings.Index(output, "<"+tag)
	if start == -1 {
		t.Fatalf("missing <%s> element", tag)
	}

	end := strings.Index(output[start:], ">")
	if end == -1 {
		t.Fatalf("unterminated <%s> element", tag)
	}

	segment := output[start : start+end+1]
	if !strings.HasSuffix(strings.TrimSpace(segment), "/>") {
		t.Fatalf("<%s> is not self-closing: %q", tag, segment)
	}

	for _, attr := range attrs {
		if !strings.Contains(segment, attr) {
			t.Fatalf("<%s> missing %q", tag, attr)
		}
	}
}
