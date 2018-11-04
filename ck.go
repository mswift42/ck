package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

type Recipe struct {
	title      string
	subtitle   string
	url        string
	thumbnail  string
	rating     string
	difficulty string
	preptime   string
}

func main() {
	doc, err := goquery.NewDocument("https://www.chefkoch.de/rs/s0/bohnen/Rezepte.html#more2")
	if err != nil {
		fmt.Println(err)
	}
	doc.Find(".search-list-item-title").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())
	})

}
