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
	"os"
	"os/exec"
	"path"
	"path/filepath"
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

	AppPath AppPath

	DB *gorm.DB //LoginDB
	Vm *service.VideoManage

	Info       AppInfo
	set        SettingManage
	taskManage service.TaskManager
}

type AppPath struct {
	CfgPath string
	CfgFile string
	LogFile string
	DBFile  string
	EsCpath string
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

	a.InitPath() //初始化Path信息

	a.set = NewSet(a.AppPath.CfgPath) //功能 - 设置管理

	a.Vm = initEs(a.AppPath.EsCpath) //初始化 Es

	a.database() // init sqlite

	a.InitLogin() // 加载缓存登录信息

	a.InitConfig() //加载配置信息

	a.taskManage.TasksDB = service.NewTaskDB(a.AppPath.DBFile) //任务管理初始化

	return
}

type Config struct {
	Elasticsearch struct {
		Address string `yaml:"address"`
		Sniff   bool   `yaml:"sniff"`
	} `yaml:"elasticsearch"`
}

func readConfig(path string) (*Config, error) {
	cfg := &Config{}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func initEs(path string) *service.VideoManage {
	cfg, err := readConfig(path)
	if err != nil {
		log.Fatal(err)
	}

	client, err := elastic.NewClient(
		elastic.SetURL(cfg.Elasticsearch.Address),
		elastic.SetSniff(cfg.Elasticsearch.Sniff),
	)
	if err != nil {
		log.Fatal(err)
	}

	data := service.VideoManage{
		ElasticIndex:     config.EsIndex,
		Client:           client,
		Err:              nil,
		PageSize:         10,
		VideoDefaultJson: config.VideoDefaultJson,
	}
	data.TotalInfo = data.GetInfo()

	return &data
}

// InitPath 初始化路径
func (a *App) InitPath() {
	//程序目录
	a.AppPath.CfgPath = utils.GetCfgPath()

	//日志文件
	a.AppPath.LogFile = fmt.Sprintf(config.LogFile, a.AppPath.CfgPath)

	//配置文件
	a.AppPath.CfgFile = fmt.Sprintf(config.CfgFile, a.AppPath.CfgPath)

	//sqlite
	a.AppPath.DBFile = fmt.Sprintf(config.DBFile, a.AppPath.CfgPath)

	dir, _ := os.Getwd()

	a.AppPath.EsCpath = filepath.Join(dir, "esConfig.yaml")

	// 日志
	a.Log = utils.NewLogger(a.AppPath.LogFile)
}

func (a *App) InitLogin() {
	//login 信息
	a.Login = &service.Login{}
	if err := yaml.Unmarshal([]byte(xfile.Read(a.AppPath.CfgFile)), a.Login); err != nil {
		a.Log.Errorf("OnStartup cfgfile err: %v", err)
	}
	//home 用户信息
	var fileList []service.LoginInfo
	a.DB.Find(&fileList)
	a.Info.User = fileList
}

func (a *App) InitConfig() {
	// 文件信息 CfgDirInfo 缓存到本地
	dirPath := fmt.Sprintf(config.CfgDirInfo, a.AppPath.CfgPath)

	// 首页读配置
	list := service.DirInfo{}

	if xfile.IsExist(dirPath) {
		//如果配置文件存在 就读文件，失败再去重新获取
		if err := yaml.Unmarshal([]byte(xfile.Read(dirPath)), &list); err != nil {

			a.Info.Dir = service.GetDirInfo(config.HomeInfoDir)
			if err = service.PutDirInfoDesk(a.Info.Dir, dirPath); err != nil {
				log.Printf(" dirPath Unmarshal failed,err: %v", err)
			}
		} else {
			//读取成功后 将值返回
			a.Info.Dir = list
		}
	} else {
		//如果配置文件不存在 重新获取
		a.Info.Dir = service.GetDirInfo(config.HomeInfoDir)
		if err := service.PutDirInfoDesk(a.Info.Dir, dirPath); err != nil {
			log.Printf(" dirPath Unmarshal failed,err: %v", err)
		}
	}
}

func (a *App) database() {

	a.Log.Info("OnStartup migrate begin")
	// 1. 如果 cantor.db 存在, 初始化, 返回
	if xfile.IsExist(a.AppPath.DBFile) {
		a.DB = service.NewDB(a.AppPath.DBFile)
	} else {
		a.DB = service.NewDB(a.AppPath.DBFile)

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
					if !xfile.IsExist(a.AppPath.CfgFile) {
						a.diag("文件不存在, 请先更新配置")
						return
					}
					_, err := exec.Command("open", a.AppPath.CfgFile).Output()
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
					if !xfile.IsExist(a.AppPath.LogFile) {
						a.diag("文件不存在, 请先更新配置")
						return
					}
					_, err := exec.Command("open", a.AppPath.LogFile).Output()
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
