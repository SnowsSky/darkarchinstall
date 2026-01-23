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

var accounts []src.Accounts

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

	case "acc":
		var acc_opt string
		form := src.Accounts_form(&acc_opt)
		err := form.Run()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		switch acc_opt {
		case "acc_add":
			var account src.Accounts
			form := src.Account_add_form(&account.Username, &account.Password, &account.Sudo)
			err := form.Run()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			accounts = append(accounts, account)

		case "acc_remove":
			var selectedUser string
			fmt.Printf("Account added: %+v\n", accounts)
			form := src.Account_remove_form(accounts, &selectedUser)
			err := form.Run()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			var updated []src.Accounts
			for _, acc := range accounts {
				if acc.Username != selectedUser {
					updated = append(updated, acc)
				}
			}

			accounts = updated
		}

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
