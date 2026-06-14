package core

/*
#cgo LDFLAGS: -L../.build/release -lMarkdownEngine -Wl,-rpath,./
#include <stdlib.h>

extern char* render_markdown_to_html(const char* input);
*/
import "C"
import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unsafe"
)

// MarkdownRenderer wraps the Swift-based markdown engine.
type MarkdownRenderer struct {
}

// NewMarkdownRenderer creates a new MarkdownRenderer.
func NewMarkdownRenderer() *MarkdownRenderer {
	return &MarkdownRenderer{}
}

// Render converts markdown to HTML using swift-markdown via FFI.
func (r *MarkdownRenderer) Render(content string) (string, string, error) {
	cContent := C.CString(content)
	defer C.free(unsafe.Pointer(cContent))

	cResult := C.render_markdown_to_html(cContent)
	if cResult == nil {
		return "", "", errors.New("swift-markdown rendering failed")
	}
	defer C.free(unsafe.Pointer(cResult))

	html := C.GoString(cResult)
	html, tocJSON := processHeadings(html)
	return html, tocJSON, nil
}

// processHeadings adds IDs to headings and generates a TOC JSON.
func processHeadings(html string) (string, string) {
	// Regex patterns for each heading level
	hPatterns := []*regexp.Regexp{
		regexp.MustCompile(`<h1>([^<]*)</h1>`),
		regexp.MustCompile(`<h2>([^<]*)</h2>`),
		regexp.MustCompile(`<h3>([^<]*)</h3>`),
		regexp.MustCompile(`<h4>([^<]*)</h4>`),
		regexp.MustCompile(`<h5>([^<]*)</h5>`),
		regexp.MustCompile(`<h6>([^<]*)</h6>`),
	}

	// Collect all matches with their levels
	type headingMatch struct{ level int; text string }
	allMatches := []headingMatch{}

	for level, pattern := range hPatterns {
		matches := pattern.FindAllStringSubmatch(html, -1)
		for _, m := range matches {
			allMatches = append(allMatches, headingMatch{level + 1, m[1]})
		}
	}

	// Track heading IDs to ensure uniqueness
	idCounts := make(map[string]int)

	// Build TOC and new HTML
	type tocEntry struct {
		Level int    `json:"level"`
		ID    string `json:"id"`
		Text  string `json:"text"`
	}
	tocEntries := []tocEntry{}
	result := html

	// Process in reverse order to preserve positions
	for i := 0; i < len(allMatches); i++ {
		level := allMatches[i].level
		text := allMatches[i].text

		// Generate unique ID from heading text
		idBase := slugify(text)
		idCounts[idBase]++
		id := idBase
		if idCounts[idBase] > 1 {
			id = idBase + "-" + fmt.Sprint(idCounts[idBase]-1)
		}

		// Add to TOC
		tocEntries = append(tocEntries, tocEntry{level, id, text})

		// Replace in HTML (surgical replace)
		old := fmt.Sprintf("<h%d>%s</h%d>", level, text, level)
		new := fmt.Sprintf("<h%d id=\"%s\">%s</h%d>", level, id, text, level)
		result = strings.Replace(result, old, new, 1)
	}

	// Build TOC JSON
	tocJSONBytes, _ := json.Marshal(tocEntries)
	tocJSON := string(tocJSONBytes)

	return result, tocJSON
}

// slugify converts text to URL-friendly ID
func slugify(text string) string {
	re := regexp.MustCompile("[^a-zA-Z0-9\u4e00-\u9fff]+")
	slug := re.ReplaceAllString(text, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = "heading"
	}
	return strings.ToLower(slug)
}
