package encript

import (
	"bytes"
	"crypto/rc4"
	"encoding/binary"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
)

// cpp 版本的info
type DiaryInfo struct {
	Year   int32
	Month  int32
	Day    int32
	Hour   int32
	Minute int32
	Length int32
}

func (d *DiaryInfo) Parse(data []byte) {
	if len(data) != 24 {
		log.Fatal("DiaryInfo fail to parse, info length is not 24 bytes")
	}
	err := binary.Read(bytes.NewBuffer(data[:4]), binary.LittleEndian, &d.Year)
	if err != nil {
		log.Fatal("binary.Read failed:", err)
	}
	err = binary.Read(bytes.NewBuffer(data[4:8]), binary.LittleEndian, &d.Month)
	if err != nil {
		log.Fatal("binary.Read failed:", err)
	}
	err = binary.Read(bytes.NewBuffer(data[8:12]), binary.LittleEndian, &d.Day)
	if err != nil {
		log.Fatal("binary.Read failed:", err)
	}
	err = binary.Read(bytes.NewBuffer(data[12:16]), binary.LittleEndian, &d.Hour)
	if err != nil {
		log.Fatal("binary.Read failed:", err)
	}
	err = binary.Read(bytes.NewBuffer(data[16:20]), binary.LittleEndian, &d.Minute)
	if err != nil {
		log.Fatal("binary.Read failed:", err)
	}
	err = binary.Read(bytes.NewBuffer(data[20:24]), binary.LittleEndian, &d.Length)
	if err != nil {
		log.Fatal("binary.Read failed:", err)
	}
}

type DiaryData struct {
	Data []byte // 密文
	text string // 明文
}

func (d *DiaryData) GetText() string {
	return d.text
}

func (d *DiaryData) SetText(text string) {
	d.text = text
}

// 解码旧版本日志的密文
func (d *DiaryData) DecryptGB18030(password string) {
	key := []byte(password) // RC4的密钥

	// 创建并初始化RC4密码
	cipher, err := rc4.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	dst := make([]byte, len(d.Data)) // 创建一个用于存放结果的切片
	cipher.XORKeyStream(dst, d.Data)
	reader := simplifiedchinese.GB18030.NewDecoder()
	decodedBytes, _ := reader.Bytes(dst)
	d.text = string(decodedBytes)
}

func (d *DiaryData) Encrypt(password string) {
	key := []byte(password) // RC4的密钥

	// 创建并初始化RC4密码
	cipher, err := rc4.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}
	textBytes := []byte(d.text)
	dst := make([]byte, len(textBytes)) // 创建一个用于存放结果的切片
	cipher.XORKeyStream(dst, textBytes)
	d.Data = dst
}

func (d *DiaryData) Decrypt(password string) {
	key := []byte(password) // RC4的密钥

	// 创建并初始化RC4密码
	cipher, err := rc4.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	dst := make([]byte, len(d.Data)) // 创建一个用于存放结果的切片
	cipher.XORKeyStream(dst, d.Data)
	d.text = string(dst)
}

func visit(infoFiles, dataFiles *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".ifo") {
			*infoFiles = append(*infoFiles, path)
		} else if strings.HasSuffix(path, ".dat") {
			*dataFiles = append(*dataFiles, path)
		}

		return nil
	}
}

type OldDiaryFormat struct {
	Info DiaryInfo
	Text string
}

func LoadOldDiary(root, password string) []*OldDiaryFormat {
	var infoFiles, dataFiles []string // 日志文件路径

	if err := filepath.Walk(root, visit(&infoFiles, &dataFiles)); err != nil {
		panic(err)
	}

	result := make([]*OldDiaryFormat, 0, len(dataFiles))
	for _, file := range dataFiles {
		dataPath := file
		infoPath := file[:len(file)-3] + "ifo"
		content, err := os.ReadFile(infoPath)
		if err != nil {
			log.Fatal(err)
		}
		diaryInfo := DiaryInfo{}
		diaryInfo.Parse(content)
		// utils.JsonPrint(diaryInfo)

		content, err = os.ReadFile(dataPath)
		if err != nil {
			log.Fatal(err)
		}
		diaryData := DiaryData{Data: content}
		diaryData.DecryptGB18030(password)
		result = append(result, &OldDiaryFormat{
			Info: diaryInfo,
			Text: diaryData.GetText(),
		})
	}
	return result
}
