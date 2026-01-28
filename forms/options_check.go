package forms

import (
	"bufio"
	"darkarchinstall/fs"
	"darkarchinstall/types"
	"fmt"
	"os"
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
		Hostname_form(&hostname).Run()
	case "rootpasswd":
		Root_passwd(&rootpasswd).Run()
	case "acc":
		var acc_opt string
		Accounts_form(&acc_opt).Run()

		switch acc_opt {
		case "acc_add":
			var account types.Accounts
			Account_add_form(&account.Username, &account.Password, &account.SudoPerms).Run()
			accounts = append(accounts, account)

		case "acc_remove":
			Account_remove_form(accounts, &selectedUser).Run()
			remove_Account()
		}
	case "locales":
		err := getLocales()
		if err != nil {
			fmt.Println("Error:", err)
		}
		Locales_form(&locales, &selected).Run()
	case "keymap":
		Keymap_form(&keymap).Run()
	case "timezone":
		Timezone_form(&timezone).Run()
	case "diskpart":
		disks, err := fs.GetDisks()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		form := Diskpart_form(disks, &selected)
		err = form.Run()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	case "cancel":
		var text string = "You want to exit installation ?"
		Sure_form(&confirm, &text).Run()
		if !confirm {
			return
		}
		fmt.Println("Installation aborted by user.")
		os.Exit(1)
		return

	}

}

func remove_Account() {
	var updated []types.Accounts
	for _, acc := range accounts {
		if acc.Username != selectedUser {
			updated = append(updated, acc)
		}
	}
	accounts = updated
}
func getLocales() error {
	file, err := os.Open("/etc/locale.gen")
	if err != nil {
		return err
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
		return err
	}
	return nil
}
