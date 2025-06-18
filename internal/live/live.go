package live

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/mihneamanolache/cert-cli/internal/types"
	"strings"
	"time"
)

// CheckDomain performs a TLS connection to the given domain and returns
// parsed certificate information.
func CheckDomain(domain string) (types.Certificate, error) {
	conn, err := tls.Dial("tcp", domain+":443", &tls.Config{
		ServerName:         domain,
		InsecureSkipVerify: true,
	})
	if err != nil {
		return types.Certificate{}, err
	}
	defer conn.Close()

	state := conn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		return types.Certificate{}, fmt.Errorf("no certificate presented")
	}

	cert := state.PeerCertificates[0]

	var keyUsage []string
	if cert.KeyUsage&x509.KeyUsageDigitalSignature != 0 {
		keyUsage = append(keyUsage, "Digital Signature")
	}
	if cert.KeyUsage&x509.KeyUsageKeyEncipherment != 0 {
		keyUsage = append(keyUsage, "Key Encipherment")
	}

	address := strings.Join(filterNonEmpty([]string{
		strings.Join(cert.Subject.PostalCode, " "),
		strings.Join(cert.Subject.StreetAddress, " "),
		strings.Join(cert.Subject.Province, " "),
		strings.Join(cert.Subject.Country, " "),
	}), " ")

	return types.Certificate{
		URL:                domain,
		Organization:       strings.Join(cert.Subject.Organization, ", "),
		CommonName:         cert.Subject.CommonName,
		SAN:                cert.DNSNames,
		Address:            address,
		Issuer:             cert.Issuer.CommonName,
		SerialNumber:       cert.SerialNumber.String(),
		NotBefore:          cert.NotBefore.Format(time.RFC3339),
		NotAfter:           cert.NotAfter.Format(time.RFC3339),
		KeyUsage:           keyUsage,
		SignatureAlgorithm: cert.SignatureAlgorithm.String(),
		Version:            cert.Version,
	}, nil
}

func filterNonEmpty(parts []string) []string {
	var result []string
	for _, part := range parts {
		if strings.TrimSpace(part) != "" {
			result = append(result, part)
		}
	}
	return result
}
