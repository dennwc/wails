package android

import (
	"errors"

	"github.com/wailsapp/wails/v2/pkg/menu"
)

func (f *Frontend) MenuSetApplicationMenu(menu *menu.Menu) {
	if menu != nil && len(menu.Items) > 0 {
		panic(errors.New("todo: MenuSetApplicationMenu not implemented"))
	}
}

func (f *Frontend) MenuUpdateApplicationMenu() {
	if f.frontendOptions.Menu != nil && len(f.frontendOptions.Menu.Items) > 0 {
		panic(errors.New("todo: MenuUpdateApplicationMenu not implemented"))
	}
}
