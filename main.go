package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/gocolly/colly"
	"github.com/webscraper/utils"
)

func main() {

	var max_depth int

	if os.Args[1] == "" {
		int_max, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatalln("Invalid max depth:", err)
		}

		

		if int_max <= 0 {
			max_depth = 1
		} else {

			max_depth = int_max
		}

	}

	var url string
	var elements = map[string]string{
		"h1":    "main headdings",
		"h2":    "sub headings",
		"img":   "images",
		"title": "title of the page",
		"p":     "all the paragraps",
	}
	var to_find = []string{}

	var results = make(chan map[string]string, 10)
	defer close(results)
	var wg = new(sync.WaitGroup)

	fmt.Print("Enter URL:- ")
	fmt.Scan(&url)

	suc, err :=utils.IsValidUrl(url)

	if !suc {
		log.Fatalln(err.Error())
	}

	utils.Generate_list(elements, &to_find)

	c := colly.NewCollector(
		colly.MaxDepth(max_depth),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("👍 Done Loading the URL")
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	for _, val := range to_find {
		wg.Add(1)

		go func(v string) {
			defer wg.Done()
			c.OnHTML(v, func(html *colly.HTMLElement) {

				results <- map[string]string{v: html.Text}
			})

		}(val)

	}

	c.Visit(url)
	wg.Wait()

	for res := range results {
		for k, v := range res {
			fmt.Printf("%s : %s\n", k, v)
		}
	}

}
