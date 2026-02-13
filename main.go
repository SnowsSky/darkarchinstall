package main

import (
	"darkarchinstall/forms"
	"darkarchinstall/network"
	"fmt"
)

var opt string
var version string = "1.0.3"

func main() {
	err := forms.CheckRoot()
	if err != nil {
		fmt.Println(err)
		return
	}
	network.NetworkCheck()
	for forms.Ininstaller {
		forms.MainForm(&opt, version).Run()
		forms.Options_check(opt, version)
	}
}
