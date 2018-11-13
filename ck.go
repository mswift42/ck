package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"regexp"
	"strings"
)

type Recipe struct {
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle"`
	Url        string `json:"url"`
	Thumbnail  string `json:"thumbnai"`
	Rating     string `json:"rating"`
	Difficulty string `json:"difficulty"`
	Preptime   string `json:"preptime"`
}

type RecipeDetail struct {
	Recipe      *Recipe  `json:"recipe"`
	Thumbnail   string   `json:"thumbnail"`
	Ingredients []string `json:"ingredients"`
	Method      string   `json:"method"`
}

const CKPrefix = "https://www.chefkoch.de"

func queryUrl(searchterm string, page string) string {
	return "https://www.chefkoch.de/rs/s" + page + "/" + searchterm + "/Rezepte.html#more2"
}

type RecipesSelection struct {
	sel *goquery.Selection
}

func (rs *RecipesSelection) title() string {
	return rs.sel.Find(".search-list-item-title").Text()
}

func (rs *RecipesSelection) subtitle() string {
	subtitle := rs.sel.Find(".search-list-item-subtitle").Text()
	trimmed := strings.Trim(subtitle, " \n")
	return strings.Replace(trimmed, "\n", " ", -1)
}

func (rs *RecipesSelection) url() string {
	return CKPrefix + rs.sel.Find(".search-list-item > a").AttrOr("href", "")
}

func (rs *RecipesSelection) thumbnail() string {
	return rs.sel.Find("picture > img").AttrOr("srcset", "")
}

func (rs *RecipesSelection) rating() string {
	rating := rs.sel.Find(".search-list-item-uservotes-stars").AttrOr("title", "")
	digregex := regexp.MustCompile(`\d\.\d*`)
	return digregex.FindString(rating)

}

func (rs *RecipesSelection) difficulty() string {
	return rs.sel.Find(".search-list-item-difficulty").Text()
}

func (rs *RecipesSelection) preptime() string {
	return rs.sel.Find(".search-list-item-preptime").Text()
}

func NewRecipe(sel *goquery.Selection) *Recipe {
	rs := &RecipesSelection{sel}
	return &Recipe{rs.title(), rs.subtitle(),
		rs.url(), rs.thumbnail(), rs.rating(), rs.difficulty(),
		rs.preptime()}
}

func allRecipes(doc *goquery.Document) []*Recipe {
	var results []*Recipe
	doc.Find(".search-list-item").Each(func(i int, s *goquery.Selection) {
		results = append(results, NewRecipe(s))
	})
	return results
}

func recipesToJson(recipes []*Recipe) ([]byte, error) {
	return json.Marshal(recipes)
}

func main() {
	res, err := http.Get("https://www.chefkoch.de/rs/s0/bohnen/Rezepte.html#more2")
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(allRecipes(doc))
}
