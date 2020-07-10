package main

import (
	"ddou123/ddcloud"
	"ddou123/ddfile"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"
)

func dump() {
	re := recover()
	if nil != re {
		path, _ := os.Getwd()
		baselog := filepath.Join(path, `.\exit.log`)
		file, err := os.OpenFile(baselog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		defer file.Close()
		if err != nil {
			log.Fatalln("Failed to open error log file", err)
			return
		}
		elog := log.New(io.MultiWriter(file, os.Stderr), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		elog.Println(fmt.Sprintf("panic:rn%v", re))
		elog.Println(string(debug.Stack()))
	}
}

func main() {
	defer dump()
	//path := filepath.VolumeName("E:\\res\\root1\\")
	//path := filepath.VolumeName("/root/b/c/d/e/f")
	//fmt.Println(path)
	//fmt.Println(err)

	ddfile.AddFolder("E:\\res\\root1")
	ddfile.AddFolder("E:\\res\\root2")
	ddfile.AddFolder("D:\\res\\root3")
	ddfile.AddFolder("C:\\res\\root4")
	var rootPaths []string
	rootPaths = append(rootPaths, "E:\\res\\root1")
	rootPaths = append(rootPaths, "E:\\res\\root2")
	rootPaths = append(rootPaths, "D:\\res\\root3")
	rootPaths = append(rootPaths, "C:\\res\\root4")
	ddcloud.Init(rootPaths)

	for {
		time.Sleep(time.Duration(2) * time.Second)
	}
}
