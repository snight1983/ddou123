package ddlib

import "go.uber.org/zap"

var (

	//CACHEMAX .
	CACHEMAX uint32 = 8192
	//ROOTFOLDER .
	ROOTFOLDER string = "ROOTFOLDER"
	//SIGNINGKEY .
	SIGNINGKEY string = "11111111"

	// GLoger log
	GLoger *zap.Logger = nil
	//GConfiger .
	GConfiger *Config = nil
)

// Config 配置信息
type Config struct {
	FileFolder []string
	DBRootPath string
}
