package diary

import (
	"crypto/sha256"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *DiaryApp) InitPassword() {
	// 创建新的输入框
	a.InputField = tview.NewInputField().SetLabel("Password: ").SetMaskCharacter('*')
	a.InputField.SetDoneFunc(
		func(key tcell.Key) {
			password := a.InputField.GetText()
			if len(password) == 0 {
				password = "default"
			}
			a.PassWord = &password
			a.App.SetRoot(a.Grid, true)
			a.App.SetFocus(a.Button1)
			a.Main.SetText(fmt.Sprintf("password hash: %x", sha256.Sum256([]byte(password))))
		})
	// 将inputField设置为应用程序的根启动应用程序
	a.App.SetRoot(a.InputField, true)
}
