// LiteDecrypt project main.go
package main

import (
	"LiteDecrypt/window"
	"fmt"
)

func main() {

	//命令行输入参数 go build -ldflags="-H windowsgui" 才能不弹出dos
	//	argsNum := len(os.Args)
	//	if argsNum == 2 {
	//		filePath.in = os.Args[1]
	//		filePath.out = "decrypt_v2.txt"
	//	} else if argsNum == 3 {
	//		filePath.in = os.Args[1]
	//		filePath.out = os.Args[2]
	//	} else {
	//		filePath.in = "logs_v2.txt"
	//		filePath.out = "decrypt_v2.txt"
	//	}
	window := &window.DecryptWindow{}
	window.OpenDecryptWindow()
	fmt.Println("The main is end")
}
