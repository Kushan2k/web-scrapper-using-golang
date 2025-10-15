package utils

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

/*
@dev function for showning the oprations and geting user inputs
@params l map[string]string, to_find *[]string
@returns void
*/
func Generate_list(l map[string]string, to_find *[]string) {
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
func IsValidUrl(add string) (bool, error) {

	u, err := url.ParseRequestURI(add)
	if err != nil {
		return false, err
	}

	if u.Scheme == "" || u.Host == "" {
		return false, fmt.Errorf("invalid url")
	}
	return true, nil

}
