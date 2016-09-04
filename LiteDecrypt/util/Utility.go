package util

import (
	"fmt"
	"io"
	"os"
	"strings"
)

/**
*获取指定长度的字符子串
 */
func SubString(str string, begin, length int) (substr string) {
	// 将字符串的转换成[]rune
	rs := []rune(str)
	lth := len(rs)

	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}

	//	json := ""
	//	for i := begin; i < end; i++ {
	//		r := rs[i]
	//		rint := int(r)
	//		if rint < 128 {
	//			json += string(r)
	//		} else {
	//			json += "\\u" + strconv.FormatInt(int64(rint), 16) // json
	//		}
	//	}

	// 返回子串
	return string(rs[begin:end])
}

/**
*判断文件是否存在
 */
func IsFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

func IsFile(file string) bool {
	f, e := os.Stat(file)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

func CopyFile(src, dest string) (length int64, err error, destPath string) {
	destDir, err := os.Stat(dest)
	if err != nil {
		return 0, err, dest
	}
	if destDir.IsDir() {
		lastIndex := strings.LastIndex(src, "\\")
		name := SubString(src, lastIndex, len(src))
		destPath = dest + name
	} else {
		destPath = dest
	}
	if strings.EqualFold(src, destPath) { //如果路径一样，则不拷贝，要不然拷贝后为0
		fmt.Println("src is equal to dest")
		return 0, nil, destPath
	}
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err, destPath
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return 0, err, destPath
	}
	defer destFile.Close()
	length, err = io.Copy(destFile, srcFile)
	return length, err, destPath
}
