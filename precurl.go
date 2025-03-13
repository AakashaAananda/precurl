package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jaytaylor/html2text"
)

func fetchPage(targetURL string) string {
	parsedURL, err := url.ParseRequestURI(targetURL)
	if err != nil {
		fmt.Println(" XXX Invalid URL format")
		os.Exit(1)
	}

	if parsedURL.Scheme == "" {
		targetURL = "https://" + targetURL
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		fmt.Println(" XXX Error creating request:", err)
		os.Exit(1)
	}

	req.Header.Set("User-Agent", "precurl/1.0 (https://github.com/AakashaAananda/precurl)")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(" XXX Error fetching the URL:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf(" XXX  HTTP Error %d: %s\n", resp.StatusCode, http.StatusText(resp.StatusCode))
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(" XXX Error reading response body:", err)
		os.Exit(1)
	}

	return string(body)
}

func renderHTMLAsText(html string) string {
	text, err := html2text.FromString(html, html2text.Options{
		PrettyTables: true,
		OmitLinks:    false,
	})
	if err != nil {
		fmt.Println(" XXX Error parsing HTML:", err)
		os.Exit(1)
	}

	formattedText := strings.ReplaceAll(text, "\n\n", "\n")
	return formattedText
}

func printWithColor(text string) {
	c := color.New(color.FgCyan, color.Bold)
	c.Println("--- Webpage Preview ---")
	fmt.Println(text)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: precurl <URL>")
		os.Exit(1)
	}
	targetURL := os.Args[1]

	htmlContent := fetchPage(targetURL)
	readableText := renderHTMLAsText(htmlContent)

	printWithColor(readableText)
}
