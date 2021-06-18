package util

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"genshin-sign-helper/util/constant"
)

func GetMD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}

func GetDs() string {
	currentTime := time.Now().Unix()
	stringRom := GetRandString(6, currentTime)
	stringAdd := fmt.Sprintf("salt=%s&t=%d&r=%s", constant.Salt, currentTime, stringRom)
	stringMd5 := GetMD5(stringAdd)
	return fmt.Sprintf("%d,%s,%s", currentTime, stringRom, stringMd5)
}

func GetRandString(len int, seed int64) string {
	bytes := make([]byte, len)
	r := rand.New(rand.NewSource(seed))
	for i := 0; i < len; i++ {
		b := r.Intn(36)
		if b > 9 {
			b += 39
		}
		b += 48
		bytes[i] = byte(b)
	}
	return string(bytes)
}

//ReadFile 读取整个文件
func ReadFile(path string) ([]byte, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("unable to read file." + err.Error())
	}
	return b, nil
}

//ReadFileAllLine 按行读取文件
func ReadFileAllLine(path string, handle func(string)) error {
	f, err := os.Open(path)
	if err != nil {
		return errors.New("unable to read file." + err.Error())
	}
	defer f.Close()

	bufReader := bufio.NewReader(f)

	line, is, err := bufReader.ReadLine()
	for ; !is && err == nil; line, is, err = bufReader.ReadLine() {
		s := string(line)
		handle(s)
	}
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}
	return nil
}

func StructToJSON(v interface{}) ([]byte, error) {
	jsonByte, err := json.Marshal(v)
	if err != nil {
		return nil, errors.New("unable convert struct to json." + err.Error())
	}
	return jsonByte, nil
}

//CheckFileIsExist 判断文件是否存在，存在返回true，不存在返回false
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
