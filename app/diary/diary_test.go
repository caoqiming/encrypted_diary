package diary

import (
	"fmt"
	"testing"
	"time"
)

func Test_SaveToLocalPath(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		SaveToLocalPath("/Users/bytedance/Documents/diary/datatest", "123", "test", time.Now())
	})
}

//go test -v -count=1 -run Test_ReadFromLocalPath
func Test_ReadFromLocalPath(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		info, text := ReadFromLocalPath("/Users/bytedance/Documents/diary/data/2024/2/8-8cad", "xxx")
		fmt.Println(info.Time)
		fmt.Println(text)
	})
}

func Test_ConvertOldDiary(t *testing.T) {
	t.Run("convert old diaries", func(t *testing.T) {
		ConvertOldDiary("/Users/bytedance/Documents/diary/old",
			"/Users/bytedance/Documents/diary/data",
			"xxx")
	})
}
