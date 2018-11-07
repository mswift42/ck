package main

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"testing"
)

var queryurls = []struct {
	st   string
	page string
	want string
}{
	{
		"rotwein",
		"0",
		"https://www.chefkoch.de/rs/s0/rotwein/Rezepte.html#more2",
	},
	{
		"rotwein",
		"60",
		"https://www.chefkoch.de/rs/s60/rotwein/Rezepte.html#more2",
	},
	{
		"bohnen",
		"120",
		"https://www.chefkoch.de/rs/s120/bohnen/Rezepte.html#more2",
	},
}

func TestQueryURL(t *testing.T) {
	for _, i := range queryurls {
		qu := queryUrl(i.st, i.page)
		if qu != i.want {
			t.Errorf("Expected qu to be %s, got: %s", i.want, qu)
		}
	}

}

var bohenrecipes = []struct {
	title      string
	subtitle   string
	url        string
	thumbnail  string
	rating     string
	difficulty string
	preptime   string
}{
	{
		"Grüne Bohnen im Speckmantel",
		"Bohnen waschen und die Spitzen abschneiden. Bohnenkraut, Knoblauch, zerdrückte Pfefferkörner und Salz mit Öl kurz anrösten. 2 Lite...",
		"https://www.chefkoch.de/rezepte/563451154612271/Gruene-Bohnen-im-Speckmantel.html",
		"https://static.chefkoch-cdn.de/rs/bilder/56345/gruene-bohnen-im-speckmantel-1124631-150x150.jpg",
		"4.32",
		"simpel",
		"30 min",
	},
}

func TestNewRecipe(t *testing.T) {
	file, err := ioutil.ReadFile("testhtml/bohnen.html")
	if err != nil {
		t.Error(err)
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(file))
	if err != nil {
		panic(err)
	}
	var results []*Recipe
	doc.Find(".search-list-item").Each(func(i int, sel *goquery.Selection) {
		results = append(results, newRecipe(sel))
	})
	if len(results) != 30 {
		t.Error("Expected length of results to be 30, got: ", len(results))
	}
	if results[0].Title != "Grüne Bohnen im Speckmantel" {
		t.Error("Expected title to be ... got: ", results[0].Title)
	}
}
