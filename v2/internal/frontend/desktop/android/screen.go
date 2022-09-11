package android

import (
	"errors"

	"github.com/wailsapp/wails/v2/internal/frontend"
)

func (f *Frontend) ScreenGetAll() ([]frontend.Screen, error) {
	return nil, errors.New("ScreenGetAll not implemented")
	/*
		r := f.mainWindow.ScreenGetAll()
		if r.Error != "" {
			return nil, errors.New(r.Error)
		}

		resultIterator := r.Result

		results := make([]frontend.Screen, resultIterator.Values())
		for resultIterator.HasValue() {
			results[resultIterator.Index()] = resultIterator.Value()
			resultIterator.Next()
		}

		return results, nil
	*/
}
