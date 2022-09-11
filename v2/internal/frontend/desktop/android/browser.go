package android

// BrowserOpenURL Use the default browser to open the url
func (f *Frontend) BrowserOpenURL(url string) {
	f.mainWindow.BrowserOpenURL(url)
}
