package utils

import (
	"github.com/mihneamanolache/cert-cli/internal/types"
	"encoding/json"
    "strings"
	"fmt"
	"io"
)

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

// ConvertSetToList converts a map[string]struct{} to a list of strings.
func ConvertSetToList(set map[string]struct{}) []string {
	list := []string{}
	for item := range set {
		list = append(list, item)
	}
	return list
}

// ConvertAddressSetToList converts a set of addresses to a list of Address structs.
func ConvertAddressSetToList(addressSet map[types.Address]struct{}) []types.Address {
	addressList := []types.Address{}
	for address := range addressSet {
		addressList = append(addressList, address)
	}
	return addressList
}

// ConvertAddressSetToStrings converts a set of addresses to a list of formatted strings.
func ConvertAddressSetToStrings(addressSet map[types.Address]struct{}) []string {
	stringList := []string{}
	for addr := range addressSet {
		// Use filterNonEmpty to remove empty parts from the address
		parts := filterNonEmpty([]string{addr.Postcode, addr.Street, addr.County, addr.Country})
		formattedAddress := strings.Join(parts, " ")
		stringList = append(stringList, formattedAddress)
	}
	return stringList
}

// PrintIndentedSet prints a map in an indented format.
func PrintIndentedSet(name string, set map[string]struct{}) {
	fmt.Printf("- %s:\n", name)
	for value := range set {
		fmt.Printf("  - %s\n", value)
	}
}

// SaveToJSON writes the query result to a JSON file.
func SaveToJSON(file io.Writer, data types.QueryResult) error {
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
