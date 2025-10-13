package main

import (
	"fmt"
	"strings"
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

	fmt.Print("Enter URL:- ")
	fmt.Scan(&url)

	generate_list(elements,&to_find)

	

	
	
	
	
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
