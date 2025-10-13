package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)


func main()  {

	var url string
	var elements = map[string]string{
		"h1":"main headdings",
		"img":"images",
		"title":"title of the page",
		"p":"all the paragraps",
	}
	var to_find = make(map[string]string)

	var results= make(chan map[string]string,10)
	var wg sync.WaitGroup

	

	fmt.Print("Enter URL:- ")
	fmt.Scan(&url)

	suc,err:=isValidUrl(url)

	if !suc {
		log.Fatalln(err.Error())
	}

	generate_list(elements,&to_find)


	c:=colly.NewCollector()

	c.OnRequest(func(r *colly.Request){
		fmt.Println("Visiting",r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("👍 Done Loading the URL")
	})

	for key,val := range to_find {
		wg.Add(1)

		func(k, v string) { 
			c.OnHTML(k, func(html *colly.HTMLElement) {
				results <- map[string]string{v: html.Text}
			})
			wg.Done()
		}(key, val)

	}

	go func() {
		wg.Wait()
		close(results)
	}()
	c.Visit(url)

	for res := range results {
		for k, v := range res {
			fmt.Printf("%s: %s\n", k, v)
		}
	}


	
}

func generate_list(l map[string]string,to_find *map[string]string){
	fmt.Println("Please enter the number what you want to find(seperate by ,):-")
	i:=1
	for key := range l {
		fmt.Printf("%d - %s\n",i,key)
		i++
	}
	var input string
	fmt.Scan(&input)
	fmt.Println(input)

	var keys = strings.Split(input, ",")
	for _,v := range keys {
		(*to_find)[v] = l[v]
	}


}

func isValidUrl(add string) (bool,error){

	u,err:=url.ParseRequestURI(add)
	if err!=nil{
		return false,err
	}

	if u.Scheme=="" || u.Host==""{
		return false,fmt.Errorf("invalid url")
	}
	return true,nil
	
}
