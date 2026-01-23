package main

import (
	"darkarchinstall/src"
	"fmt"
)

var (
	opt        string
	confirm    bool
	hostname   string = "darkarch"
	rootpasswd string
)

// Forms

func options_check() {
	//checks
	switch opt {
	case "hostname":
		form := src.Hostname_form(&hostname)
		err := form.Run()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Hostname successfully changed.")
	case "rootpasswd":
		form := src.Root_passwd(&rootpasswd)
		err := form.Run()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Root password sucessfully changed.")

	case "cancel":
		var text string = "You want to exit installation ?"
		form := src.Sure_form(&confirm, &text)
		err := form.Run()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if !confirm {
			main()
			return
		}
		fmt.Println("Installation aborted by user.")
		return
	case "acc":
		fmt.Println("nob")
	}

	main()

}
func main() {
	opt = ""
	form := src.Main_form(&opt)
	err := form.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	options_check()
}
