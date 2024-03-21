package diary

func (a *DiaryApp) InitGrid() {
	a.Grid.
		SetRows(0).
		SetColumns(30, 0).
		SetBorders(true).
		AddItem(a.Menu, 0, 0, 1, 1, 0, 0, false).
		AddItem(a.Main, 0, 1, 1, 1, 0, 0, false)
}
