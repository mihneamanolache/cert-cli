package types

// MatchType defines the allowed match types for the search.
type MatchType string

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

// Address represents a structured address.
type Address struct {
	Postcode string `json:"postcode"`
	Street   string `json:"street"`
	County   string `json:"county"`
	Country  string `json:"country"`
}

// QueryResult holds all the findings for JSON output.
type QueryResult struct {
	Query         string    `json:"query"`
	URL           string    `json:"url"`
	Organizations []string  `json:"organizations"`
	Addresses     []Address `json:"addresses"`
	Domains       []string  `json:"domains"`
	Entries       []string  `json:"entries"`
}
