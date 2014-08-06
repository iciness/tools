package log2file

import (
	"fmt"
	"log"
	"os"
)

var fileName map[int]string = map[int]string{}

func SetFileName(idx int, filename string, flag int) {
	fileName[idx] = filename
	if filename != "" && flag == 0 {
		os.Remove(filename)
	}
}

func L2F(fileIdx int, content string) {
	if fileName[fileIdx] != "" {
		file, err := os.OpenFile(fileName[fileIdx], os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()
		lf := log.New(file, "", log.LstdFlags)
		lf.SetFlags(0)
		lf.Printf("%s\r\n", content)
	}
	fmt.Printf("%s\r\n", content)
}
