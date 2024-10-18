package utils

import (
	"encoding/json"
	"fmt"
	"github.com/mihneamanolache/cert-cli/internal/types"
	"io"
    "strings"
)

// escapeXML replaces special characters with their corresponding XML escape codes.
func EscapeXML(input string) string {
    replacer := strings.NewReplacer(
		// "&", "&amp;",
		// "<", "&lt;",
		// ">", "&gt;",
		// "\"", "&quot;",
		// "'", "&apos;",
	)
	return replacer.Replace(input)
}

// SaveToJSON writes the query result to a JSON file.
func SaveToJSON(file io.Writer, data types.QueryResult) error {
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// PrintResults prints the certificate information to the console.
func PrintResults(certificates []types.Certificate) {
	if len(certificates) == 0 {
		fmt.Println("No certificates found.")
		return
	}

	// Create sets to aggregate unique organizations, addresses, and domains
	organizations := make(map[string]struct{})
	addresses := make(map[string]struct{})
	domains := make(map[string]struct{})
	sans := make(map[string]struct{})

	// Loop through certificates and populate the sets
	for _, cert := range certificates {
		// Add organization if not empty
		if cert.Organization != "" {
			organizations[cert.Organization] = struct{}{}
		}

		// Add address if not empty
		if cert.Address != "" {
			addresses[cert.Address] = struct{}{}
		}

		// Add common name (domain)
		if cert.CommonName != "" {
			domains[cert.CommonName] = struct{}{}
		}

		// Add all SANs (Subject Alternative Names)
		for _, san := range cert.SAN {
			sans[san] = struct{}{}
		}
	}

	// Remove domains from the SANs set (to leave only SANs that are not in domains)
	for domain := range domains {
		delete(sans, domain)
	}

	// Print aggregated results
	fmt.Println(types.Bold + "Found:" + types.Reset)

	// Print organizations
	fmt.Println("- Organizations:")
	if len(organizations) > 0 {
		for org := range organizations {
			fmt.Printf("  - %s\n", org)
		}
	} else {
		fmt.Println("  - None")
	}

	// Print addresses
	fmt.Println("- Addresses:")
	if len(addresses) > 0 {
		for addr := range addresses {
			fmt.Printf("  - %s\n", addr)
		}
	} else {
		fmt.Println("  - None")
	}

	// Print domains
	fmt.Println("- Domains:")
	if len(domains) > 0 {
		for domain := range domains {
			fmt.Printf("  - %s\n", domain)
		}
	} else {
		fmt.Println("  - None")
	}

	// Print SANs (only those that are not listed as domains)
	fmt.Println("- Subject Alternative Names (SANs not in Domains):")
	if len(sans) > 0 {
		for san := range sans {
			fmt.Printf("  - %s\n", san)
		}
	} else {
		fmt.Println("  - None")
	}
}
