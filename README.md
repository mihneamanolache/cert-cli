# `cert-cli`: OSINT Certificate Analysis CLI Tool

`cert-cli` is an OSINT command-line tool that allows you to query and retrieve certificate information from the public certificate transparency logs, leveraging [crt.sh](https://crt.sh/), which currently hosts over **14 billion** certificates.

This tool helps you gather domains, organizations, and addresses from X.509 certificates associated with a particular company or domain name. It is useful for cybersecurity investigations, security research, and general OSINT purposes.

## Features
- Search for certificates based on domain or organization names
- Supports proxy configurations for enhanced privacy or bypassing restrictions.
- Extract unique domains, organization names, addresses, and other metadata from certificates
- Supports multiple query types:  `=`, `ILIKE`, `LIKE`, `single`, `any`, `FTS`
- Save the results in JSON format
- Cross-platform support (Linux, macOS, Windows)

## Installation

### Download the Binary
You can download the latest binary from the [releases page](https://github.com/mihneamanolache/cert-cli/releases) and add it to your PATH.
```bash
# Download the binary
wget https://github.com/mihneamanolache/cert-cli/releases/download/$VERSION/cert-cli-linux-amd64.tar.gz
# Extract the binary
tar -xvzf cert-cli-linux-amd64.tar.gz
# Make the binary executable
chmod +x cert-cli
# Move the binary to your PATH
sudo mv cert-cli /usr/local/bin/cert-cli
```

### Build from Source
To build from source, clone the repository and run the following commands:

```bash
# Clone the repository
git clone https://github.com/mihneamanolache/cert-cli
# Change the directory
cd cert-cli
# Build the binary
go build -o cert-cli cmd/main.go
# Move the binary to your PATH
mv cert-cli /usr/local/bin
```

## Usage
```bash
cert-cli -q "<company or domain>" -match "<query-type>" -o <output-file> -proxy "<proxy-url>"
```

### Flags
- `-q`: The company or domain name to search for
- `-match`: The query type to use. Supported values are `LIKE`, `ANY`, `ILIKE`
- `-o`: The output file to save the results in JSON format
- `-proxy`: The proxy URL to use for the request

### Example Commands
Search for certificates for `Dreamwors` with the match type `LIKE`:
```bash
cert-cli -q "Dreamworks" -match "LIKE" 

# Output
[i] Fetching feed: https://crt.sh/atom?q=Dreamworks&match=LIKE
[i] Query: Dreamworks
[i] Parsed 20 certificates

Found:
- Organizations:
  - Dreamworks Model Products LLC
  - Dreamworks Animation
  - Dreamworks
  - Shadow Dreamworks LLC
  - Dreamworks Animation SKG Inc.
- Addresses:
  - FL US
  - California US
  - Osaka JP
  - Georgia US
  - Florida US
- Domains:
  - www.dreamworks.co.jp
  - owa.dwextra.com
  - chadwrichardson.com
  - saml.dreamworks.com
  - connect.pdi.com
  - testsaml.dreamworks.com
  - *.dreamworks.com
  - connect.ddu.dreamworks.com
  - www.dreamplaceexperience.com
  - www.dreamworksrc.com
  - connect.dreamworks.com
  - *.dreamworksanimation.com
- Subject Alternative Names (SANs not in Domains):
  - shrews.win.dreamworks.com
  - www.chadwrichardson.com
  - dreamworks.com
  - georgia.win.dreamworks.com
  - pool.win.dreamworks.com
  - royals.win.dreamworks.com
  - pilgrims.win.dreamworks.com
  - dreamworksrc.com
  - dreamplaceexperience.com
  - northcarolina.win.dreamworks.com
  - dreamworksanimation.com
  - devilrays.win.dreamworks.com
```

Should you need to save the results in a file, you can use the `-o` flag:
```bash
cert-cli -q "Dreamworks" -match "LIKE" -o dreamworks.json
```
Output:
```json
{
  "query": "Dreamworks",
  "url": "https://crt.sh/atom?q=Dreamworks\u0026match=LIKE",
  "certificates": [
    {
      "url": "https://crt.sh/?id=1167265#q;Dreamworks",
      "organization": "Dreamworks Model Products LLC",
      "commonName": "www.dreamworksrc.com",
      "san": [
        "www.dreamworksrc.com",
        "dreamworksrc.com"
      ],
      "address": "FL US",
      "issuer": "Go Daddy Secure Certification Authority",
      "serialNumber": "1234066529932341",
      "notBefore": "2012-01-13T21:52:25Z",
      "notAfter": "2015-02-20T23:22:40Z",
      "keyUsage": [
        "Digital Signature",
        "Key Encipherment"
      ],
      "signatureAlgorithm": "SHA1-RSA",
      "version": 3
    },
    {
      "url": "https://crt.sh/?id=1583134#q;Dreamworks",
      "organization": "Dreamworks Animation",
      "commonName": "connect.dreamworks.com",
      "san": null,
      "address": "California US",
      "issuer": "VeriSign Class 3 International Server CA - G3",
      "serialNumber": "130848744974652552427047020601007864849",
      "notBefore": "2012-04-04T00:00:00Z",
      "notAfter": "2016-04-24T23:59:59Z",
      "keyUsage": [
        "Digital Signature",
        "Key Encipherment"
      ],
      "signatureAlgorithm": "SHA1-RSA",
      "version": 3
    }
  ]
}
```

## Contributing
We welcome contributions! Feel free to open issues or submit pull requests to improve the tool.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer
This tool is intended for educational and research purposes only. The author does not bear any responsibility for any misuse of the tool. Use it responsibly and at your own risk. We're not asscoiated with crt.sh in any way and we strongly recommend you to read their [Terms of Service](https://crt.sh/tos) before using this tool. 
