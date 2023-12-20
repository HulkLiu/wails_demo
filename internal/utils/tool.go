package utils

import (
	"changeme/config"
	"fmt"
	"github.com/evercyan/brick/xfile"

	"os"
	"os/user"
	"strings"
)

func GetCfgPath() string {
	userPath, err := user.Current()
	if err != nil {
		panic("获取应用配置目录失败: " + err.Error())
	}
	cfgPath := fmt.Sprintf("%s/.%s", userPath.HomeDir, strings.ToLower(config.App))
	if !xfile.IsExist(cfgPath) {
		os.Mkdir(cfgPath, os.ModePerm)
	}
	return cfgPath
}
