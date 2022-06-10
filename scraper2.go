package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
	"strings"
	"time"
)

type dog struct {
	name             string
	numberOfMentions int
}

func newDog(name string) *dog {
	d := dog{name: name}
	return &d
}

var allDogs []dog

var negativeWords = []string{"no", "never", "don't", "rarely", "seldom", "won't"}

func countBarks(textBlob string) int {
	numberOfMentions := 0
	words := strings.Fields(textBlob)
	negatives := false
	for i, word := range words {

		if strings.Contains(word, "bark") || strings.Contains(word, "yap") {
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
			}

			fmt.Println(word)

		}

	}

	return numberOfMentions
}

func main() {

	f, err := os.Create("barkData.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

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
		d := dog{name: e.Text}
		allDogs = append(allDogs, d)

		barkCollector.Visit(link)
	})

	linkCollector.OnScraped(func(r *colly.Response) {
		fmt.Fprintln(f, "All done!")
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

		numberOfMentions := countBarks(textBlob)

		var nowDog *dog = &allDogs[len(allDogs)-1]
		nowDog.numberOfMentions += numberOfMentions
	})

	barkCollector.OnHTML("div.breed-data-item-content", func(e *colly.HTMLElement) {
		textBlob := e.Text

		numberOfMentions := countBarks(textBlob)

		var nowDog *dog = &allDogs[len(allDogs)-1]
		nowDog.numberOfMentions += numberOfMentions

	})

	barkCollector.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
		fmt.Println(allDogs)
		fmt.Fprintln(f, allDogs[len(allDogs)-1])

	})

	linkCollector.Visit("https://dogtime.com/dog-breeds/profiles")

}
