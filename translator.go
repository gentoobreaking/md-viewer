package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// =============================================================================
// Config
// =============================================================================

func getBackend() string {
	backend := os.Getenv("TRANSLATE_BACKEND")
	if backend == "" {
		return "mymemory"
	}
	return backend
}

// =============================================================================
// MyMemory (default, no API key needed)
// =============================================================================

type MyMemoryResponse struct {
	ResponseData struct {
		TranslatedText string `json:"translatedText"`
		QuotaFinished   bool   `json:"quotaFinished"`
	} `json:"responseData"`
	ResponseStatus any    `json:"responseStatus"`
}

// MyMemoryTranslate translates text using MyMemory API with retry logic for 429 errors
func MyMemoryTranslate(text, source, target string) (string, error) {
	if text == "" {
		return "", nil
	}

	if source != "auto" {
		source = normalizeChineseCode(source)
	}
	target = normalizeChineseCode(target)

	langPair := source + "|" + target
	if source == "auto" {
		source = "en"
		langPair = source + "|" + target
	}

	apiURL := fmt.Sprintf(
		"https://api.mymemory.translated.net/get?q=%s&langpair=%s",
		url.QueryEscape(text),
		url.QueryEscape(langPair),
	)

	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 1s, 2s...
			time.Sleep(time.Duration(attempt) * time.Second)
		}

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Get(apiURL)
		if err != nil {
			lastErr = err
			continue
		}
		
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode == 429 {
			lastErr = fmt.Errorf("rate limited (429)")
			debugLog("[MyMemory] Rate limited, attempt %d/3...", attempt+1)
			continue
		}

		var result MyMemoryResponse
		if err := json.Unmarshal(body, &result); err != nil {
			lastErr = err
			continue
		}

		// Check status
		statusOK := false
		switch v := result.ResponseStatus.(type) {
		case float64: statusOK = (int(v) == 200)
		case string: statusOK = (v == "200")
		}

		if !statusOK {
			lastErr = fmt.Errorf("API error: %v", result.ResponseStatus)
			continue
		}

		if result.ResponseData.QuotaFinished {
			return "", fmt.Errorf("MyMemory daily quota exceeded")
		}

		return result.ResponseData.TranslatedText, nil
	}

	return "", fmt.Errorf("MyMemory failed after retries: %v", lastErr)
}

// =============================================================================
// DeepL API (requires DEEPL_API_KEY)
// =============================================================================

type DeepLRequest struct {
	Text       []string `json:"text"`
	SourceLang string   `json:"source_lang"`
	TargetLang string   `json:"target_lang"`
}

type DeepLResponse struct {
	Translations []struct {
		Text string `json:"text"`
	} `json:"translations"`
}

// DeepLTranslate translates text using DeepL API
func DeepLTranslate(text, source, target string) (string, error) {
	if text == "" {
		return "", nil
	}

	apiKey := os.Getenv("DEEPL_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("DEEPL_API_KEY not set")
	}

	// DeepL uses different lang codes
	realSource := deepLLang(source)
	realTarget := deepLLang(target)

	apiURL := "https://api-free.deepl.com/v2/translate"
	if strings.Contains(apiKey, ":fx") {
		apiURL = "https://api-pro.deepl.com/v2/translate"
	}

	reqBody := DeepLRequest{
		Text:       []string{text},
		SourceLang: realSource,
		TargetLang: realTarget,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("DeepL marshal failed: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("DeepL request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "DeepL-Auth-Key "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("DeepL network failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("DeepL read failed: %w", err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("DeepL error %d: %s", resp.StatusCode, string(body))
	}

	var result DeepLResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("DeepL parse failed: %w", err)
	}

	if len(result.Translations) == 0 {
		return "", fmt.Errorf("DeepL no translations returned")
	}

	return result.Translations[0].Text, nil
}

func deepLLang(code string) string {
	switch normalizeChineseCode(code) {
	case "zh-TW":
		return "ZH-HANT"
	case "zh-CN":
		return "ZH-HANS"
	case "auto":
		return "auto"
	default:
		return strings.ToUpper(code)
	}
}

// =============================================================================
// LibreTranslate (requires LIBRETRANSLATE_API_KEY)
// =============================================================================

type LibreTranslateRequest struct {
	Q      string `json:"q"`
	Source string `json:"source"`
	Target string `json:"target"`
	Format string `json:"format"`
	APIKey string `json:"api_key,omitempty"`
}

type LibreTranslateResponse struct {
	TranslatedText string `json:"translatedText"`
}

// LibreTranslate translates text using LibreTranslate API
func LibreTranslate(text, source, target string) (string, error) {
	if text == "" {
		return "", nil
	}

	src, tgt := source, target
	if src != "auto" {
		src = normalizeChineseCode(src)
	}
	tgt = normalizeChineseCode(tgt)
	if src == "zh-CN" {
		src = "zh"
	}
	if tgt == "zh-CN" {
		tgt = "zh"
	}
	if src == "zh-TW" {
		src = "zt"
	}
	if tgt == "zh-TW" {
		tgt = "zt"
	}

	apiKey := os.Getenv("LIBRETRANSLATE_API_KEY")
	apiURL := os.Getenv("LIBRETRANSLATE_URL")
	if apiURL == "" {
		apiURL = "https://libretranslate.com/translate"
	}

	reqBody := LibreTranslateRequest{
		Q:      text,
		Source: src,
		Target: tgt,
		Format: "text",
		APIKey: apiKey,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("LibreTranslate marshal failed: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("LibreTranslate request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("LibreTranslate network failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("LibreTranslate read failed: %w", err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("LibreTranslate error %d: %s", resp.StatusCode, string(body))
	}

	var result LibreTranslateResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("LibreTranslate parse failed: %w", err)
	}

	return result.TranslatedText, nil
}

// =============================================================================
// Unified Translate interface
// =============================================================================

func debugLog(format string, a ...interface{}) {
	if currentConfig.DebugMode {
		fmt.Printf("[DEBUG] "+format+"\n", a...)
	}
}

// normalizeChineseCode maps UI / legacy codes to explicit script variants for translation APIs.
// Unqualified "zh" is treated as Simplified (zh-CN), matching common API defaults.
func normalizeChineseCode(code string) string {
	c := strings.ToLower(strings.TrimSpace(code))
	c = strings.ReplaceAll(c, "_", "-")
	switch c {
	case "zh-tw", "zh-hant", "zht":
		return "zh-TW"
	case "zh-cn", "zh-hans", "zh", "zhs":
		return "zh-CN"
	default:
		return code
	}
}

// GoogleTranslate translates text using Google Translate's gtx API (more reliable free tier)
func GoogleTranslate(text, source, target string) (string, error) {
	if text == "" {
		return "", nil
	}

	if source != "auto" {
		source = normalizeChineseCode(source)
	}
	target = normalizeChineseCode(target)

	apiURL := fmt.Sprintf(
		"https://translate.googleapis.com/translate_a/single?client=gtx&sl=%s&tl=%s&dt=t&q=%s",
		url.QueryEscape(source),
		url.QueryEscape(target),
		url.QueryEscape(text),
	)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Google error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Google returns a complex nested array: [[[ "translated", "source", ... ]]]
	var result []interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if len(result) > 0 {
		if inner, ok := result[0].([]interface{}); ok {
			translatedText := ""
			for _, item := range inner {
				if entry, ok := item.([]interface{}); ok && len(entry) > 0 {
					if t, ok := entry[0].(string); ok {
						translatedText += t
					}
				}
			}
			return translatedText, nil
		}
	}

	return "", fmt.Errorf("failed to parse Google response")
}

// Translate calls the configured backend to translate text
func Translate(text, source, target string) (string, error) {
	if text == "" {
		return "", nil
	}

	backend := getBackend()
	// Default to google if no backend set or if mymemory is failing
	if backend == "" || backend == "mymemory" {
		backend = "google"
	}

	debugLog("[Translate] backend=%s, text=%q, %s→%s", backend, text[:min(20, len(text))], source, target)

	switch backend {
	case "google":
		return GoogleTranslate(text, source, target)
	case "deepl":
		return DeepLTranslate(text, source, target)
	case "libre":
		return LibreTranslate(text, source, target)
	case "mymemory":
		return MyMemoryTranslate(text, source, target)
	default:
		return GoogleTranslate(text, source, target)
	}
}

// TranslateMarkdown translates Markdown content with batching
func TranslateMarkdown(md, source, target string, onProgress func(int)) (string, error) {
	if md == "" { return "", nil }
	reader := strings.NewReader(md)
	doc, err := html.Parse(reader)
	if err != nil { return md, err }

	var nodes []*html.Node
	var collect func(*html.Node)
	collect = func(n *html.Node) {
		if n.Type == html.TextNode {
			if strings.TrimSpace(n.Data) != "" && !isCodeBlock(n) {
				nodes = append(nodes, n)
			}
			return
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling { collect(child) }
	}
	collect(doc)

	if len(nodes) > 0 {
		currentBatch := []*html.Node{}
		currentLen := 0
		processedCount := 0
		var lastErr error
		
		processNodes := func(batch []*html.Node) error {
			if len(batch) == 0 { return nil }
			texts := make([]string, len(batch))
			for i, n := range batch { texts[i] = n.Data }
			
			joined := strings.Join(texts, "\n---\n")
			debugLog("[Translate] Batch HTML: %d nodes (%d chars)...", len(batch), len(joined))
			translated, err := Translate(joined, source, target)
			time.Sleep(800 * time.Millisecond)
			
			if err != nil {
				return err
			}
			parts := strings.Split(translated, "\n---\n")
			for i, n := range batch {
				if i < len(parts) { n.Data = strings.TrimSpace(parts[i]) }
			}
			
			processedCount += len(batch)
			if onProgress != nil {
				percent := (processedCount * 100) / len(nodes)
				onProgress(percent)
			}
			return nil
		}

		for _, n := range nodes {
			if currentLen + len(n.Data) > 800 {
				if err := processNodes(currentBatch); err != nil { lastErr = err }
				currentBatch = []*html.Node{}
				currentLen = 0
			}
			currentBatch = append(currentBatch, n)
			currentLen += len(n.Data)
		}
		if err := processNodes(currentBatch); err != nil { lastErr = err }
		
		if lastErr != nil && len(nodes) > 5 {
			return "", fmt.Errorf("翻譯失敗 (API 限制): %v", lastErr)
		}
	}

	var buf bytes.Buffer
	html.Render(&buf, doc)
	return buf.String(), nil
}

// TranslateMarkdownSource translates raw Markdown text using batching
func TranslateMarkdownSource(md, source, target string, onProgress func(int)) (string, error) {
	if md == "" {
		return "", nil
	}

	// 1. Protect code blocks
	codeBlocks := []string{}
	codeRegex := regexp.MustCompile("(?s)```.*?```|`.*?`")
	mdProtected := codeRegex.ReplaceAllStringFunc(md, func(match string) string {
		placeholder := fmt.Sprintf("[[CP_%d]]", len(codeBlocks))
		codeBlocks = append(codeBlocks, match)
		return placeholder
	})

	// 2. Protect Links and Images
	links := []string{}
	linkRegex := regexp.MustCompile(`(!?\[)(.*?)(\]\()(.*?)(\))`)
	mdProtected = linkRegex.ReplaceAllStringFunc(mdProtected, func(match string) string {
		placeholder := fmt.Sprintf("[[LP_%d]]", len(links))
		links = append(links, match)
		return placeholder
	})

	// 3. Batch translation
	lines := strings.Split(mdProtected, "\n")
	translatedLines := make([]string, len(lines))
	
	const maxBatchChars = 800
	currentBatch := []int{}
	currentBatchLen := 0
	processedLines := 0
	var lastErr error

	processBatch := func(indices []int) error {
		if len(indices) == 0 { return nil }
		
		textsToTranslate := []string{}
		prefixes := make([]string, len(indices))
		
		for i, idx := range indices {
			line := lines[idx]
			prefix := ""
			listRegex := regexp.MustCompile(`^(\s*([-*+]|\d+\.)\s+|>+\s*)+`)
			if loc := listRegex.FindStringIndex(line); loc != nil {
				prefix = line[loc[0]:loc[1]]
				line = line[loc[1]:]
			}
			prefixes[i] = prefix
			textsToTranslate = append(textsToTranslate, line)
		}

		joinedText := strings.Join(textsToTranslate, "\n---\n")
		debugLog("[Translate] Batch Source: %d lines (%d chars)...", len(indices), len(joinedText))
		translatedJoined, err := Translate(joinedText, source, target)
		
		time.Sleep(800 * time.Millisecond)

		if err != nil {
			return err
		}

		parts := strings.Split(translatedJoined, "\n---\n")
		for i, idx := range indices {
			if i < len(parts) {
				translatedLines[idx] = prefixes[i] + strings.TrimSpace(parts[i])
			} else {
				translatedLines[idx] = lines[idx]
			}
		}

		processedLines += len(indices)
		if onProgress != nil {
			percent := (processedLines * 100) / len(lines)
			onProgress(percent)
		}
		return nil
	}

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || (strings.HasPrefix(trimmed, "[[") && strings.HasSuffix(trimmed, "]]")) {
			translatedLines[i] = line
			continue
		}

		if currentBatchLen + len(line) > maxBatchChars && len(currentBatch) > 0 {
			if err := processBatch(currentBatch); err != nil { lastErr = err }
			currentBatch = []int{}
			currentBatchLen = 0
		}
		currentBatch = append(currentBatch, i)
		currentBatchLen += len(line)
	}
	if err := processBatch(currentBatch); err != nil { lastErr = err }

	if lastErr != nil && len(lines) > 5 {
		return "", fmt.Errorf("另存翻譯失敗: %v", lastErr)
	}

	translatedMD := strings.Join(translatedLines, "\n")

	// 4. Restore Links and translate Alt-text
	for i, original := range links {
		placeholder := fmt.Sprintf("[[LP_%d]]", i)
		subMatches := linkRegex.FindStringSubmatch(original)
		if len(subMatches) >= 6 {
			prefix, altText, middle, urlPart, suffix := subMatches[1], subMatches[2], subMatches[3], subMatches[4], subMatches[5]
			if altText != "" {
				tAlt, _ := Translate(altText, source, target)
				time.Sleep(200 * time.Millisecond) 
				translatedMD = strings.ReplaceAll(translatedMD, placeholder, prefix+tAlt+middle+urlPart+suffix)
			} else {
				translatedMD = strings.ReplaceAll(translatedMD, placeholder, original)
			}
		} else {
			translatedMD = strings.ReplaceAll(translatedMD, placeholder, original)
		}
	}

	// 5. Restore code blocks
	for i, original := range codeBlocks {
		placeholder := fmt.Sprintf("[[CP_%d]]", i)
		translatedMD = strings.ReplaceAll(translatedMD, placeholder, original)
	}

	return translatedMD, nil
}

// isCodeBlock returns true if node is inside a <pre>/<code> block
func isCodeBlock(n *html.Node) bool {
	for p := n.Parent; p != nil; p = p.Parent {
		if strings.ToLower(p.Data) == "pre" || strings.ToLower(p.Data) == "code" {
			return true
		}
	}
	return false
}

// GetSupportedLanguages returns supported language list
func GetSupportedLanguages() []map[string]string {
	return []map[string]string{
		{"code": "auto", "name": "Auto-detect"},
		{"code": "en", "name": "English"},
		{"code": "zh-TW", "name": "繁體中文"},
		{"code": "zh-CN", "name": "简体中文"},
		{"code": "ja", "name": "日本語"},
		{"code": "ko", "name": "한국어"},
		{"code": "es", "name": "Español"},
		{"code": "fr", "name": "Français"},
		{"code": "de", "name": "Deutsch"},
		{"code": "ru", "name": "Русский"},
		{"code": "ar", "name": "العربية"},
		{"code": "pt", "name": "Português"},
		{"code": "it", "name": "Italiano"},
		{"code": "nl", "name": "Nederlands"},
		{"code": "pl", "name": "Polski"},
		{"code": "uk", "name": "Українська"},
	}
}

// =============================================================================
// Utilities
// =============================================================================

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
