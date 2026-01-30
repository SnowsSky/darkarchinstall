package forms

import (
	"darkarchinstall/types"

	"github.com/charmbracelet/huh"
)

func MainForm(opt *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Welcome to darkarchinstall !").
				Options(
					huh.NewOption("Select Hostname <"+hostname+">", "hostname"),
					huh.NewOption("Set root password", "rootpasswd"),
					huh.NewOption("Accounts", "acc"),
					huh.NewOption("Select locales", "locales"),
					huh.NewOption("Keymap <"+keymap+">", "keymap"),
					huh.NewOption("Timezone <"+timezone+">", "timezone"),
					huh.NewOption("Disk partitioning", "diskpart"),
					huh.NewOption("Install Dark Arch", "install"),
					huh.NewOption("Cancel & exit", "cancel"),
				).
				Value(opt),
		),
	)
}

func HostnameForm(hostname *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter Hostname:").
				Placeholder(*hostname).
				Value(hostname),
		),
	)
}

func RootPasswd(rootpasswd *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter root password:").
				EchoMode(huh.EchoModePassword).
				Value(rootpasswd),
		),
	)
}

func ConfirmForm(confirm *bool, text *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(*text).
				Affirmative("Yes!").
				Negative("No.").
				Value(confirm),
		),
	)
}

func AccountsForm(acc_opt *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Darkinstall - Accounts").
				Options(
					huh.NewOption("Add an account", "acc_add"),
					huh.NewOption("Remove an account", "acc_remove"),
					huh.NewOption("Go back ?", "back"),
				).
				Value(acc_opt),
		),
	)
}

func AccountAddForm(username *string, password *string, sudo *bool) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter username:").
				Placeholder("darkarch").
				Value(username),
			huh.NewInput().
				Title("Enter password:").
				EchoMode(huh.EchoModePassword).
				Value(password),
			huh.NewConfirm().
				Title("Grant sudo privileges?").
				Affirmative("Of course !").
				Negative("No !").
				Value(sudo),
		),
	)
}

func AccountRemoveForm(accounts []types.Accounts, selected *string) *huh.Form {
	options := make([]huh.Option[string], 0)

	for _, acc := range accounts {
		options = append(options,
			huh.NewOption(acc.Username, acc.Username),
		)
	}
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select user to delete").
				Options(options...).
				Value(selected),
		),
	)
}

func LocalesForm(locales *[]string, selected *string) *huh.Form {
	options := make([]huh.Option[string], 0)
	for _, loc := range *locales {
		options = append(options, huh.NewOption(loc, loc))
	}
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select locales").
				Options(options...).
				Value(selected),
		),
	)
}

func KeymapForm(keymap *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter keymap: https://man.archlinux.org/man/vconsole.conf.5").
				Placeholder("de-latin1").
				Value(keymap),
		),
	)
}

func TimezoneForm(timezone *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter timezone:").
				Placeholder("UTC").
				Value(timezone),
		),
	)
}

func DiskpartForm(disks []string, selected *string) *huh.Form {
	options := make([]huh.Option[string], 0)

	for _, disk := range disks {
		options = append(options, huh.NewOption(disk, disk))
	}
	options = append(options, huh.NewOption("Go back ?", "back"))
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a disk to edit.").
				Options(options...).
				Value(selected),
		),
	)
}
