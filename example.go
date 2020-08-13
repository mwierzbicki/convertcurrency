package main

import (
	"fmt"

	cc "github.com/mwierzbicki/convertcurrency/convectcurrency"
)

func main() {
	val, err := cc.ConvertCurrency(14.50, "USD", "CHF", "2020-06-08")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(val)
	}
}
