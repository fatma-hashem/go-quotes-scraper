package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// Quote struct for JSON
type Quote struct {
	Text   string `json:"text"`
	Author string `json:"author"`
}

func scrapeQuotes() ([]Quote, error) {
	var quotes []Quote

	// Fetch the page
	res, err := http.Get("https://quotes.toscrape.com/")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	// Extract quotes
	doc.Find(".quote").Each(func(i int, s *goquery.Selection) {
		text := s.Find(".text").Text()
		author := s.Find(".author").Text()
		quotes = append(quotes, Quote{Text: text, Author: author})
	})

	return quotes, nil
}

func quotesHandler(w http.ResponseWriter, r *http.Request) {
	quotes, err := scrapeQuotes()
	if err != nil {
		http.Error(w, "Failed to scrape quotes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quotes)
}

func main() {
	http.HandleFunc("/quotes", quotesHandler)
	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

