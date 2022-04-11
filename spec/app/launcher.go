package app

import "github.com/gotrino/fusion/spec/svg"

type Icon struct {
	Icon  svg.SVG
	Title string
	Hint  string
	Link  string
}

func (Icon) IsLauncher() bool {
	return true
}
