package internal

import (
	"changeme/config"
	"changeme/internal/service"
	"changeme/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/evercyan/brick/xfile"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
	"log"
	"os/exec"
	"path"
	"reflect"
	"strings"
	"sync"
)

var logger = utils.NewLogger("app")

// App struct
type App struct {
	Ctx   context.Context
	Log   *logrus.Logger
	Login *service.Login

	DB      *gorm.DB //LoginDB
	CfgFile string
	LogFile string
	DBFile  string

	Vm   service.VideoManage
	Info AppInfo

	set SettingManage
}
type AppInfo struct {
	Video service.VideoTotalInfo
	User  []service.LoginInfo
	Dir   service.DirInfo
}

func NewApp() *App {
	return &App{}
}

func (a *App) OnStartup(ctx context.Context) {
	a.Ctx = ctx
	//res, _ := service.CmdDockerStart()

	//路径等配置项----------------------------------------------------
	a.set = NewSet()
	// 初始化 initEsData 和视频部分 ----------------------------------------------------
	a.Vm = a.NewVideo()
	fmt.Printf("%v - %v", &a.set.Data.Bak.EsIndex, a.Vm.ElasticIndex)
	//----------------------------------------------------
	cfgPath := utils.GetCfgPath()
	// 日志
	a.LogFile = fmt.Sprintf(config.LogFile, cfgPath)
	a.Log = utils.NewLogger(a.LogFile)
	a.Log.Info("OnStartup begin")

	//login 信息
	a.CfgFile = fmt.Sprintf(config.CfgFile, cfgPath)
	a.Login = &service.Login{}
	if err := yaml.Unmarshal([]byte(xfile.Read(a.CfgFile)), a.Login); err != nil {
		a.Log.Errorf("OnStartup cfgfile err: %v", err)
	}

	// 数据库文件
	a.DBFile = fmt.Sprintf(config.DBFile, cfgPath)
	a.database()

	//home 用户信息
	fileList := make([]service.LoginInfo, 0)
	a.DB.Find(&fileList)
	a.Info.User = fileList

	//文件信息

	// CfgDirInfo 缓存到本地
	dirPath := fmt.Sprintf(config.CfgDirInfo, cfgPath)

	// CfgDirInfo 读配置
	list := service.DirInfo{}
	if err := yaml.Unmarshal([]byte(xfile.Read(dirPath)), &list); err != nil {
		// CfgDirInfo 信息
		a.Info.Dir = service.GetDirInfo(config.HomeInfoDir)
		if err = service.PutDirInfoDesk(a.Info.Dir, dirPath); err != nil {
			log.Printf(" dirPath Unmarshal failed,err: %v", err)
		}
	} else {
		//读取成功后 将值返回
		a.Info.Dir = list
	}

	return

}

func (a *App) NewVideo() service.VideoManage {
	var data = service.VideoManage{}

	//res, err := service.CmdEsRun()
	//if err != nil {
	//	log.Fatal("CmdEsRun failed,err:%v", err)
	//}

	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		if err != nil {
			logger.Printf("initEsData connect failed,err:%v", err)
			return data
		}
	}
	data = service.VideoManage{
		ElasticIndex:     a.set.Data.Bak.EsIndex,
		Client:           client,
		Err:              nil,
		PageSize:         10,
		VideoDefaultJson: config.VideoDefaultJson,
	}
	data.TotalInfo = data.GetInfo()

	return data
}
func (a *App) database() {
	a.Log.Info("OnStartup migrate begin")
	// 1. 如果 cantor.db 存在, 初始化, 返回
	if xfile.IsExist(a.DBFile) {
		a.DB = service.NewDB(a.DBFile)

	} else {
		a.DB = service.NewDB(a.DBFile)
		//list := make([]service.LoginInfo,0)
		//无本地文件 初始化固定字符串至本地 sqliteDB 中
		jsonContent := service.InitLoginTable()
		list := make([]service.LoginInfo, 0)
		if err := json.Unmarshal([]byte(jsonContent), &list); err != nil {
			a.Log.Errorf("OnStartup migrate json err: %v", err)
			return
		}
		fmt.Printf("--------------------- %v", list)
		success := 0
		for i := len(list) - 1; i >= 0; i-- {
			res := a.DB.Create(&list[i])
			if res.Error == nil {
				success++
			}
		}
		a.Log.Infof("OnStartup migrate json success: %d", success)
	}

	return
}

// diag ...
func (a *App) diag(message string, buttons ...string) (string, error) {
	if len(buttons) == 0 {
		buttons = []string{
			config.BtnConfirmText,
		}
	}
	return runtime.MessageDialog(a.Ctx, runtime.MessageDialogOptions{
		Type:          runtime.InfoDialog,
		Title:         config.Title,
		Message:       message,
		CancelButton:  config.BtnConfirmText,
		DefaultButton: config.BtnConfirmText,
		Buttons:       buttons,
	})
}

// Menu 应用菜单
func (a *App) Menu() *menu.Menu {
	return menu.NewMenuFromItems(
		menu.SubMenu("文件", menu.NewMenuFromItems(
			menu.Text("关于 华章文件系统", nil, func(_ *menu.CallbackData) {
				a.diag(config.Description)
			}),
			menu.Text("检查更新", nil, func(_ *menu.CallbackData) {
				a.diag(config.VersionNewMsg)
			}),
			menu.Separator(),
			menu.Text("上传图片", keys.CmdOrCtrl("O"), func(_ *menu.CallbackData) {
				runtime.EventsEmit(a.Ctx, config.EventUploadBegin)
				resp := a.BatchUploadFile()
				if resp.Code == 0 {
					runtime.EventsEmit(a.Ctx, config.EventUploadSuccess, resp.Data)
				} else {
					runtime.EventsEmit(a.Ctx, config.EventUploadFail, resp.Msg)
				}
			},
			),
			menu.Separator(),
			menu.Text("退出", keys.CmdOrCtrl("Q"), func(_ *menu.CallbackData) {
				runtime.Quit(a.Ctx)
			}),
		)),
		menu.EditMenu(),
		menu.SubMenu("帮助", menu.NewMenuFromItems(
			menu.Text(
				"打开配置文件",
				keys.Combo("C", keys.CmdOrCtrlKey, keys.ShiftKey),
				func(_ *menu.CallbackData) {
					if !xfile.IsExist(a.CfgFile) {
						a.diag("文件不存在, 请先更新配置")
						return
					}
					_, err := exec.Command("open", a.CfgFile).Output()
					if err != nil {
						a.diag("操作失败: " + err.Error())
						return
					}
				},
			),
			menu.Text(
				"打开日志文件",
				keys.Combo("L", keys.CmdOrCtrlKey, keys.ShiftKey),
				func(_ *menu.CallbackData) {
					if !xfile.IsExist(a.LogFile) {
						a.diag("文件不存在, 请先更新配置")
						return
					}
					_, err := exec.Command("open", a.LogFile).Output()
					if err != nil {
						a.diag("操作失败: " + err.Error())
						return
					}
				},
			),
			menu.Separator(),
			menu.Text(
				"打开应用主页",
				keys.Combo("H", keys.CmdOrCtrlKey, keys.ShiftKey),
				func(_ *menu.CallbackData) {
					runtime.BrowserOpenURL(a.Ctx, config.GitRepoURL)
				},
			),
		)),
	)
}

// BatchUploadFile ...
func (a *App) BatchUploadFile() *utils.Response {
	//if t.Git.Repo == "" {
	//	return internal.Fail("请先更新配置")
	//}
	files, err := runtime.OpenMultipleFilesDialog(a.Ctx, runtime.OpenDialogOptions{
		Title: "选择图片",
		Filters: []runtime.FileFilter{{
			DisplayName: "Images (*.png;*.jpg;*.jpeg;*.gif)",
			Pattern:     "*.png;*.jpg;*.jpeg;*.gif",
		}},
	})
	if err != nil {
		return utils.Fail(err.Error())
	}
	if len(files) == 0 {
		return utils.Fail("请选择至少一张图片")
	}
	if len(files) > config.MaxFileCount {
		return utils.Fail(fmt.Sprintf("最多可选择 %d 张图片", config.MaxFileCount))
	}
	for _, file := range files {
		err := a.CheckFile(file)
		if err != nil {
			return utils.Fail(fmt.Sprintf("%s: %s", path.Base(file), err.Error()))
		}
	}

	var wg sync.WaitGroup
	var mx sync.Mutex
	count := 0
	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			_, err := a.Upload(file)
			if err != nil {
				a.Log.Errorf("BatchUploadFile file: %s , err: %s", file, err.Error())
				return
			}
			mx.Lock()
			count++
			mx.Unlock()
		}(file)
	}
	wg.Wait()

	if count == 0 {
		return utils.Fail("上传图片失败")
	}
	//go a.SyncDatabase() 同步本地数据库文件到远程仓库
	return utils.Success(fmt.Sprintf("上传图片 %d 张, 成功 %d 张", len(files), count))
}

// IsContains ...
func IsContains(src, v interface{}) bool {
	srcValue := reflect.ValueOf(src)
	switch reflect.TypeOf(src).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < srcValue.Len(); i++ {
			if srcValue.Index(i).Interface() == v {
				return true
			}
		}
	case reflect.Map:
		return srcValue.MapIndex(reflect.ValueOf(v)).IsValid()
	}
	return false
}

// CheckFile  校验上传文件
func (a *App) CheckFile(filepath string) error {
	if filepath == "" {
		return fmt.Errorf("请选择图片文件")
	}
	if a.Login.Repo == "" {
		return fmt.Errorf("请设置 Login 配置")
	}
	// 文件格式校验
	fileExt := strings.ToLower(path.Ext(filepath))
	if IsContains(config.AllowFileExts, fileExt) {
		return fmt.Errorf("仅支持以下格式: %s", strings.Join(config.AllowFileExts, ", "))
	}
	// 文件大小校验
	fileSize := xfile.Size(filepath)
	a.Log.Infof("UploadFile fileSize: %v", fileSize)
	if fileSize > config.MaxFileSize {
		return fmt.Errorf("最大支持 4M 的文件")
	}
	return nil
}

// OnDomReady ...
func (a *App) OnDomReady(ctx context.Context) {
	a.Log.Info("OnDomReady")
	return
}

// OnShutdown ...
func (a *App) OnShutdown(ctx context.Context) {
	a.Log.Info("OnShutdown")
	return
}

// OnBeforeClose ...
func (a *App) OnBeforeClose(ctx context.Context) bool {
	a.Log.Info("OnBeforeClose")
	// 返回 true 将阻止程序关闭
	return false
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
