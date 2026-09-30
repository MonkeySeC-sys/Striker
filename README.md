# DNS Exfiltration Toolkit — A Covert Channel via TXT Records

[![Go Version](https://shields.io)](https://go.dev)
[![License](https://shields.io)](LICENSE)
[![Purpose](https://shields.io)](#disclaimer)

A lightweight Go-based utility for encoding arbitrary data into DNS queries and decoding exfiltrated payloads from TXT records. Designed for security research, penetration testing, and incident response training, this toolkit demonstrates how DNS tunneling can be leveraged as a covert channel when traditional C2 paths are noisy or blocked.

---

## 📖 Table of Contents
- [Introduction](#-introduction)
- [Technical Details](#-technical-details)
  - [Encoding Logic](#encoding-logic)
  - [Decoding Logic](#decoding-logic)
  - [Resolver Configuration](#resolver-configuration)
- [Installation & Setup](#-installation--setup)
- [Usage](#-usage)
  - [Encode a Payload](#encode-a-payload)
  - [Decode from DNS](#decode-from-dns)
- [Security Considerations](#-security-considerations)
- [Intended Use Cases](#-intended-use-cases)
- [Disclaimer & Legal](#-disclaimer--legal)

---

## 🚀 Introduction

In modern networked environments, defenders increasingly monitor outbound traffic to detect data exfiltration attempts. However, DNS remains one of the few protocols that is universally permitted through firewalls and rarely inspected at depth. 

This toolkit exploits that trust by encoding arbitrary payloads into DNS queries—specifically using TXT records—and decoding them on the receiving side. The approach is intentionally simple: no custom protocol, just standard DNS over UDP/TCP with a configurable resolver address.

---

## ⚙️ Technical Details

### Encoding Logic
The `Encode` function takes raw bytes, converts them to hexadecimal (two hex chars per byte), then splits the result into chunks of **254 characters**—the maximum length allowed in a single DNS label. Each chunk becomes a subdomain part. Parts are joined with dots and suffixed with the configured exfiltration domain (`exfil.example.com` by default). This yields a fully qualified domain name (FQDN) that, when queried for TXT records, returns the encoded data.

### Decoding Logic
The `Decode` function reverses this process: it strips the suffix, adds a leading dot if needed to ensure proper parsing, then splits on periods and queries each part via a resolver configured with a 30-second dial timeout. The TXT records are concatenated in order and hex-decoded back to bytes. Errors during lookup are wrapped and returned for caller inspection.

### Resolver Configuration
By default, the tool uses a `net.Resolver` pointing to `0.0.0.0:53`. Users can override this via the `DNS_ADDR` environment variable (e.g., `export DNS_ADDR=8.8.8.8:53`). The resolver’s dial timeout is hard-coded to 30 seconds; on high-latency networks this may cause false negatives, so users should consider adjusting it if needed.

---

## 🛠️ Installation & Setup

Ensure you have [Go](https://go.devdoc/install) installed on your system.

1. **Clone the repository:**
   ```bash
   git clone https://github.com
   cd dns-exfil-toolkit
   ```

2. **Build the binary:**
   ```bash
   go build -o Striker Striker.go
   ```

---

## 💻 Usage

### Encode a Payload
By default, the utility uses `exfil.example.com` as a placeholder (replace this in the source code before deployment).

```bash
./Striker
```
**Example Output:**
```text
Encode: deadbeef -> deadbeef.deadbeef.exfil.example.com
```

### Decode from DNS
Requires a live resolver. By default, it queries `0.0.0.0:53`—override this by setting the `DNS_ADDR` variable:

```bash
export DNS_ADDR=8.8.8.8:53
./Striker decode <encoded-dns-string>
```

**Example Output (Before DNS configuration):**
```text
decode error: decode: lookup TXT deadbeef.deadbeef.exfil.example.com: no such host
```
*(Note: Decode will fail until the targeted domain is registered and populated with the appropriate TXT records containing the payload.)*

---

## 🛡️ Security Considerations

* **Exfil Domain:** The constant `exfilDomain` is a placeholder. Always swap it with your actual, authorized C2/logging endpoint before deployment.
* **Resolver Dial Timeout:** Hardcoded at 30s. On high-latency networks, this can cause false negatives. Adjust the resolver's `Dial` function in the source code if timeouts occur.
* **Visibility:** DNS queries are logged by most enterprise resolvers and firewalls. Expect detection unless you are routing through a stealthy or trusted upstream path.

---

## 🎯 Intended Use Cases

* **Penetration Testing:** Simulate data exfiltration paths to validate detection controls and network segmentation.
* **Incident Response Training:** Practice identifying anomalous DNS query patterns during Blue Team exercises or tabletop scenarios.
* **Threat Modeling:** Evaluate the feasibility of DNS-based covert channels against your infrastructure.
* **Red Team Operations:** Test adversary emulation scenarios where traditional web-based C2 channels are blocked.

---

## ⚠️ Disclaimer & Legal

This tool is intended for **authorized security research and educational purposes only**. Users must ensure they have explicit, written authorization per local laws (e.g., Computer Fraud and Abuse Act in the US, Computer Misuse Act in the UK). 

Always document scope and targets, and obtain explicit permission before testing any environment not owned by you or your organization. The author assumes no liability for misuse or damage caused by this utility.
