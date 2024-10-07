package main

import (
	"bufio"
	"flag"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sort"
)

// prints the version message
const version = "0.0.2"

func printVersion() {
	fmt.Printf("Current uforall version %s\n", version)
}

// Prints the Colorful banner
func printBanner() {
	banner := `
          ____                     __ __
  __  __ / __/____   _____ ____ _ / // /
 / / / // /_ / __ \ / ___// __  // // / 
/ /_/ // __// /_/ // /   / /_/ // // /  
\__,_//_/   \____//_/    \__,_//_//_/                                        
`
fmt.Printf("%s\n%50s\n\n", banner, "Current uforall version "+version)

}

func getArchiveUrls(domain string) {
	url := fmt.Sprintf("http://web.archive.org/cdx/search/cdx?url=*.%s/*&output=text&fl=original&collapse=urlkey", domain)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "HTTP error:", resp.Status)
		return
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading response body:", err)
	}
}


type OtxUrlData struct {
	Url string `json:"url"`
}

type OtxResponseData struct {
	FullSize int        `json:"full_size"`
	UrlList  []OtxUrlData  `json:"url_list"`
}

func getOtxUrls(domain string) {
	url := fmt.Sprintf("https://otx.alienvault.com/api/v1/indicators/domain/%s/url_list", domain)
	client := &http.Client{}

	// Initial request to get total pages
	resp, err := client.Get(url + "?limit=500&page=1")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "HTTP error:", resp.Status)
		return
	}

	var data OtxResponseData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Fprintln(os.Stderr, "JSON decode error:", err)
		return
	}

	totalUrls := data.FullSize
	totalPages := (totalUrls + 499) / 500 // Using ceiling to calculate total pages

	for i := 1; i <= totalPages; i++ {
		resp, err := client.Get(fmt.Sprintf("%s?limit=500&page=%d", url, i))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Fprintln(os.Stderr, "HTTP error:", resp.Status)
			return
		}

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			fmt.Fprintln(os.Stderr, "JSON decode error:", err)
			return
		}

		for _, urlData := range data.UrlList {
			fmt.Println(urlData.Url)
		}
	}
}


type UrlscanResult struct {
	Page struct {
		URL string `json:"url"`
	} `json:"page"`
	Task struct {
		URL string `json:"url"`
	} `json:"task"`
}

type UrlscanResponseData struct {
	Results []UrlscanResult `json:"results"`
}

func getUrlscanUrls(domain string) {
	url := fmt.Sprintf("https://urlscan.io/api/v1/search/?q=domain:%s&size=10000", domain)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "HTTP error:", resp.Status)
		return
	}

	var data UrlscanResponseData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Fprintln(os.Stderr, "JSON decode error:", err)
		return
	}

	// Extract URLs
	var urls []string
	for _, result := range data.Results {
		if result.Page.URL != "" {
			urls = append(urls, result.Page.URL)
		}
		if result.Task.URL != "" {
			urls = append(urls, result.Task.URL)
		}
	}

	// Remove duplicates and sort URLs
	urlSet := make(map[string]struct{})
	for _, url := range urls {
		urlSet[url] = struct{}{}
	}

	var uniqueUrls []string
	for url := range urlSet {
		uniqueUrls = append(uniqueUrls, url)
	}
	sort.Strings(uniqueUrls)

	// Print URLs
	for _, url := range uniqueUrls {
		fmt.Println(url)
	}
}


// Fetch URLs from Common Crawl
type CdxApi struct {
	CDXAPI string `json:"cdx-api"`
}

func getCommonCrawlUrls(domain string) {
	resp, err := http.Get("https://index.commoncrawl.org/collinfo.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "RequestException: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading response body: %v\n", err)
		return
	}

	var data []CdxApi
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Fprintf(os.Stderr, "JSON decode error: %v\n", err)
		return
	}

	// Process each URL sequentially
	for _, item := range data {
		apiUrl := item.CDXAPI
		urls := fetchCommonCrawlUrls(apiUrl, domain)
		for _, url := range urls {
			fmt.Println(url)
		}
	}
}

func fetchCommonCrawlUrls(apiUrl, domain string) []string {
	var urls []string
	param := fmt.Sprintf("?url=*.%s&fl=url&output=json", domain)
	indexUrl := apiUrl + param

	resp, err := http.Get(indexUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ConnectionError: %v\n", err)
		return urls
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading response body: %v\n", err)
		return urls
	}

	urlPattern := regexp.MustCompile(`"url": "([^"]+)"`)
	matches := urlPattern.FindAllStringSubmatch(string(body), -1)
	for _, match := range matches {
		if len(match) > 1 {
			urls = append(urls, match[1])
		}
	}

	return urls
}


func main() {
	toolFlag := flag.String("t", "all", "Comma-separated list of tools to run: 'otx', 'archive', 'urlscan', 'commoncrawl', or 'all'")
	version := flag.Bool("version", false, "Print the version of the tool and exit.")
	silent := flag.Bool("silent", false, "silent mode.")
	flag.Parse()

	// Print version and exit if -version flag is provided
	if *version {
		printBanner()
		printVersion()
		return
	}

	// Don't Print banner if -silnet flag is provided
	if !*silent {
		printBanner()
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())
		
		// Handle each tool based on the flag
		if *toolFlag == "all" || strings.Contains(*toolFlag, "archive") {
			getArchiveUrls(domain)
		}
		if *toolFlag == "all" || strings.Contains(*toolFlag, "otx") {
			getOtxUrls(domain)
		}
		if *toolFlag == "all" || strings.Contains(*toolFlag, "urlscan") {
			getUrlscanUrls(domain)
		}
		if *toolFlag == "all" || strings.Contains(*toolFlag, "commoncrawl") {
			getCommonCrawlUrls(domain)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}
}
