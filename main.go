package main

import (
	"ddou123/ddfile"
	"ddou123/router"
	"fmt"
	"time"
)

func main() {

	var rootPaths []string
	rootPaths = append(rootPaths, "E:\\res\\root1")
	rootPaths = append(rootPaths, "E:\\res\\root2")
	rootPaths = append(rootPaths, "d:\\res\\root3")
	rootPaths = append(rootPaths, "c:\\res\\root4")
	ddfile.Init(rootPaths)
	router.InitRouter()

	path, size := ddfile.GetCacheFolder()
	fmt.Println(path)
	fmt.Println(size)

	//path, ishave := ddfile.IsExit("92a2cf5d21750c997dec7969c43ff7e182bc927a")
	//ddfile.SaveFile("E:\\于文文+-+体面.flac")
	//ddfile.SaveFile("E:\\于文文+-+体面.flac")
	//if ishave {
	//	fmt.Println(path)
	//}
	for {
		time.Sleep(time.Duration(2) * time.Second)
	}
}
