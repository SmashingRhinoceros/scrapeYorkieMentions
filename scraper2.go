package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"strings"
	"time"
)

func getBarkDetails(link string) {

}

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
		barkCollector.Visit(link)

	})

	//Bark callbacks

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
		fmt.Println("Number of mentions is ", numberOfMentions)

	})

	linkCollector.Visit("https://dogtime.com/dog-breeds/profiles")

}
