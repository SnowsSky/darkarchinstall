package src

import (
	"github.com/charmbracelet/huh"
)

func Main_form(opt *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Welcome to darkarchinstall !").
				Options(
					huh.NewOption("Select Hostname", "hostname"),
					huh.NewOption("Set root password", "rootpasswd"),
					huh.NewOption("Accounts", "acc"),
					huh.NewOption("Cancel & exit", "cancel"),
				).
				Value(opt),
		),
	)
}

func Hostname_form(hostname *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter Hostname:").
				Placeholder(*hostname).
				Value(hostname),
		),
	)
}

func Root_passwd(rootpasswd *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter root password:").
				Value(rootpasswd),
		),
	)
}

func Sure_form(confirm *bool, text *string) *huh.Form {
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
