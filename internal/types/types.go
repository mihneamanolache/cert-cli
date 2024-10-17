package types

// AtomEntry represents a single entry in the Atom feed.
type AtomEntry struct {
	ID      string `xml:"id"`
	Summary string `xml:"summary"`
	Title   string `xml:"title"`
}

// AtomFeed represents the Atom feed.
type AtomFeed struct {
	Entries []AtomEntry `xml:"entry"`
}

// Certificate represents the certificate data extracted from a feed.
type Certificate struct {
    URL              string   `json:"url"`
	Organization     string   `json:"organization"`
	CommonName       string   `json:"commonName"`
	SAN              []string `json:"san"`
	Address          string   `json:"address"`
	Issuer           string   `json:"issuer"`
	SerialNumber     string   `json:"serialNumber"`
	NotBefore        string   `json:"notBefore"`
	NotAfter         string   `json:"notAfter"`
	KeyUsage         []string `json:"keyUsage"`
	SignatureAlgorithm string `json:"signatureAlgorithm"`
	Version          int      `json:"version"`
}

// QueryResult holds all the findings for JSON output.
type QueryResult struct {
	Query        string        `json:"query"`
	URL          string        `json:"url"`
	Certificates []Certificate `json:"certificates"`
}
