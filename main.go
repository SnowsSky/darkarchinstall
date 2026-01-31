package main

import (
	"darkarchinstall/forms"
	"darkarchinstall/network"
	"fmt"
)

var opt string

func main() {
	err := forms.CheckRoot()
	if err != nil {
		fmt.Println(err)
		return
	}
	network.NetworkCheck()
	for forms.Ininstaller {
		forms.MainForm(&opt).Run()
		forms.Options_check(opt)
	}
}
