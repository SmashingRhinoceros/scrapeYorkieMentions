package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Dog struct {
	gorm.Model
	Name             string
	NumberOfMentions int
	BarkMentions     string
}

func newDog(name string) *Dog {
	d := Dog{Name: name}
	return &d
}

var allDogs []Dog

var negativeWords = []string{"no", "never", "don't", "rarely", "seldom", "won't"}

var barkWords = []string{"bark", "yap", "noise", "noisy"}

func getBarks(textBlob string) (int, string) {
	numberOfMentions := 0
	var barkMentions string
	words := strings.Fields(textBlob)
	negatives := false
	for i, word := range words {

		for _, bW := range barkWords {

			if strings.Contains(word, bW) {
				fmt.Println("NEW BARK", word, i)
				for n := i - 1; n >= 0 && n >= i-3 && negatives == false; n-- {
					fmt.Println("i: ", i, words[i], "n: ", n, words[n])

					for _, neg := range negativeWords {
						if words[n] == neg {
							negatives = true
						} else {
						}
					}

				}
				if !negatives {
					numberOfMentions += 1
					barkMentions += word
					barkMentions += "_"
				}

				fmt.Println(word)

			}
		}

	}

	return numberOfMentions, barkMentions
}

func main() {

	/*f, err := os.Create("barkData2.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	*/

	db, err := gorm.Open(sqlite.Open("newBarkData.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	db.AutoMigrate(&Dog{})

	//Link collector
	linkCollector := colly.NewCollector(
		colly.MaxDepth(1))
	linkCollector.SetRequestTimeout(120 * time.Second)

	//Bark collector

	barkCollector := colly.NewCollector()
	barkCollector.SetRequestTimeout(120 * time.Second)

	//Link callbacks
	linkCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Link collector visiting the page ", r.URL)
	})
	linkCollector.OnResponse(func(r *colly.Response) {
		fmt.Println("Link collector got response from", r.Request.URL)
	})
	linkCollector.OnError(func(r *colly.Response, e error) {
		fmt.Println("Link collector got the error", e)
	})

	linkCollector.OnHTML("a.list-item-title", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		d := Dog{Name: e.Text}
		allDogs = append(allDogs, d)

		barkCollector.Visit(link)
	})

	linkCollector.OnScraped(func(r *colly.Response) {

		fmt.Println("Finished", r.Request.URL)
	})

	//Bark callbacks

	barkCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Bark collector visiting the page ", r.URL)

	})
	barkCollector.OnResponse(func(r *colly.Response) {
		fmt.Println("Bark collector got response from", r.Request.URL)

	})
	barkCollector.OnError(func(r *colly.Response, e error) {
		fmt.Println("Bark collector got the error", e)
	})

	barkCollector.OnHTML("div.breeds-single-intro", func(e *colly.HTMLElement) {

		textBlob := e.Text

		numberOfMentions, barkMentions := getBarks(textBlob)

		var nowDog *Dog = &allDogs[len(allDogs)-1]
		nowDog.NumberOfMentions += numberOfMentions
		nowDog.BarkMentions += barkMentions
	})

	barkCollector.OnHTML("div.breed-data-item-content", func(e *colly.HTMLElement) {
		textBlob := e.Text

		numberOfMentions, barkMentions := getBarks(textBlob)

		var nowDog *Dog = &allDogs[len(allDogs)-1]
		nowDog.NumberOfMentions += numberOfMentions
		nowDog.BarkMentions += barkMentions

	})

	barkCollector.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
		fmt.Println(allDogs)
		var nowDog = allDogs[len(allDogs)-1]
		db.Create(&Dog{Name: nowDog.Name, NumberOfMentions: nowDog.NumberOfMentions, BarkMentions: nowDog.BarkMentions})

	})

	linkCollector.Visit("https://dogtime.com/dog-breeds/profiles")

}
