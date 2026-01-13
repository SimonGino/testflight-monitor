package monitor

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// TestFlightInfo contains parsed info from TestFlight page
type TestFlightInfo struct {
	AppID     string
	AppName   string
	IconURL   string
	Available bool // true if beta has open slots
	Message   string
}

// Checker handles TestFlight availability checking
type Checker struct {
	client *http.Client
}

// NewChecker creates a new checker with optional proxy
func NewChecker(proxyURL string) *Checker {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err == nil {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxy),
			}
		}
	}

	return &Checker{client: client}
}

func ParseURL(testFlightURL string) (string, error) {
	re := regexp.MustCompile(`testflight\.apple\.com/join/([a-zA-Z0-9]+)`)
	matches := re.FindStringSubmatch(testFlightURL)
	if len(matches) < 2 {
		return "", fmt.Errorf("invalid TestFlight URL: %s", testFlightURL)
	}
	return matches[1], nil
}

func parseAppNameFromTitle(title string) string {
	title = strings.TrimSpace(title)
	title = strings.TrimPrefix(title, "Join the ")
	title = strings.TrimSuffix(title, " beta")
	title = strings.TrimSuffix(title, " - TestFlight - Apple")
	return strings.TrimSpace(title)
}

// Check fetches TestFlight page and parses availability
func (c *Checker) Check(appID string) (*TestFlightInfo, error) {
	url := fmt.Sprintf("https://testflight.apple.com/join/%s", appID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set headers to mimic browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	info := &TestFlightInfo{
		AppID: appID,
	}

	ogTitle, _ := doc.Find("meta[property='og:title']").Attr("content")
	if ogTitle != "" {
		info.AppName = parseAppNameFromTitle(ogTitle)
	}
	if info.AppName == "" {
		twitterTitle, _ := doc.Find("meta[name='twitter:title']").Attr("content")
		if twitterTitle != "" {
			info.AppName = parseAppNameFromTitle(twitterTitle)
		}
	}
	if info.AppName == "" {
		titleText := doc.Find("title").First().Text()
		info.AppName = parseAppNameFromTitle(titleText)
	}

	ogImage, _ := doc.Find("meta[property='og:image']").Attr("content")
	if ogImage != "" {
		info.IconURL = ogImage
	}
	if info.IconURL == "" {
		twitterImage, _ := doc.Find("meta[name='twitter:image']").Attr("content")
		if twitterImage != "" {
			info.IconURL = twitterImage
		}
	}

	// Check availability by looking for specific messages
	pageText := strings.ToLower(doc.Text())

	// "This beta is full" or "beta is full" means not available
	if strings.Contains(pageText, "beta is full") {
		info.Available = false
		info.Message = "Beta is full"
		return info, nil
	}

	// "This beta isn't accepting" means not available
	if strings.Contains(pageText, "isn't accepting") || strings.Contains(pageText, "not accepting") {
		info.Available = false
		info.Message = "Not accepting testers"
		return info, nil
	}

	// Look for "Start Testing" or "Accept" button - means available
	hasAcceptButton := false
	doc.Find("a, button").Each(func(i int, s *goquery.Selection) {
		text := strings.ToLower(s.Text())
		if strings.Contains(text, "start testing") || strings.Contains(text, "accept") {
			hasAcceptButton = true
		}
	})

	if hasAcceptButton {
		info.Available = true
		info.Message = "Beta available"
		return info, nil
	}

	// Check for "View in App Store" which might indicate a redirect
	if strings.Contains(pageText, "view in app store") || strings.Contains(pageText, "open in testflight") {
		info.Available = true
		info.Message = "Available"
		return info, nil
	}

	// Default to unavailable if we can't determine
	info.Available = false
	info.Message = "Unable to determine availability"

	return info, nil
}
