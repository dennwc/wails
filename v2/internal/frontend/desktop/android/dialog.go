package android

import (
	"errors"

	"github.com/wailsapp/wails/v2/internal/frontend"
)

func (f *Frontend) OpenFileDialog(dialogOptions frontend.OpenDialogOptions) (string, error) {
	return "", errors.New("OpenFileDialog not implemented")
}

func (f *Frontend) OpenMultipleFilesDialog(dialogOptions frontend.OpenDialogOptions) ([]string, error) {
	return nil, errors.New("OpenMultipleFilesDialog not implemented")
}

func (f *Frontend) OpenDirectoryDialog(dialogOptions frontend.OpenDialogOptions) (string, error) {
	return "", errors.New("OpenDirectoryDialog not implemented")
}

func (f *Frontend) SaveFileDialog(dialogOptions frontend.SaveDialogOptions) (string, error) {
	return "", errors.New("SaveFileDialog not implemented")
}

func (f *Frontend) MessageDialog(dialogOptions frontend.MessageDialogOptions) (string, error) {
	return "", errors.New("MessageDialog not implemented")
}

/*
func (f *Frontend) OpenFileDialog(dialogOptions frontend.OpenDialogOptions) (string, error) {
	r := f.mainWindow.OpenFileDialog(android.NewOpenDialogOptions(dialogOptions))
	if r.Error != "" {
		return "", errors.New(r.Error)
	}
	return r.Result, nil
}

func (f *Frontend) OpenMultipleFilesDialog(dialogOptions frontend.OpenDialogOptions) ([]string, error) {
	r := f.mainWindow.OpenMultipleFilesDialog(android.NewOpenDialogOptions(dialogOptions))
	if r.Error != "" {
		return nil, errors.New(r.Error)
	}

	resultIterator := r.Result

	results := make([]string, resultIterator.Values())
	for resultIterator.HasValue() {
		results[resultIterator.Index()] = resultIterator.Value()
		resultIterator.Next()
	}

	return results, nil
}

func (f *Frontend) OpenDirectoryDialog(dialogOptions frontend.OpenDialogOptions) (string, error) {
	r := f.mainWindow.OpenDirectoryDialog(android.NewOpenDialogOptions(dialogOptions))
	if r.Error != "" {
		return "", errors.New(r.Error)
	}
	return r.Result, nil
}

func (f *Frontend) SaveFileDialog(dialogOptions frontend.SaveDialogOptions) (string, error) {
	r := f.mainWindow.SaveFileDialog(android.NewSaveDialogOptions(dialogOptions))
	if r.Error != "" {
		return "", errors.New(r.Error)
	}
	return r.Result, nil
}

func (f *Frontend) MessageDialog(dialogOptions frontend.MessageDialogOptions) (string, error) {
	r := f.mainWindow.MessageDialog(android.NewMessageDialogOptions(dialogOptions))
	if r.Error != "" {
		return "", errors.New(r.Error)
	}
	return r.Result, nil
}
*/
