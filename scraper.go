package main

import (
	"github.com/gocolly/colly"
	"fmt"
	"strings"
	"time"

)

func main() {
	linkCollector := colly.NewCollector(
		colly.MaxDepth(1))
	linkCollector.SetRequestTimeout(120*time.Second)


	linkCollector.OnRequest(func(r *colly.Request){
		fmt.Println("Link collector visiting the page ", r.URL)
	})
	linkCollector.OnResponse(func(r *colly.Response){
		fmt.Println("Link collector got response from", r.Request.URL)
	})
	linkCollector.OnError(func(r *colly.Response, e error){
		fmt.Println("Link collector got the error", e)
	})

	barkCollector := colly.NewCollector()
	barkCollector.SetRequestTimeout(120*time.Second)

	//Callbacks
	textBlob := ""
	numberMentions := 0
	barkCollector.OnHTML("div.breeds-single-intro", func(e *colly.HTMLElement) {

		e.ForEach("p", func(_ int, element *colly.HTMLElement) {
			textBlob += e.Text
		})
		barkCollector.OnHTML("ul.breed-data js-accordion item-expandable-container profile-descriptions-list", func(e *colly.HTMLElement) {

			e.ForEach("li.breed-data-item js-accordion-item item-expandable-content", func(_ int, element *colly.HTMLElement) {

			})

		})
		words := strings.Fields(textBlob)

		for _, word := range words {
			if strings.Contains(word, "bark") {
				fmt.Println(word)
				numberMentions += 1
			}
			if strings.Contains(word, "yap") {
				fmt.Println(word)
				numberMentions += 1
			}
		}
		fmt.Println("number mentions", numberMentions)

		barkCollector.OnRequest(func(r *colly.Request) {
			fmt.Println("Bark collector visiting the page ", r.URL)
		})
		barkCollector.OnResponse(func(r *colly.Response) {
			fmt.Println("Bark collector got response from", r.Request.URL)
		})
		barkCollector.OnError(func(r *colly.Response, e error) {
			fmt.Println("Bark collector got the error", e)
		})

		linkCollector.OnHTML("a.list-item-title", func(e *colly.HTMLElement) {
			link := e.Attr("href")
			fmt.Println(link)
			barkCollector.Visit(link)
		})

		linkCollector.Visit("https://dogtime.com/dog-breeds/profiles")

	}
}
