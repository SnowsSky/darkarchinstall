package main

import (
	"darkarchinstall/forms"
	"darkarchinstall/network"
)

var opt string

func main() {
	network.NetworkCheck()
	for {
		forms.Main_form(&opt).Run()
		forms.Options_check(opt)
	}
}
