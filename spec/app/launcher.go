package app

type Icon struct {
	Icon  string
	Title string
	Hint  string
	Link  string
}

func (Icon) IsLauncher() bool {
	return true
}
