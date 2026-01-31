package network

import (
	"fmt"
	"net/http"
	"os"

	"github.com/charmbracelet/huh/spinner"
)

func NetworkCheck() {
	action := func() {
		_, err := http.Get("https://github.com/darkarchlinux")
		if err != nil {
			fmt.Println("You need to be connected to the internet to install Dark Arch.")
			os.Exit(1)
		}
	}
	if err := spinner.New().Title("Checking network...").Action(action).Run(); err != nil {
		fmt.Println("Failed:", err)
		os.Exit(1)

	}

}
