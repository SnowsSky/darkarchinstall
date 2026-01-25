package main

import (
	"darkarchinstall/forms"
)

var opt string

func main() {
	for {
		forms.Main_form(&opt).Run()
		forms.Options_check(opt)
	}
}
