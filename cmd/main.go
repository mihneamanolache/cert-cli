package main

import (
	"cert-cli/internal/feed"
	"cert-cli/internal/types"
	"cert-cli/internal/utils"
	"flag"
	"fmt"
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

	// Validate the query
	if *query == "" {
		log.Fatalf(types.Red + "You must provide a query using --q flag" + types.Reset)
	}

	// Construct the feed URL
	feedURL := feed.ConstructFeedURL(*query, *match)

	// Fetch the Atom feed with retries
    fmt.Printf(types.Bold + "[i] Fetching feed: %s\n" + types.Reset, feedURL)
	atomFeed, err := feed.FetchFeedWithRetry(feedURL, 3, *proxy)
	if err != nil {
		log.Fatalf(types.Red + "Error fetching feed: %v" + types.Reset, err)
	}

	// Process the feed
	organizations, addresses, domains, entryURLs := feed.ProcessFeed(atomFeed)

	// Convert address set to string map for output
    stringAddresses := utils.ConvertAddressSetToStrings(addresses)
    stringAddressMap := make(map[string]struct{})
    for _, addr := range stringAddresses {
        stringAddressMap[addr] = struct{}{}
    }

	// Print results
	fmt.Printf(types.Bold + "[i] Query: %s\n" + types.Reset, *query)
	fmt.Printf("[i] Parsed %d certificates\n\n", len(entryURLs))
	fmt.Println(types.Bold + "Found:" + types.Reset)
	utils.PrintIndentedSet("Organizations", organizations)
	utils.PrintIndentedSet("Addresses", stringAddressMap)
	utils.PrintIndentedSet("Domains", domains)

	// Save results to JSON if requested
	if *jsonOutput != "" {
		results := types.QueryResult{
			Query:         *query,
			URL:           feedURL,
			Organizations: utils.ConvertSetToList(organizations),
			Addresses:     utils.ConvertAddressSetToList(addresses),
			Domains:       utils.ConvertSetToList(domains),
			Entries:       entryURLs,
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

