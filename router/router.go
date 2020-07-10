package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//InitRouter ...
func InitRouter() *gin.Engine {

	//第一个参数是api 第二个静态问价的文件夹相对目录
	//router.StaticFS("/data", http.Dir("./data"))

	//第一个参数是api 第二个参数是具体的文件名字
	//router.StaticFile("/favicon.ico", "./resources/favicon.ico")
	gin.SetMode(gin.ReleaseMode)
	engineRoot := gin.New()
	engineRoot.StaticFS("/", http.Dir("./website/unpackage/dist/build/h5/"))
	//engineRoot.StaticFS("/", http.Dir("./website/views/"))

	engineRoot.Run(":9090")
	return engineRoot
}
