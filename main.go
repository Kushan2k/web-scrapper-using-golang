package main

import (
	"encoding/json"
	"fmt"
	"io"
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
	var elements = map[string]string{}
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

	file, err := os.Open("./tags.json")

	if err != nil {
		log.Fatalln("Error opening tags.json:", err)
	}
	defer file.Close()

	bytes,err:=io.ReadAll(file)

	if err != nil {
		log.Fatalln("Error reading tags.json:", err)
	}
	err = json.Unmarshal(bytes, &elements)
	if err != nil {
		log.Fatalln("Error unmarshalling tags.json:", err)
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

				switch v {
				case "img":
					results <- map[string]string{v: html.Attr("src")}
				case "a":
					results <- map[string]string{v: html.Attr("href")}
				case "title":
					results <- map[string]string{v: html.Text}

				default:
					results <- map[string]string{v: html.Text}

				}
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
