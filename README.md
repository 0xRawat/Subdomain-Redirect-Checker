# Subdomain-Redirect-Checker
A fast and concurrent redirect checker. This Go tool checks for JavaScript-based and meta-refresh redirects across a list of subdomains using a headless browser powered by chromedp. It's ideal for bug bounty recon, redirect detection. 

# 🚦 Headless Redirect Checker using ChromeDP

A fast and concurrent redirect checker built using [chromedp](https://github.com/chromedp/chromedp).

---

##  🚀 Features
- ✅ Detects redirect chains using real browser rendering

- 🌐 Supports both HTTP and HTTPS fallback

- 🎨 Color-coded output:

- 🟨 Normal redirects

- 🟩 HTTP → HTTPS redirects

- 🧠 Filters out known safe redirects like www.domain.com → domain.com

- 📦 Groups output by final redirect destination

- 🧵 Handles multiple subdomains concurrently (configurable)

- 💾 Saves results to redirects.txt in a grouped format


---

## 📦 Requirements

- [Go](https://golang.org/doc/install) 1.18 or higher
- Headless-compatible environment (e.g., desktop or server with Chrome)

---

## 📥 Installation

    git clone https://github.com/yourusername/redirect-check.git
    cd redirect-check
    go build -o redirect-check


---

## 🚀 Usage

    ./redirect-check -l subdomains.txt -o redirects.txt
    
    


| Flag | Description                                                                      |
| ---- | -------------------------------------------------------------------------------- |
| `-l` | Path to input file containing a list of subdomains (one per line)                |
| `-o` | Output file where redirected subdomains will be saved (default: `redirects.txt`) |


---

## ⚠️ Notes

This tool uses a headless browser, so it's slower than curl-based tools — but far more accurate.

The output intentionally ignores trivial changes (like www. removal) to reduce noise.

Works best on public internet-facing domains.

Avoid using on very large lists without rate-limiting or proxies.



---

## ❤️ Contributions
Contributions, issues, and feature requests are welcome!
