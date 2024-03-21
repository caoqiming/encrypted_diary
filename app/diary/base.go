package diary

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type DiaryApp struct {
	App        *tview.Application
	Grid       *tview.Grid
	Button1    *tview.Button // 新日记
	Button2    *tview.Button // 查看历史日记
	Button3    *tview.Button // 保存
	Menu       *tview.Flex
	Main       *tview.TextView
	TextArea   *tview.TextArea   // 输入日记
	InputField *tview.InputField // 输入密码
	PassWord   *string
	Root       *timePoint // 日志info信息树
}

var Diary DiaryApp

func Init() {
	Diary.App = tview.NewApplication()
	Diary.Grid = tview.NewGrid()
	Diary.Button1 = tview.NewButton("new diary")
	Diary.Button2 = tview.NewButton("view history")
	Diary.Button3 = tview.NewButton("save")
	Diary.Menu = tview.NewFlex()
	Diary.Main = tview.NewTextView().SetTextAlign(tview.AlignLeft)
	Diary.InitMenu()
	Diary.InitGrid()
	Diary.InitButton()
	Diary.InitPassword()
}

func (a *DiaryApp) Run() {
	if err := a.App.
		SetRoot(a.InputField, true).
		SetFocus(a.InputField).
		EnableMouse(true).
		Run(); err != nil {
		panic(err)
	}
}

func (a *DiaryApp) InitButton() {
	a.Button1.SetSelectedFunc(a.ActivateTextArea())
	a.Button1.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyUp || event.Key() == tcell.KeyDown {
			a.App.SetFocus(a.Button2)
		}
		return event
	})
	a.Button2.SetSelectedFunc(a.ShowHistoryDiaries())
	a.Button2.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyUp || event.Key() == tcell.KeyDown {
			a.App.SetFocus(a.Button1)
		}
		return event
	})
	a.Button3.SetSelectedFunc(a.Save())
}

func (a *DiaryApp) InitMenu() {
	a.Menu.
		SetDirection(tview.FlexRow).
		AddItem(a.Button1, 1, 1, true).
		AddItem(nil, 1, 1, false).
		AddItem(a.Button2, 1, 1, false)
}
