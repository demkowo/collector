package service

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/demkowo/goquery/model"
)

type CollectorRepository interface {
	GatherAllLinks(links []*model.Links) error
}

type CollectorService interface {
	GatherAllLinks([]string, ...string) error
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

func (s collectorService) GatherAllLinks(urls []string, limits ...string) error {
	for _, url := range urls {

		// Create goquery doc from url
		doc := createGoqueryDoc(url)

		linksSlice := []*model.Links{}

		// Find all <a> elements
		doc.Find("a").Each(func(index int, element *goquery.Selection) {
			// Get the href attribute value from each <a> element
			link, _ := element.Attr("href")

			// Check limits
			if len(limits) == 0 || checkLimits(link, limits) {
				// If link fullfull conditions, store it in the map
				//links[link] = true
				links := &model.Links{
					Url:     link,
					Details: false,
				}
				linksSlice = append(linksSlice, links)
			}
		})

		if err := s.repo.GatherAllLinks(linksSlice); err != nil {
			return err
		}
	}
	return nil
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
