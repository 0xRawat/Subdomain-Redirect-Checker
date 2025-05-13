package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

type RedirectChain struct {
	Original string
	Chain    []string
}

var safeRedirects = []func(string, string) bool{
	isWWWRedirect,
}

// === Main Execution ===

func main() {
	inputFile := flag.String("l", "", "Input file with list of subdomains")
	outputFile := flag.String("o", "redirects.txt", "Output file to save redirected subdomains")
	flag.Parse()

	if *inputFile == "" {
		fmt.Println("[-] Please provide an input file with -l flag")
		os.Exit(1)
	}

	// Read input
	domains, err := readLines(*inputFile)
	if err != nil {
		fmt.Printf("[-] Error reading file: %v\n", err)
		os.Exit(1)
	}

	var mu sync.Mutex
	grouped := make(map[string][]string)
	var wg sync.WaitGroup
	sem := make(chan struct{}, 5)

	fmt.Println("🔍 Checking redirect chains...")

	for _, domain := range domains {
		wg.Add(1)
		sem <- struct{}{}

		go func(domain string) {
			defer wg.Done()
			defer func() { <-sem }()

			chain := getRedirectChain(domain)
			if len(chain) < 2 {
				return
			}

			origHost := extractHost(chain[0])
			finalHost := extractHost(chain[len(chain)-1])

			if isSafeRedirect(origHost, finalHost) {
				return
			}

			color := "\033[33m" // Yellow
			if !strings.HasPrefix(chain[0], "https://") && strings.HasPrefix(chain[len(chain)-1], "https://") {
				color = "\033[32m" // Green for HTTP→HTTPS
			}

			fmt.Printf("🔁 %s \033[31m→\033[0m %s%s\033[0m\n", domain, color, finalHost)

			mu.Lock()
			grouped[finalHost] = append(grouped[finalHost], domain)
			mu.Unlock()
		}(domain)
	}

	wg.Wait()
	writeGroupedOutput(grouped, *outputFile)
	fmt.Printf("\n✅ Done. Grouped redirect info saved to %s\n", *outputFile)
}

// === Redirect Logic ===

func getRedirectChain(domain string) []string {
	var chain []string

	for _, scheme := range []string{"http://", "https://"} {
		startURL := ensureScheme(scheme, domain)
		ctx, cancel := chromedp.NewContext(context.Background())
		ctx, cancelTimeout := context.WithTimeout(ctx, 15*time.Second)

		var finalURL string
		err := chromedp.Run(ctx,
			chromedp.Navigate(startURL),
			chromedp.Sleep(4*time.Second),
			chromedp.Location(&finalURL),
		)

		cancelTimeout()
		cancel()

		if err == nil && finalURL != "" && finalURL != startURL {
			chain = append(chain, startURL, finalURL)
			// Optionally follow more redirects here if needed
			break
		}
	}
	return chain
}

func ensureScheme(scheme, url string) string {
	url = strings.TrimSpace(url)
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return scheme + url
	}
	return url
}

func extractHost(rawurl string) string {
	u, err := url.Parse(rawurl)
	if err != nil {
		return ""
	}
	return u.Hostname()
}

// === Helpers ===

func readLines(filename string) ([]string, error) {
	var lines []string
	file, err := os.Open(filename)
	if err != nil {
		return lines, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text != "" {
			lines = append(lines, text)
		}
	}
	return lines, scanner.Err()
}

func writeGroupedOutput(data map[string][]string, outputFile string) {
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("[-] Failed to create output file: %v\n", err)
		return
	}
	defer file.Close()

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, finalHost := range keys {
		file.WriteString(fmt.Sprintf("\n%s redirects\n", finalHost))
		for _, sub := range data[finalHost] {
			file.WriteString(fmt.Sprintf("%s\n", sub))
		}
	}
}

// === Filtering Safe Redirects ===

func isSafeRedirect(from, to string) bool {
	for _, f := range safeRedirects {
		if f(from, to) {
			return true
		}
	}
	return false
}

func isWWWRedirect(from, to string) bool {
	if strings.HasPrefix(from, "www.") {
		from = strings.TrimPrefix(from, "www.")
	}
	if strings.HasPrefix(to, "www.") {
		to = strings.TrimPrefix(to, "www.")
	}
	return from == to
}
