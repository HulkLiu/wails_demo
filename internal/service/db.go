package service

import (
	"github.com/evercyan/brick/xfile"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// 	 LoginArr {Repo:1 Owner:1 Email:1 AccessToken:1}
type LoginInfo struct {
	ID          int    `json:"id" gorm:"column:id;AUTO_INCREMENT;not null"`
	Repo        string `json:"user_repo" gorm:"column:user_repo;not null"`
	Email       string `json:"file_email" gorm:"column:user_email;not null"`
	Owner       string `json:"file_owner" gorm:"column:user_owner;not null"`
	AccessToken string `json:"file_accessToken" gorm:"column:accessToken;not null"`
	CreateAt    string `json:"create_at" gorm:"column:create_at;not null"`
}

// TableName ...
func (f *LoginInfo) TableName() string {
	return "user"
}

// NewDB ...
func NewDB(dbFilePath string) *gorm.DB {
	isExist := xfile.IsExist(dbFilePath)
	db, err := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{
		Logger: gormLogger.Discard,
	})
	if err != nil {
		panic("创建数据库失败")
	}
	if !isExist {
		// db 文件在 open 前不存在时, 需要创建表
		db.AutoMigrate(&LoginInfo{})
	}
	return db
}
