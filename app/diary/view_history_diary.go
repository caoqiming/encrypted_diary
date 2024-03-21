package diary

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/rivo/tview"
)

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if info.IsDir() {
			return nil
		}
		if strings.Contains(path, ".git") {
			return nil
		}
		*files = append(*files, path)
		return nil
	}
}

type pathAndInfo struct {
	Path string
	Info *DiaryInfo
}

func GetPathAndInfo(files []string) []*pathAndInfo {
	result := make([]*pathAndInfo, 0, len(files))
	for _, f := range files {
		result = append(result, &pathAndInfo{
			Path: f,
			Info: ReadInfoFromLocalPath(f),
		})
	}
	return result
}

type timePoint struct {
	Year       int
	Month      int
	Father     *timePoint
	Childs     map[int]*timePoint // root 下面用年做索引，年下面用月做索引
	ChildsList []*timePoint       // 内容同Childs，用于存排序后的结果
	Diaries    []*pathAndInfo
}

func newTimePoint() *timePoint {
	return &timePoint{
		Childs: make(map[int]*timePoint),
	}
}

// 将Childs排序后放入ChildsList，将Diaries排序
func (t *timePoint) Sort() {
	t.ChildsList = make([]*timePoint, 0, len(t.Childs))
	for _, p := range t.Childs {
		t.ChildsList = append(t.ChildsList, p)
	}
	sort.Slice(t.ChildsList, func(i, j int) bool {
		// 年份倒序排
		return t.ChildsList[i].Year > t.ChildsList[j].Year || t.ChildsList[i].Month < t.ChildsList[j].Month
	})

	sort.Slice(t.Diaries, func(i, j int) bool {
		return t.Diaries[i].Info.Time < t.Diaries[j].Info.Time
	})
}

func sortPathAndInfoTree(p *timePoint) {
	p.Sort()
	for _, c := range p.Childs {
		sortPathAndInfoTree(c)
	}
}

func BuildPathAndInfoTree(pathAndInfo []*pathAndInfo) *timePoint {
	root := newTimePoint()
	for _, p := range pathAndInfo {
		t := time.Unix(p.Info.Time, 0)
		year, ok := root.Childs[t.Year()]
		if !ok {
			year = newTimePoint()
			year.Year = t.Year()
			year.Father = root
			root.Childs[t.Year()] = year
		}
		month, ok := year.Childs[int(t.Month())]
		if !ok {
			month = newTimePoint()
			month.Month = int(t.Month())
			month.Father = year
			year.Childs[int(t.Month())] = month
		}
		month.Diaries = append(month.Diaries, p)
	}
	// 排序
	sortPathAndInfoTree(root)
	return root
}

// 列举所有的历史日记，会根据timestamp而不是文件名来排序，虽然理论上两者是一致的
func (a *DiaryApp) ShowHistoryDiaries() func() {
	return func() {
		a.RootSelected()
	}
}

// 显示历史日记的时候还未选择年份
func (a *DiaryApp) RootSelected() {
	if a.Root == nil {
		var files []string
		rootPath := "../data"
		if err := filepath.Walk(rootPath, visit(&files)); err != nil {
			panic(err)
		}

		pathAndInfo := GetPathAndInfo(files)
		a.Root = BuildPathAndInfoTree(pathAndInfo)
	}

	a.Menu.Clear()
	list := tview.NewList()
	list.AddItem("..", "", rune('`'), a.BackToMenu())
	index := 0
	for _, p := range a.Root.ChildsList {
		list.AddItem(fmt.Sprintf("%d", p.Year), "", rune('a'+index), a.YearSelected(p))
		index++
	}
	a.Menu.AddItem(list, 0, 1, true)
	a.App.SetFocus(list)
}

// 显示历史日记的时候还未选择月份
func (a *DiaryApp) YearSelected(p *timePoint) func() {
	return func() {
		a.Menu.Clear()
		list := tview.NewList()
		list.AddItem("..", "", rune('`'), a.BackSelected(p))
		index := 0
		for _, m := range p.ChildsList {
			list.AddItem(fmt.Sprintf("%d", m.Month), "", rune('a'+index), a.MonthSelected(m))
			index++
		}
		a.Menu.AddItem(list, 0, 1, true)
		a.App.SetFocus(list)
	}
}

// 显示历史日记的时候还未选择具体日记
func (a *DiaryApp) MonthSelected(p *timePoint) func() {
	return func() {
		a.Menu.Clear()
		list := tview.NewList()
		list.AddItem("..", "", rune('`'), a.BackSelected(p))
		index := 0
		for _, d := range p.Diaries {
			dt := time.Unix(d.Info.Time, 0)
			list.AddItem(fmt.Sprintf("%d", dt.Day()), "", rune('a'+index), a.DiarySelected(d))
			index++
		}
		a.Menu.AddItem(list, 0, 1, true)
		a.App.SetFocus(list)
	}
}

func (a *DiaryApp) DiarySelected(d *pathAndInfo) func() {
	return func() {
		info, text := ReadFromLocalPath(d.Path, *a.PassWord)
		t := time.Unix(info.Time, 0)
		text = fmt.Sprintf("%s\n\n%s", t.Format("2006-01-02 15:04:05"), text)
		a.Main.SetText(text)
	}
}

func (a *DiaryApp) BackSelected(p *timePoint) func() {
	return func() {
		p = p.Father
		if p.Year != 0 {
			a.YearSelected(p)()
			return
		}
		a.RootSelected()
	}
}

func (a *DiaryApp) BackToMenu() func() {
	return func() {
		a.Menu.Clear()
		a.InitMenu()
		a.App.SetFocus(a.Button1)
	}
}
