package forms

import (
	"darkarchinstall/types"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var theme = huh.ThemeBase()

func MainForm(opt *string) *huh.Form {
	// Theme setup
	theme.Focused.TextInput.Cursor = lipgloss.NewStyle().Foreground(lipgloss.Color("#eb0e0e"))
	theme.Focused.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("#e7130cd8")).Bold(true)
	theme.Focused.SelectedOption = lipgloss.NewStyle().Bold(true)
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
					huh.NewOption("Select bootloader <"+bootloader+">", "bootloader"),
					huh.NewOption("Select Desktop Environment", "de"),
					huh.NewOption("Install Dark Arch", "install"),
					huh.NewOption("Cancel & exit", "cancel"),
				).
				Value(opt),
		),
	).WithTheme(theme)
}

func HostnameForm(hostname *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter Hostname:").
				Placeholder(*hostname).
				Value(hostname),
		),
	).WithTheme(theme)
}

func RootPasswd(rootpasswd *string) *huh.Form {

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter root password:").
				EchoMode(huh.EchoModePassword).
				Value(rootpasswd),
		),
	).WithTheme(theme)
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
	).WithTheme(theme)
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
	).WithTheme(theme)
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
	).WithTheme(theme)
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
	).WithTheme(theme)
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
	).WithTheme(theme)
}

func KeymapForm(keymap *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter keymap: https://man.archlinux.org/man/vconsole.conf.5").
				Placeholder("de-latin1").
				Value(keymap),
		),
	).WithTheme(theme)
}

func TimezoneForm(timezone *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter timezone:").
				Placeholder("UTC").
				Value(timezone),
		),
	).WithTheme(theme)
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
	).WithTheme(theme)
}

func BootLoaderForm(bootloader *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select A Bootloader").
				Options(
					huh.NewOption("Grub", "grub"),
					huh.NewOption("SystemD-boot", "systemd-boot"),
					huh.NewOption("Limine", "limine"),
				).
				Value(bootloader),
		),
	).WithTheme(theme)
}

func SelectDEForm(de *[]string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Options(
					huh.NewOption("KDE plasma", "plasma").Selected(true),
					huh.NewOption("XFCE", "xfce"),
					huh.NewOption("GNOME", "gnome"),
				).
				Title("Select Desktop Environments").
				Value(de),
		),
	).WithTheme(theme)
}
