package main

import (
	"fmt"
	"os"
)

func pathExist(path string) bool {
	_, err := os.Stat(path)
	// 文件不存在
	if err != nil {
		return false
	}
	// 目录不存在
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// Mkdir 创建文件夹
func Mkdir(date string, rootAddress string) string {
	year := date[0:4] + "/"
	month := date[4:6] + "/"
	if pathExist(rootAddress+year) == false {
		err := os.Mkdir(rootAddress+year, 0755)
		if err != nil {
			fmt.Println("文件夹创建失败 ", err)
		}
	}
	if pathExist(rootAddress+year+month) == false {
		err := os.Mkdir(rootAddress+year+month, 0755) // 不应为0777, 否则会有权限报错
		if err != nil {
			fmt.Println("文件夹创建失败 ", err)
		}
	}

	return rootAddress + year + month
}
