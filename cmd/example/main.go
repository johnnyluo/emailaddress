package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/johnnyluo/emailaddress"
)

func main() {
	fmt.Println(`
##########################################################
please type in email address you would like to validate....
##########################################################
`)
	reader := bufio.NewReader(os.Stdin)
	for {
		buf, _, err := reader.ReadLine()
		if nil != err {
			panic(err)
		}
		input := string(buf)
		b, err := emailaddress.Validate(input)
		if nil != err {
			fmt.Println(err)
		}
		if b {
			d := color.New(color.FgGreen)
			d.Printf("%s is a valid email address\n", input)
		} else {
			d := color.New(color.FgRed)
			d.Printf("%s is invalid email address\n", input)
		}
	}
}
