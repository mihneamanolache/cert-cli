package feed

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"cert-cli/internal/types"
	"fmt"
	"regexp"
	"strings"
	"time"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
)

// ConstructFeedURL builds the Atom feed URL based on the query and match type.
func ConstructFeedURL(query string, match string) string {
	encodedQuery := url.QueryEscape(query)
	return fmt.Sprintf("https://crt.sh/atom?q=%s&match=%s", encodedQuery, match)
}

// ProcessFeed processes the Atom feed and extracts data from each certificate.
func ProcessFeed(feed *types.AtomFeed) (map[string]struct{}, map[types.Address]struct{}, map[string]struct{}, []string) {
	organizations := make(map[string]struct{})
	addresses := make(map[types.Address]struct{})
	domains := make(map[string]struct{})
	entryURLs := []string{}

	for _, entry := range feed.Entries {
		entryURLs = append(entryURLs, entry.ID)
		certPem, err := extractCertFromSummary(entry.Summary)
		if err != nil {
			fmt.Printf(types.Red+"Error extracting certificate: %v\n"+types.Reset, err)
			continue
		}
		parseAndCollectInfo(certPem, organizations, addresses, domains)
	}
	return organizations, addresses, domains, entryURLs
}

// FetchFeed fetches the Atom feed from the provided URL.
func FetchFeed(feedURL string, proxyURL string) (*types.AtomFeed, error) {
    var client *http.Client
    if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL: %w", err)
		}
		// Set up the transport with the proxy
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
        transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
		client = &http.Client{
			Transport: transport,
		}
        fmt.Println( types.Yellow + "[i] Using proxy for this request!" + types.Reset)
	} else {
		// Use default HTTP client (no proxy)
		client = &http.Client{}
	}
	resp, err := client.Get(feedURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the feed: %w", err)
	}
	defer resp.Body.Close()

    
    status := resp.StatusCode
    if status != 200 {
        return nil, fmt.Errorf("failed to fetch the feed. status_code: %d", status)
    }

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var feed types.AtomFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("failed to unmarshal XML: %w", err)
	}

	return &feed, nil
}

// FetchFeedWithRetry attempts to fetch the Atom feed, retrying on failure.
func FetchFeedWithRetry(feedURL string, retries int, proxyURL string) (*types.AtomFeed, error) {
	var feed *types.AtomFeed
	var err error

	for i := 0; i <= retries; i++ {
		feed, err = FetchFeed(feedURL, proxyURL)
		if err == nil {
			return feed, nil
		}
		fmt.Printf(types.Red+"Error fetching feed (attempt %d/%d): %v\n"+types.Reset, i+1, retries+1, err)
		time.Sleep(5 * time.Second)
	}
	return nil, fmt.Errorf("failed to fetch feed after %d attempts: %w", retries+1, err)
}

func extractCertFromSummary(summary string) (string, error) {
	re := regexp.MustCompile(`-----BEGIN CERTIFICATE-----(.*?)-----END CERTIFICATE-----`)
	matches := re.FindStringSubmatch(summary)
	if len(matches) > 0 {
		cert := strings.ReplaceAll(matches[0], "<br>", "\n")
		return cert, nil
	}
	return "", fmt.Errorf("no certificate found in the summary")
}

func parseAndCollectInfo(certPem string, organizations map[string]struct{}, addresses map[types.Address]struct{}, domains map[string]struct{}) error {
	block, _ := pem.Decode([]byte(certPem))
	if block == nil || block.Type != "CERTIFICATE" {
		return fmt.Errorf("failed to decode PEM block containing certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}

	if cert.Subject.CommonName != "" {
		domains[cert.Subject.CommonName] = struct{}{}
	}
	for _, org := range cert.Subject.Organization {
		organizations[org] = struct{}{}
	}

	address := types.Address{
		Postcode: strings.Join(cert.Subject.PostalCode, " "),
		Street:   strings.Join(cert.Subject.StreetAddress, " "),
		County:   strings.Join(cert.Subject.Province, " "),
		Country:  strings.Join(cert.Subject.Country, " "),
	}
	if address != (types.Address{}) {
		addresses[address] = struct{}{}
	}

	return nil
}
