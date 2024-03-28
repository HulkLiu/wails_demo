package config

// 应用
const (
	App     = "VideoManager"
	Version = "v0.1.0"
)

const (
	Title          = App + " " + Version
	Description    = "一个简单好用的视频应用"
	VersionNewMsg  = "当前已经是最新版本!"
	VersionOldMsg  = "最新版本: %s, 是否立即更新?"
	BtnConfirmText = "确定"
	BtnCancelText  = "取消"
)

// 窗口尺寸
const (
	VideoDefaultJson = "VideoDefaultJson.conf"

	Width  = 1024
	Height = 768

	EventUploadBegin   = "event.upload.begin"
	EventUploadSuccess = "event.upload.success"
	EventUploadFail    = "event.upload.fail"

	//git
	GitApiURL   = "https://api.github.com/repos/%s/%s/contents/%s"
	GitTagURL   = "https://api.github.com/repos/evercyan/cantor/tags"
	GitFileURL  = "https://cdn.jsdelivr.net/gh/%s/%s/%s"
	GitDBFile   = "service/cantor.db"
	GitFilePath = "service/%s/%s%s"
	GitMarkFile = "mark"
	GitMessage  = "upload by cantor"
	//GitRepoURL  = "https://github.com/evercyan/cantor"
	GitRepoURL = "https://github.com/HulkLiu/wails_demo"
	GitAppURL  = GitRepoURL + "/releases/tag/%s"
)

var (
	EsIndex     = "vjshi10"
	ExcelDir    = "data/excel/"
	HomeInfoDir = "G:\\书"
	ContainerID = "71b0a4b9b297c10500a33588579a8426b3db47b5f6c1af152664dca8dea29c91"
)

// 图片配置
var (
	AllowFileExts       = []string{".png", ".gif", ".jpg", ".jpeg"}
	MaxFileSize   int64 = 4 * 1024 * 1024
	MaxFileCount        = 10
)

// 文件配置
var (
	CfgFile     = "%s/config.yaml"
	CfgDirInfo  = "%s/dirInfo.yaml"
	LogFile     = "%s/app.log"
	DBFile      = "%s/videoDB.db"
	SettingFile = "%s/setting.yaml"
	//LoginFile = "%s/login.yaml"
)

var (
	TimeFormat = "20060102150405"
)
