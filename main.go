package main

import (
	"darkarchinstall/src"
	"darkarchinstall/src/checks"
	"fmt"
)

var opt string

func main() {
	opt = ""
	for {
		form := src.Main_form(&opt)
		err := form.Run()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		checks.Options_check(opt)
	}
}
