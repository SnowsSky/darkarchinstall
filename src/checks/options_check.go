package checks

import (
	"bufio"
	"darkarchinstall/src"
	"darkarchinstall/src/types"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var accounts []types.Accounts

var (
	confirm      bool
	hostname     string = "darkarch"
	rootpasswd   string
	selectedUser string
	selected     string
	locales      []string
	keymap       string
	timezone     string
)

func Options_check(opt string) {
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
			var account types.Accounts
			form := src.Account_add_form(&account.Username, &account.Password, &account.Sudo)
			err := form.Run()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			accounts = append(accounts, account)

		case "acc_remove":
			fmt.Printf("Account added: %+v\n", accounts)
			form := src.Account_remove_form(accounts, &selectedUser)
			err := form.Run()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			var updated []types.Accounts
			for _, acc := range accounts {
				if acc.Username != selectedUser {
					updated = append(updated, acc)
				}
			}

			accounts = updated
		case "back":
			return
		}

	case "locales":
		file, err := os.Open("/etc/locale.gen")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer file.Close()
		lineCount := 0
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			// Skip the first 15 lines (header information)
			lineCount++
			if lineCount <= 15 {
				continue
			}

			line := scanner.Text()
			line = strings.ReplaceAll(line, "#", "")
			locales = append(locales, line)
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}

		form := src.Locales_form(&locales, &selected)
		err = form.Run()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	case "keymap":
		form := src.Keymap_form(&keymap)
		err := form.Run()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	case "timezone":
		form := src.Timezone_form(&timezone)
		err := form.Run()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	case "diskpart":
		out, _ := exec.Command("lsblk", "-dn", "-o", "NAME").Output()
		disks := strings.Fields(string(out))
		form := src.Diskpart_form(disks, &selected)
		err := form.Run()
		if err != nil {
			fmt.Println("Error:", err)
			return
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
			return
		}
		fmt.Println("Installation aborted by user.")
		os.Exit(1)
		return

	}

}
