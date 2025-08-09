package diary

import (
	"testing"

	"github.com/caoqiming/encrypted_diary/app/diary/pkg/utils"
)

func TestViewHistory(t *testing.T) {
	Init()
	Diary.RootPath = utils.Ptr("/Users/glimmer/Documents/code/encrypted_diary/diary-data")
	Diary.RootSelected()
}
