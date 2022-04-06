package app

type Icon struct {
	Name  string
	Title string
	Hint  string
}

func (Icon) IsLauncher() bool {
	return true
}
