package diary

import (
	"crypto/sha256"
	"diary/app/diary/pkg/encript"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

type DiaryInfo struct {
	Time int64
}

// 存储日志的格式，详见readme
type DiaryFormat struct {
	Info DiaryInfo
	Data string
}

func SaveToLocalPath(rootPath, password, text string, diaryTime time.Time) {
	diaryFormat := DiaryFormat{
		Info: DiaryInfo{Time: diaryTime.Unix()},
	}

	diaryData := encript.DiaryData{}
	diaryData.SetText(text)
	diaryData.Encrypt(password)
	hash := sha256.Sum256(diaryData.Data)

	diaryFormat.Data = base64.StdEncoding.EncodeToString(diaryData.Data)
	// Create a directory path down to the day and use the current time for the file name
	dayFolderPath := path.Join(rootPath, fmt.Sprint(diaryTime.Year()), fmt.Sprintf("%02d", diaryTime.Month()), fmt.Sprintf("%02d", diaryTime.Day()))
	fileName := fmt.Sprintf("%02d-%02d-%02d-%x.json", diaryTime.Hour(), diaryTime.Minute(), diaryTime.Second(), hash[:2]) // Using .json extension for clarity
	filePath := path.Join(dayFolderPath, fileName)
	err := os.MkdirAll(dayFolderPath, 0750) // Ensure all directories up to the day are created
	if err != nil {
		log.Fatal("fail to mkdir", err)
	}

	DiaryFormatJson, err := json.Marshal(diaryFormat)
	if err != nil {
		log.Fatal("fail to marshal diary format", err)
	}

	// Write to the file
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	_, err = file.Write(DiaryFormatJson)
	if err != nil {
		log.Fatal(err)
	}

	file.Sync()
}

// 从本地读取一个日记
func ReadFromLocalPath(path, password string) (*DiaryInfo, string) {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	diaryFormat := DiaryFormat{}
	if err = json.Unmarshal(content, &diaryFormat); err != nil {
		log.Fatal(err)
	}

	diaryData := encript.DiaryData{}
	diaryData.Data, err = base64.StdEncoding.DecodeString(diaryFormat.Data)
	if err != nil {
		log.Fatal("fail to decode base64", err)
	}
	diaryData.Decrypt(password)
	return &diaryFormat.Info, diaryData.GetText()
}

// 从本地读取一个日记的info
func ReadInfoFromLocalPath(path string) *DiaryInfo {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	diaryFormat := DiaryFormat{}
	if err = json.Unmarshal(content, &diaryFormat); err != nil {
		log.Fatal(err)
	}

	return &diaryFormat.Info
}

func ConvertOldDiary(oldPath, newPath, password string) {
	oldDiary := encript.LoadOldDiary(oldPath, password)

	// 定义布局，该布局必须使用 "2006-01-02 15:04:05" 格式
	layout := "2006-01-02 15:04:05"
	// 加载北京时区
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println(err)
	}

	for _, od := range oldDiary {
		timeStr := fmt.Sprintf("%d-%02d-%02d %02d:%02d:00", od.Info.Year, od.Info.Month, od.Info.Day, od.Info.Hour, od.Info.Minute)
		odt, err := time.ParseInLocation(layout, timeStr, location)
		if err != nil {
			log.Fatal(err)
		}
		SaveToLocalPath(newPath, password, od.Text, odt)
	}
}
