package main

import (
	"flag"
	"fmt"
	"github.com/mihneamanolache/cert-cli/internal/feed"
	"github.com/mihneamanolache/cert-cli/internal/types"
	"github.com/mihneamanolache/cert-cli/internal/utils"
	"log"
	"os"
)

func main() {
	// Parse command-line flags
	query := flag.String("q", "", "organization to search for in the certificate's subject field")
	match := flag.String("match", "LIKE", "Match type (allowed values: =, ILIKE, LIKE, single, any, FTS)")
	proxy := flag.String("proxy", "", "proxy")
	jsonOutput := flag.String("o", "", "save findings to JSON file")
	flag.Parse()

	if *query == "" {
		log.Fatalf(types.Red + "You must provide a query using --q flag" + types.Reset)
	}

	feedURL := feed.ConstructFeedURL(*query, *match)
	fmt.Printf(types.Bold + "[i] Fetching feed: %s\n" + types.Reset, feedURL)
	atomFeed, err := feed.FetchFeedWithRetry(feedURL, 3, *proxy)
	if err != nil {
		log.Fatalf(types.Red + "Error fetching feed: %v" + types.Reset, err)
	}

	certificates, entryURLs := feed.ProcessFeed(atomFeed)

	fmt.Printf(types.Bold + "[i] Query: %s\n" + types.Reset, *query)
	fmt.Printf("[i] Parsed %d certificates\n\n", len(entryURLs))

    utils.PrintResults(certificates)

	// Save results to JSON if requested
	if *jsonOutput != "" {
		results := types.QueryResult{
			Query:        *query,
			URL:          feedURL,
			Certificates: certificates,
		}

		file, err := os.Create(fmt.Sprintf("%s.json", *jsonOutput))
		if err != nil {
			log.Fatalf(types.Red + "Error creating JSON file: %v" + types.Reset, err)
		}
		defer file.Close()

		utils.SaveToJSON(file, results)
		fmt.Println(types.Yellow + "\n[i] Findings have been saved to " + file.Name() + types.Reset)
	}
}

