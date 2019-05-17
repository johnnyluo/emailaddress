package main

import (
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/fatih/color"
)

func main() {
	exp, err := ioutil.ReadFile("../../data/regex_expression.dat")
	if nil != err {
		panic(err)
	}
	fmt.Println("regex:" + string(exp))
	re := regexp.MustCompile(string(exp))
	for {
		var input string
		n, err := fmt.Scan(&input)
		if nil != err {
			panic(err)
		}
		if n == 0 { // let's exist
			return
		}
		if !re.MatchString(input) {
			d := color.New(color.FgRed)
			d.Printf("%s is invalid email address\n", input)
		} else {
			d := color.New(color.FgGreen)
			d.Printf("%s is a valid email address\n", input)
		}
	}
}
