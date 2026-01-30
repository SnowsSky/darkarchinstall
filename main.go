package main

import (
	"darkarchinstall/forms"
	"darkarchinstall/network"
)

var opt string

func main() {
	network.NetworkCheck()
	for {
		forms.MainForm(&opt).Run()
		forms.Options_check(opt)
	}
}
