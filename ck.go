package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
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
	Title       string              `json:"title"`
	Rating      string              `json:"rating"`
	Difficulty  string              `json:"difficulty"`
	Preptime    string              `json:"preptime"`
	Cookingtime string              `json:"cookingtime"`
	Thumbnail   string              `json:"thumbnail"`
	Ingredients []*RecipeIngredient `json:"ingredients"`
	Method      string              `json:"method"`
}

type RecipeIngredient struct {
	Amount     string `json:"amount"`
	Ingredient string `json:"ingredient"`
}

type RecipeDetailDocument struct {
	doc *goquery.Document
}

func (rdd *RecipeDetailDocument) newRecipeDetail() *RecipeDetail {
	title := rdd.title()
	ingredients := rdd.ingredients()
	method := rdd.method()
	rating := rdd.rating()
	difficulty := rdd.difficulty()
	preptime := rdd.preptime()
	cookingtime := rdd.cookingtime()
	thumbnail := rdd.thumbnail()
	return &RecipeDetail{title, rating, difficulty,
		preptime, cookingtime, thumbnail, ingredients, method}
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

func (rdd *RecipeDetailDocument) title() string {
	return rdd.doc.Find(".page-title").Text()
}

func (rdd *RecipeDetailDocument) rating() string {
	rat := rdd.doc.Find(".rating__average-rating").Text()
	rat = strings.Replace(rat, "Ø", "", 1)
	return strings.Replace(rat, ",", ".", 1)
}

func (rdd *RecipeDetailDocument) difficulty() string {
	return strings.Trim(rdd.prepinfo()[2], " \n")
}

func (rdd *RecipeDetailDocument) preptime() string {
	return strings.Trim(rdd.prepinfo()[0], " \n")
}

func (rdd *RecipeDetailDocument) cookingtime() string {
	return strings.Trim(rdd.prepinfo()[1], " \n")
}

func (rdd *RecipeDetailDocument) prepinfo() []string {
	prep := rdd.doc.Find("#preparation-info").Text()
	prep = strings.Replace(prep, "Arbeitszeit: ", "", 1)
	prep = strings.Replace(prep, "Koch-/Backzeit: ", "", 1)
	prep = strings.Replace(prep, "Schwierigkeitsgrad: ", "", 1)
	return strings.Split(prep, "/")
}

func (rdd *RecipeDetailDocument) thumbnail() string {
	return rdd.doc.Find(".slideshow-image").AttrOr("src", "")
}

func (rdd *RecipeDetailDocument) ingredients() []*RecipeIngredient {
	ingtable := rdd.doc.Find(".incredients>tbody>tr")
	var ingredients []*RecipeIngredient
	ingtable.Each(func(i int, s *goquery.Selection) {
		amount := strings.Trim(s.Find(".amount").Text(), " \n")
		ing := strings.Trim(s.Find("td:nth-child(2)").Text(), " \n")
		ingredients = append(ingredients, &RecipeIngredient{amount, ing})
	})
	return ingredients
}

func (rdd *RecipeDetailDocument) method() string {
	text := rdd.doc.Find("#rezept-zubereitung").Text()
	return strings.Trim(text, " \n")
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

func recipeDetailToJson(rd *RecipeDetail) ([]byte, error) {
	return json.Marshal(rd)
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
