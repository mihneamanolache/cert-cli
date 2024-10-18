package feed

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/mihneamanolache/cert-cli/internal/types"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
	"encoding/xml"
	"io"
)

// ConstructFeedURL builds the Atom feed URL based on the query and match type.
func ConstructFeedURL(query string, match string) string {
	encodedQuery := url.QueryEscape(query)
	return fmt.Sprintf("https://crt.sh/atom?q=%s&match=%s", encodedQuery, match)
}

// filterNonEmpty removes empty parts from a slice of strings.
func filterNonEmpty(parts []string) []string {
	var result []string
	for _, part := range parts {
		if strings.TrimSpace(part) != "" { // Ensure that parts containing only whitespaces are excluded
			result = append(result, part)
		}
	}
	return result
}

// parseCertInfo parses the X.509 certificate and collects relevant information.
func parseCertInfo(certPem string, entryURL string) (types.Certificate, error) {
	block, _ := pem.Decode([]byte(certPem))
	if block == nil || block.Type != "CERTIFICATE" {
		return types.Certificate{}, fmt.Errorf("failed to decode PEM block containing certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return types.Certificate{}, fmt.Errorf("failed to parse certificate: %w", err)
	}

	// Extract fields
	organization := strings.Join(cert.Subject.Organization, ", ")
	commonName := cert.Subject.CommonName
	issuer := cert.Issuer.CommonName
	serialNumber := cert.SerialNumber.String()
	notBefore := cert.NotBefore.Format(time.RFC3339)
	notAfter := cert.NotAfter.Format(time.RFC3339)
	version := cert.Version

	// SAN (Subject Alternative Names)
	var san []string
	for _, dns := range cert.DNSNames {
		san = append(san, dns)
	}

	// Key Usage
	var keyUsage []string
	if cert.KeyUsage&x509.KeyUsageDigitalSignature != 0 {
		keyUsage = append(keyUsage, "Digital Signature")
	}
	if cert.KeyUsage&x509.KeyUsageKeyEncipherment != 0 {
		keyUsage = append(keyUsage, "Key Encipherment")
	}

	// Address
	address := strings.Join(filterNonEmpty([]string{
		strings.Join(cert.Subject.PostalCode, " "),
		strings.Join(cert.Subject.StreetAddress, " "),
		strings.Join(cert.Subject.Province, " "),
		strings.Join(cert.Subject.Country, " "),
	}), " ")

	// Signature Algorithm
	signatureAlgorithm := cert.SignatureAlgorithm.String()

	// Create the Certificate object
	certificate := types.Certificate{
		URL:               entryURL,
		Organization:      organization,
		CommonName:        commonName,
		SAN:               san,
		Address:           address,
		Issuer:            issuer,
		SerialNumber:      serialNumber,
		NotBefore:         notBefore,
		NotAfter:          notAfter,
		KeyUsage:          keyUsage,
		SignatureAlgorithm: signatureAlgorithm,
		Version:           version,
	}

	return certificate, nil
}

// ProcessFeed processes the Atom feed and extracts data from each certificate.
func ProcessFeed(feed *types.AtomFeed) ([]types.Certificate, []string) {
    var certificates []types.Certificate
    entryURLs := []string{}

    for _, entry := range feed.Entries {
        entryURLs = append(entryURLs, entry.ID)

        certPem, err := extractCertFromSummary(entry.Summary)
        if err != nil {
            fmt.Printf(types.Red+"Error extracting certificate: %v\n"+types.Reset, err)
            continue
        }

        certInfo, err := parseCertInfo(certPem, entry.ID)
        if err != nil {
            fmt.Printf(types.Red+"Error parsing certificate: %v\n"+types.Reset, err)
            continue
        }

        certificates = append(certificates, certInfo)
    }

    return certificates, entryURLs
}

// FetchFeed fetches the Atom feed from the provided URL.
func FetchFeed(feedURL string, proxyURL string) (*types.AtomFeed, error) {
	var client *http.Client
	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL: %w", err)
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
		client = &http.Client{
			Transport: transport,
		}
		fmt.Println(types.Yellow + "[i] Using proxy for this request!" + types.Reset)
	} else {
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

	decoder := xml.NewDecoder(strings.NewReader(string(body)))
	decoder.Strict = false

	var feed types.AtomFeed
	if err := decoder.Decode(&feed); err != nil {
		return nil, fmt.Errorf("failed to unmarshal XML: %w", err)
	}

	return &feed, nil
}

// FetchFeedWithRetry attempts to fetch the Atom feed, retrying on failure.
func FetchFeedWithRetry(feedURL string, retries int, proxyURL string) (*types.AtomFeed, error) {
	var feed *types.AtomFeed
	var err error

	for i := 0; i <= retries; i++ {
        // if we have more than 2 retries, exclude expired certificates to increase the chances of getting a response
        if i > 2 {
			feedURL += "&exclude=expired"
		}
        if i > 3 {
            feedURL += "&deduplicate=Y"
        }
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

func parseCertificate(certPem string, entryID string) (types.Certificate, error) {
	block, _ := pem.Decode([]byte(certPem))
	if block == nil || block.Type != "CERTIFICATE" {
		return types.Certificate{}, fmt.Errorf("failed to decode PEM block containing certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return types.Certificate{}, fmt.Errorf("failed to parse certificate: %w", err)
	}

	// Construct the concatenated address
	address := strings.Join([]string{
		strings.Join(cert.Subject.PostalCode, " "),
		strings.Join(cert.Subject.StreetAddress, " "),
		strings.Join(cert.Subject.Province, " "),
		strings.Join(cert.Subject.Country, " "),
	}, " ")

	return types.Certificate{
		URL:          entryID,
		Organization: strings.Join(cert.Subject.Organization, ", "),
		CommonName:   cert.Subject.CommonName,
		SAN:          cert.DNSNames,
		Address:      address,
	}, nil
}
