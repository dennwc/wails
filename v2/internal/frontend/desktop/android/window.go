package android

import (
	"errors"

	"github.com/wailsapp/wails/v2/pkg/options"
)

func (f *Frontend) WindowSetTitle(title string) {
	panic(errors.New("todo: WindowSetTitle not implemented"))
}

func (f *Frontend) WindowShow() {
	panic(errors.New("todo: WindowShow not implemented"))
}

func (f *Frontend) WindowHide() {
	panic(errors.New("todo: WindowHide not implemented"))
}

func (f *Frontend) WindowCenter() {
	panic(errors.New("todo: WindowCenter not implemented"))
}

func (f *Frontend) WindowToggleMaximise() {
	panic(errors.New("todo: WindowToggleMaximise not implemented"))
}

func (f *Frontend) WindowMaximise() {
	panic(errors.New("todo: WindowMaximise not implemented"))
}

func (f *Frontend) WindowUnmaximise() {
	panic(errors.New("todo: WindowUnmaximise not implemented"))
}

func (f *Frontend) WindowMinimise() {
	panic(errors.New("todo: WindowMinimise not implemented"))
}

func (f *Frontend) WindowUnminimise() {
	panic(errors.New("todo: WindowUnminimise not implemented"))
}

func (f *Frontend) WindowSetAlwaysOnTop(b bool) {
	panic(errors.New("todo: WindowSetAlwaysOnTop not implemented"))
}

func (f *Frontend) WindowSetPosition(x int, y int) {
	panic(errors.New("todo: WindowSetPosition not implemented"))
}

func (f *Frontend) WindowGetPosition() (int, int) {
	panic(errors.New("todo: WindowGetPosition()  not implemented"))
}

func (f *Frontend) WindowSetSize(width int, height int) {
	panic(errors.New("todo: WindowSetSize not implemented"))
}

func (f *Frontend) WindowGetSize() (int, int) {
	panic(errors.New("todo: WindowGetSize()  not implemented"))
}

func (f *Frontend) WindowSetMinSize(width int, height int) {
	panic(errors.New("todo: WindowSetMinSize not implemented"))
}

func (f *Frontend) WindowSetMaxSize(width int, height int) {
	panic(errors.New("todo: WindowSetMaxSize not implemented"))
}

func (f *Frontend) WindowFullscreen() {
	panic(errors.New("todo: WindowFullscreen not implemented"))
}

func (f *Frontend) WindowUnfullscreen() {
	panic(errors.New("todo: WindowUnfullscreen not implemented"))
}

func (f *Frontend) WindowSetBackgroundColour(col *options.RGBA) {
	panic(errors.New("todo: WindowSetBackgroundColour not implemented"))
}

func (f *Frontend) WindowReload() {
	panic(errors.New("todo: WindowReload not implemented"))
}

func (f *Frontend) WindowReloadApp() {
	panic(errors.New("todo: WindowReloadApp not implemented"))
}

func (f *Frontend) WindowSetSystemDefaultTheme() {
	panic(errors.New("todo: WindowSetSystemDefaultTheme not implemented"))
}

func (f *Frontend) WindowSetLightTheme() {
	panic(errors.New("todo: WindowSetLightTheme not implemented"))
}

func (f *Frontend) WindowSetDarkTheme() {
	panic(errors.New("todo: WindowSetDarkTheme not implemented"))
}

func (f *Frontend) WindowIsMaximised() bool {
	panic(errors.New("todo: WindowIsMaximised not implemented"))
}

func (f *Frontend) WindowIsMinimised() bool {
	panic(errors.New("todo: WindowIsMinimised not implemented"))
}

func (f *Frontend) WindowIsNormal() bool {
	panic(errors.New("todo: WindowIsNormal not implemented"))
}

func (f *Frontend) WindowIsFullscreen() bool {
	return f.mainWindow.IsFullScreen()
}
