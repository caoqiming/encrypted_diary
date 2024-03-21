package diary

import (
	"time"

	"github.com/rivo/tview"
)

// 激活输入框，修改左侧按钮
func (a *DiaryApp) ActivateTextArea() func() {
	return func() {
		a.Grid.RemoveItem(a.Main)
		a.TextArea = tview.NewTextArea()
		a.Grid.AddItem(a.TextArea, 0, 1, 1, 1, 0, 0, false)
		a.Menu.Clear()
		a.Menu.AddItem(a.Button3, 1, 1, false)
		a.App.SetFocus(a.TextArea)
	}
}

// 保存日志
func (a *DiaryApp) Save() func() {
	return func() {
		text := a.TextArea.GetText()
		a.App.Stop()
		SaveToLocalPath("../data", *a.PassWord, text, time.Now())
	}
}
