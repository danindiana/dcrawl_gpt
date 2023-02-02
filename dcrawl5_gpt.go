//Write multi-threaded web crawler for randomly gathering lists of unique domain names in golang. The program takes one site URL as input and detects all <a href=...> links in the site's body. Each found link is put into the queue. Successively, each queued link is crawled in the same way, branching out to more URLs found in links on each site's body.

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

// global variables
var domainList = make(map[string]int)
var mutex = &sync.Mutex{}
var wg sync.WaitGroup
var queue = make(chan string)
var visited = make(map[string]bool)
var maxDepth = 3
var maxGoroutines = 10

func main() {
	start := time.Now()
	wg.Add(1)
	wg.Wait()
	fmt.Println(domainList)
	fmt.Printf("Time taken: %s	", time.Since(start))
}

func crawl(url string, depth int) {
	defer wg.Done()
	if depth > maxDepth {
		return
	}
	visited[url] = true
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	links := extractLinks(string(body))
	for _, link := range links {
		if !visited[link] {
			wg.Add(1)
			go crawl(link, depth+1)
		}
	}
	domain := "https://www.algolia.com" // dummy domain
	if domain != "" {
		mutex.Lock()
		if _, ok := domainList[domain]; !ok {
			domainList[domain] = 1
		}
		mutex.Unlock()
	}
}

// create a URLqueue
type URLqueue struct {
	list []string
	lock sync.Mutex
}

// add a URL to the queue
func (u *URLqueue) Push(url string) {
	u.lock.Lock()
	u.list = append(u.list, url)
	u.lock.Unlock()
}

// remove a URL from the queue
func (u *URLqueue) Pop() string {
	u.lock.Lock()
	url := u.list[0]
	u.list = u.list[1:]
	u.lock.Unlock()
	return url
}

// get the length of the queue
func (u *URLqueue) Length() int {
	return len(u.list)
}

// fetch the HTML body from the URL
func getHTML(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	// read the response body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// extract all the links from the HTML body
func extractLinks(html string) []string {
	links := []string{}
	// parse the HTML and extract all links
	// ...
	return links
}

// main crawler function
func crawler(url string, queue *URLqueue) {
	// fetch the HTML from the URL
	html, err := getHTML(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	// extract all links from the HTML
	links := extractLinks(html)
	for _, link := range links {
		// add the link to the queue
		queue.Push(link)
	}

	// extract the domain from the URL
	domain := "example.com" // dummy domain
	if domain != "" {
	} // check if the domain is not empty
	{
		// add the domain to the list
		mutex.Lock()
		if _, ok := domainList[domain]; !ok {
			domainList[domain] = 1
		}
		mutex.Unlock()
	}

	// crawl the next URL
	if queue.Length() > 0 {
		nextURL := queue.Pop()
		crawler(nextURL, queue)
	}
}

func Wait() { // create a wait group
	start := time.Now()

	// create a URL queue
	queue := &URLqueue{}

	// add the initial URL to the queue
	queue.Push("http://example.com")

	// spawn multiple goroutines to crawl the URLs
	for i := 0; i < 10; i++ {
		go crawler(queue.Pop(), queue)
	}

	// wait until all goroutines are done
	x := sync.WaitGroup{}
	x.Add(1) // add a counter
	x.Wait() // wait until the counter is zero

	// print the list of unique domains
	fmt.Println(domainList)

	fmt.Printf("Time taken: %s\n", time.Since(start))
}
