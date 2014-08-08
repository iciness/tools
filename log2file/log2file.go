package log2file

import (
	"fmt"
	"log"
	"os"
)

type Log2File struct {
	FileName string
}

func NewLog2File(fileName string) *Log2File {
	return &Log2File{fileName}
}

func (log2File *Log2File) Println(a ...interface{}) {
	if log2File.FileName != "" {
		file, err := os.OpenFile(log2File.FileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()
		lf := log.New(file, "", log.LstdFlags)
		lf.SetFlags(0)
		lf.Printf(fmt.Sprintln(a))
	}
	fmt.Println(a)
}
