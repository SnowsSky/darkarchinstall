package main

import (
	"darkarchinstall/forms"
	"fmt"
)

var opt string

func main() {
	opt = ""
	for {
		form := forms.Main_form(&opt)
		err := form.Run()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		forms.Options_check(opt)
	}
}
