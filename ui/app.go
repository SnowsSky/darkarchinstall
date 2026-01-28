package ui

import (
	"darkarchinstall/types"

	"github.com/charmbracelet/huh"
)

type Application struct {
	types.Config
}

func NewApplication() *Application {
	return &Application{}
}
