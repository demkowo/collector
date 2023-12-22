package service

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type CollectorRepository interface {
	FindAllLinks()
}

type CollectorService interface {
	FindAllLinks([]string, ...string)
}

type collectorService struct {
	repo CollectorRepository
}

func New(repository CollectorRepository) CollectorService {
	log.Println("--- service/New() ---")

	return &collectorService{
		repo: repository,
	}
}

func (s collectorService) FindAllLinks(urls []string, limits ...string) {
	for _, url := range urls {

		// Create goquery doc from url
		doc := createGoqueryDoc(url)

		// Slice to store links matching limits
		uniqueLinks := make(map[string]bool)

		// Find all <a> elements
		doc.Find("a").Each(func(index int, element *goquery.Selection) {
			// Get the href attribute value from each <a> element
			link, _ := element.Attr("href")

			// Check limits
			if len(limits) == 0 || checkLimits(link, limits) {
				// If link fullfull conditions, store it in the map
				uniqueLinks[link] = true
			}
		})

		// Convert the unique links map to a slice
		filteredLinks := make([]string, 0, len(uniqueLinks))
		for link := range uniqueLinks {
			filteredLinks = append(filteredLinks, link)
		}

		fmt.Println(len(uniqueLinks))
		// Print or use the filtered links as needed
		for _, link := range filteredLinks {
			fmt.Println("Link:", link)
		}
	}
}

func createGoqueryDoc(url string) *goquery.Document {
	// Make an HTTP GET request to fetch the HTML content
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Parse the HTML content to goquery doc
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

func checkLimits(link string, limits []string) bool {
	for _, limit := range limits {
		if !strings.Contains(link, limit) {
			return false
		}
	}
	return true
}
