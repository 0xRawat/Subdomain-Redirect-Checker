# redirect_check
A fast and concurrent redirect checker. This tool uses a headless browser to detect client-side (JavaScript-based) and server-side redirects from a list of subdomains.

# 🚦 Headless Redirect Checker using ChromeDP

A fast and concurrent redirect checker built using [chromedp](https://github.com/chromedp/chromedp).

---

## 🔧 Features

- ✅ Detects both server-side and JavaScript-based redirects
- 🚀 Uses `chromedp` with headless Chrome for accurate navigation
- 🔁 Shows only actual redirected subdomains (not same-host redirects)
- 🧠 Automatically adds `http://` scheme if missing
- 🌈 Color-coded terminal output for clarity
- 🧵 Concurrent with configurable goroutine limit


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

Works best on public internet-facing domains.
Avoid using on very large lists without rate-limiting or proxies.
Requires a working Chrome/Chromium installation (or compatible headless browser).


---

## ❤️ Contributions
Contributions, issues, and feature requests are welcome!
