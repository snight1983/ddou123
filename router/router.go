package router

import "github.com/gin-gonic/gin"

//InitRouter ...
func InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engineRoot := gin.New()
	//engineRoot.Group("/")
	return engineRoot
}
