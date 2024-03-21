package encript

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Encript(t *testing.T) {

	test := []struct {
		originalText string
	}{
		{originalText: "hello world"},
		{originalText: "你好，世界"},
		{originalText: "test \n 你好"},
	}

	for i, tt := range test {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			d := DiaryData{}
			d.SetText(tt.originalText)
			d.Encrypt("test-password")
			d2 := DiaryData{}
			d2.Data = d.Data
			d2.Decrypt("test-password")
			assert.Equal(t, tt.originalText, d2.GetText())
			d3 := DiaryData{}
			d3.Data = d.Data
			d3.Decrypt("test-password-wrong")
			assert.NotEqual(t, tt.originalText, d3.GetText())
		})
	}

}
