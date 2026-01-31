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
	bootloader   string = "grub"
)

func Options_check(opt string) {
	//checks
	switch opt {
	case "hostname":
		HostnameForm(&hostname).Run()
	case "rootpasswd":
		RootPasswd(&rootpasswd).Run()
	case "acc":
		var acc_opt string
		AccountsForm(&acc_opt).Run()

		switch acc_opt {
		case "acc_add":
			var account types.Accounts
			AccountAddForm(&account.Username, &account.Password, &account.SudoPerms).Run()
			accounts = append(accounts, account)

		case "acc_remove":
			AccountRemoveForm(accounts, &selectedUser).Run()
			remove_Account()
		}
	case "locales":
		err := getLocales()
		if err != nil {
			fmt.Println("Error:", err)
		}
		LocalesForm(&locales, &selected).Run()
	case "keymap":
		KeymapForm(&keymap).Run()
	case "timezone":
		TimezoneForm(&timezone).Run()
	case "diskpart":
		disks, err := fs.GetDisks()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		DiskpartForm(disks, &selected).Run()
		if selected == "back" {
			return
		}
		err = fs.EditDisk(selected)
		if err != nil {
			fmt.Println("Error editing disk:", err)
			return
		}
		for {
			MainForm(&opt).Run()
			Options_check(opt)
		}
	case "bootloader":
		BootLoaderForm(&bootloader).Run()
	case "cancel":
		var text string = "You want to exit installation ?"
		ConfirmForm(&confirm, &text).Run()
		if !confirm {
			return
		}
		fmt.Println("Installation aborted by user.")
		os.Exit(1)
		return

	}

}

func CheckRoot() error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("You need to run this program as root")
	}
	return nil
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
		// Skip the first 17 lines (header information)
		lineCount++
		if lineCount <= 17 {
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
