package services

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

func scrape_devpost_projects() {
	url := "https://hackumass-xii.devpost.com/project-gallery?page=1"

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (compatible; CollyScraper/1.0)"),
		colly.AllowedDomains("hackumass-xii.devpost.com", "devpost.com"),
	)

	// Rate limiting for safety
	_ = c.Limit(&colly.LimitRule{
		DomainGlob:  "*devpost.*",
		Parallelism: 2,
		Delay:       800 * time.Millisecond,
		RandomDelay: 400 * time.Millisecond,
	})

	var titles []string

	// Each project name is inside <div class="software-entry-name"><h5>...</h5></div>
	c.OnHTML("div.software-entry-name h5", func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.Text)
		if title != "" {
			fmt.Println(title)
			titles = append(titles, title)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Printf("Visiting %s\n", r.URL.String())
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Request error: %s -> %v\n", r.Request.URL, err)
	})

	if err := c.Visit(url); err != nil {
		log.Fatalf("Visit error: %v", err)
	}
	c.Wait()

	// Write to file
	f, err := os.Create("titles_page1.txt")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer f.Close()

	for _, t := range titles {
		f.WriteString(t + "\n")
	}
	log.Printf("✅ Wrote %d project titles to titles_page1.txt", len(titles))
}
