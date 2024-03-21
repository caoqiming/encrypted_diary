package main

import "diary/app/diary"

func main() {
	diary.Init()
	diary.Diary.Run()
}
