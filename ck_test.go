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
	for ind, i := range bohenrecipes {
		if results[ind].Title != i.title {
			t.Errorf("Expected title to be %q, got: %q", i.title, results[ind].Title)
		}
		if results[ind].Subtitle != i.subtitle {
			t.Errorf("Expected subtitle to be %q, got: %q", i.subtitle,
				results[ind].Subtitle)
		}
		if results[ind].Url != i.url {
			t.Errorf("Expected url to be %q, got: %q", i.url,
				results[ind].Url)
		}
		if results[ind].Thumbnail != i.thumbnail {
			t.Errorf("Expected thumbnail to be %q, got: %q", i.thumbnail,
				results[ind].Thumbnail)
		}
		if results[ind].Rating != i.rating {
			t.Errorf("Expected rating to be %q, got %q", i.rating,
				results[ind].Rating)
		}
		if results[ind].Difficulty != i.difficulty {
			t.Errorf("Expected difficulty to be %q, got %q", i.difficulty,
				results[ind].Difficulty)
		}
		if results[ind].Preptime != i.preptime {
			t.Errorf("Expected preptime to be %q, got %q", i.preptime,
				results[ind].Preptime)
		}

	}
}
