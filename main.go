package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/gocolly/colly"
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

	suc, err := isValidUrl(url)

	if !suc {
		log.Fatalln(err.Error())
	}

	generate_list(elements, &to_find)

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

/*
@dev function for showning the oprations and geting user inputs
@params l map[string]string, to_find *[]string
@returns void
*/
func generate_list(l map[string]string, to_find *[]string) {
	fmt.Println("Please enter the number what you want to find(seperate by ,):-")
	i := 0

	keys := make([]string, 0, 100)

	for k, val := range l {
		fmt.Printf("%d - %s\n", i+1, val)
		keys = append(keys, k)
		i++
	}

	var input string
	fmt.Scan(&input)

	var splited = strings.Split(input, ",")

	var intSplited = make([]int, 0, len(splited))

	for _, s := range splited {
		num, err := strconv.Atoi(s)
		if err != nil {
			fmt.Println("Error converting:", s, err)
			continue
		}
		intSplited = append(intSplited, num)
	}

	for _, val := range intSplited {
		*to_find = append(*to_find, keys[val-1])
	}

}

/*
@dev function for validating a url
@params  add string
@returns  bool,error  true is success and error will not be nill if an error
*/
func isValidUrl(add string) (bool, error) {

	u, err := url.ParseRequestURI(add)
	if err != nil {
		return false, err
	}

	if u.Scheme == "" || u.Host == "" {
		return false, fmt.Errorf("invalid url")
	}
	return true, nil

}
