package main

import (
	"fmt"
	"github.com/gocolly/colly"
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

func main() {

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
		fmt.Println(link)
		d := dog{name: e.Text}
		allDogs = append(allDogs, d)
		fmt.Println(allDogs)
		barkCollector.Visit(link)
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
		textBlob := ""
		numberOfMentions := 0
		e.ForEach("p", func(_ int, element *colly.HTMLElement) {
			textBlob += element.Text
		})

		words := strings.Fields(textBlob)

		for _, word := range words {
			if strings.Contains(word, "bark") {
				fmt.Println(word)
				numberOfMentions += 1
			}
			if strings.Contains(word, "yap") {
				fmt.Println(word)
				numberOfMentions += 1
			}

		}
		fmt.Println(numberOfMentions)
		var nowDog *dog = &allDogs[len(allDogs)-1]
		nowDog.numberOfMentions += numberOfMentions

		fmt.Println("NOWDOG", nowDog)
	})

	barkCollector.OnHTML("ul.breed-data js-accordion item-expandable-container profile-descriptions-list", func(e *colly.HTMLElement) {
		fmt.Println("HELLO", e.Name)

		textBlob := ""
		numberOfMentions := 0
		e.ForEach("li.breed-data-item js-accordion-item item-expandable-content", func(_ int, li *colly.HTMLElement) {
			li.ForEach("div.breed-data-item-content js-breed-data-section", func(_ int, div *colly.HTMLElement) {
				div.ForEach("p", func(_ int, p *colly.HTMLElement) {
					textBlob += div.Text
				})
			})
		})

		words := strings.Fields(textBlob)

		for _, word := range words {
			if strings.Contains(word, "bark") {
				fmt.Println(word)
				numberOfMentions += 1
			}
			if strings.Contains(word, "yap") {
				fmt.Println(word)
				numberOfMentions += 1
			}
		}

		var nowDog *dog = &allDogs[len(allDogs)-1]
		nowDog.numberOfMentions += numberOfMentions
		fmt.Println("NOWDOG", nowDog)

	})

	linkCollector.Visit("https://dogtime.com/dog-breeds/profiles")

}
