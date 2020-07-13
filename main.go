package main

import (
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

	/*
		dataA := dduser.UserBase{}
		ret, err := msgpack.Marshal(dataA)
		fmt.Println(ret)
		fmt.Println(err)
		var buffer bytes.Buffer
		err1 := binary.Write(&buffer, binary.BigEndian, &dataA)
		fmt.Println(err1)
	*/

	defer dump()

	ddfile.AddFolder("E:\\res\\root1")
	ddfile.AddFolder("E:\\res\\root2")
	ddfile.AddFolder("D:\\res\\root3")
	ddfile.AddFolder("C:\\res\\root4")

	for {
		time.Sleep(time.Duration(2) * time.Second)
	}
}
