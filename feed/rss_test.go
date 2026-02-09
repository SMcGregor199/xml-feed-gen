package feed

import (
	"os"
	"strings"
	"testing"
)

func TestWriteRSSXML_Smoke(t *testing.T) {
	posts := []Post{
		{
			ID:            "2e40721d-8085-8043-b11d-fee776963643",
			Tag:           "Story Time",
			Title:         "When Speed Became the Right Tradeoff",
			Summary:       "On making a fast decision under pressure and living with it responsibly",
			Link:          "when-speed-became-the-right-tradeoff",
			Thumbnail:     "https://shaynemcgregordev-be.netlify.app/.netlify/functions/notion-image?blockId=abc",
			PublishedDate: "2026-01-10T13:58:00.000Z",
			UpdatedDate:   "2026-01-10T13:58:00.000Z",
			Body: []struct {
				Heading string   `json:"heading"`
				Paras   []string `json:"paras"`
			}{},
		},
	}

	tempDir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get wd: %v", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() {
		_ = os.Chdir(originalDir)
	}()

	if err := WriteRSSXML(posts); err != nil {
		t.Fatalf("WriteRSSXML: %v", err)
	}

	data, err := os.ReadFile("rss.xml")
	if err != nil {
		t.Fatalf("read rss.xml: %v", err)
	}

	output := string(data)
	required := []string{
		`<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">`,
		`<channel>`,
		`<title>`,
		`<link>`,
		`<description>`,
		`<lastBuildDate>`,
		`<language>en-us</language>`,
		`<ttl>60</ttl>`,
		`<item>`,
	}

	for _, snippet := range required {
		if !strings.Contains(output, snippet) {
			t.Fatalf("missing %q in output", snippet)
		}
	}

	assertSelfClosingElement(t, output, "atom:link", []string{
		`href="https://shaynemcgregor.dev/rss.xml"`,
		`rel="self"`,
		`type="application/rss+xml"`,
	})
	assertSelfClosingElement(t, output, "enclosure", []string{
		`url="https://shaynemcgregordev-be.netlify.app/.netlify/functions/notion-image?blockId=abc"`,
		`type="image/webp"`,
	})
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
