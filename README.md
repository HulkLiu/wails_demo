# videoTools

## 基于Golang Wails 开发的视频管理工具 

## 技术栈

后端：golang

数据库：sqlite

前端：vue3、element-plus、Naive ui

## 安装

1. 安装 golang 1.18 + ，下载地址： https://golang.google.cn/dl/
2. 安装 node 15 + ，下载地址：https://nodejs.org/en/
3. 安装 xcode 命令行 `xcode-select --install`
4. 安装 wails `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

## 功能模块

​	用户设置
​	首页
​	视频模板
​	视频列表
​	任务管理
​	设置

## 构建

先初始化 ESData 并开启 docker 
docker run -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" docker.elastic.co/elasticsearch/elasticsearch:7.4.2

## 软件截图

