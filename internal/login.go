package internal

import (
	"changeme/internal/service"
	"changeme/internal/utils"
	"encoding/json"
	"github.com/evercyan/brick/xencoding"
	"github.com/evercyan/brick/xfile"
	"gopkg.in/yaml.v3"
	"log"
)

// GetConfig 获取配置信息
func (a *App) GetConfig() *utils.Response {
	a.Log.Infof("GetConfig resp: %s", xencoding.JSONEncode(a.Login))
	return utils.Success(a.Login)
}

// ----------------------------------------------------------------

// SetConfig 更新配置信息
func (a *App) SetConfig(content string) *utils.Response {
	a.Log.Infof("SetConfig content: %v", content)
	Login := &service.Login{}
	if err := json.Unmarshal([]byte(content), Login); err != nil {
		a.Log.Errorf("SetConfig Unmarshal err: %v", err)
		return utils.Fail(err.Error())
	}

	//if err := Login.Update(config.GitMarkFile, xtime.Format(time.Now(), "ymdhis")); err != nil {
	//	a.Log.Errorf("SetConfig GitHub err: %v", err)
	//	return utils.Fail("无效的 Git 配置")
	//}

	a.Login = Login

	fileList := make([]service.LoginInfo, 0)
	a.DB.Order("create_at DESC").Find(&fileList)
	a.Log.Infof("GetLoginList count: %v", len(fileList))
	if !a.Login.Check(fileList) {
		return utils.Fail("账号密码错误！")
	}

	b, _ := yaml.Marshal(a.Login)
	log.Printf("login:%+v", a.Login)
	if err := xfile.Write(a.AppPath.CfgFile, string(b)); err != nil {
		a.Log.Errorf("SetConfig Write err: %v", err)
		return utils.Fail(err.Error())
	}
	// 数据库
	//a.database()
	return utils.Success("操作成功")
}

// ----------------------------------------------------------------
