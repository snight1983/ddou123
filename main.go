package main

import (
	"ddou123/ddfile"
	"time"
)

func main() {

	/* test
	f, _ := os.Open("D:\\92a2cf5d21750c997dec7969c43ff7e182bc927a.mp3")
	buf := make([]byte, 1024*1024*4)
	var id int = 0
	for {
		n, _ := f.Read(buf)
		if n == 0 {
			break
		}

		part := filepath.Join("D:\\tmp", strconv.Itoa(id))
		destTmp, _ := os.Create(part)
		destTmp.Write(buf[:n])
		destTmp.Close()
		id++
	}
	*/
	//ddnet.Init()

	var rootPaths []string

	rootPaths = append(rootPaths, "E:\\res\\root1")
	rootPaths = append(rootPaths, "E:\\res\\root2")
	rootPaths = append(rootPaths, "D:\\res\\root3")
	rootPaths = append(rootPaths, "C:\\res\\root4")

	ddfile.Init(rootPaths)
	//router.InitRouter()
	ddfile.SaveFileSingle("E:\\于文文+-+体面.flac")
	ddfile.SaveFileFolder("D:\\tmp", 7)

	//path, size := ddfile.GetCacheFolder()
	//fmt.Println(path)
	//fmt.Println(size)

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
