package main

import "github.com/caoqiming/encrypted_diary/app/diary"

func main() {
	diary.Init()
	diary.Diary.Run()
}
