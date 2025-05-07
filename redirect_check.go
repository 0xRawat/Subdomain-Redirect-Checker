package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	inputFile := flag.String("l", "", "Input file with list of subdomains")
	outputFile := flag.String("o", "redirects.txt", "Output file to save redirected subdomains")
	flag.Parse()

	if *inputFile == "" {
		fmt.Println("[-] Please provide an input file with -l flag")
		os.Exit(1)
	}

	// Open input file
	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("[-] Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Prepare output file
	out, err := os.Create(*outputFile)
	if err != nil {
		fmt.Printf("[-] Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer out.Close()

	scanner := bufio.NewScanner(file)
	var wg sync.WaitGroup
	sem := make(chan struct{}, 5) // concurrency limit
	mu := sync.Mutex{}

	fmt.Println("🔍 Checking redirections...")

	for scanner.Scan() {
		subdomain := strings.TrimSpace(scanner.Text())
		if subdomain == "" {
			continue
		}

		wg.Add(1)
		sem <- struct{}{}

		go func(domain string) {
			defer wg.Done()
			defer func() { <-sem }()

			finalURL := checkRedirect(domain)
			if finalURL == "" {
				return
			}

			originalHost := extractHost(ensureScheme(domain))
			redirectHost := extractHost(finalURL)

			// Only print/write if host changed
			if redirectHost != "" && originalHost != "" && !strings.EqualFold(originalHost, redirectHost) {
				fmt.Printf("🔁 %s \033[31m->\033[0m \033[33m%s\033[0m\n", domain, redirectHost)
				mu.Lock()
				out.WriteString(fmt.Sprintf("%s -> %s\n", domain, redirectHost))
				mu.Unlock()
			}
		}(subdomain)
	}

	wg.Wait()
	fmt.Println("✅ Done. Redirected subdomains saved to redirected.txt.")
}

func checkRedirect(domain string) string {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var finalURL string
	domain = ensureScheme(domain)
	err := chromedp.Run(ctx,
		chromedp.Navigate(domain),
		chromedp.Sleep(4*time.Second),
		chromedp.Location(&finalURL),
	)
	if err != nil {
		return ""
	}
	return finalURL
}

func ensureScheme(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "http://" + url
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

