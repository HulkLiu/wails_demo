package internal

import (
	"changeme/config"
	"changeme/internal/define"
	"changeme/internal/service"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func (a *App) VideoList() define.H {

	data, err := a.Vm.VideoList()
	//data, err := service.DefaultList()

	if err != nil {
		log.Printf("%v", err)
		return define.M{
			"code": -1,
			"msg":  "ERROR:" + err.Error(),
		}
	}
	return define.M{
		"code": 200,
		"data": data,
	}
}

func (a *App) ExportVideoList(form string) define.H {
	data, err := a.Vm.VideoList2(form)
	tmpTime := fmt.Sprintf("%v", time.Now().Unix())

	nowPath, _ := os.Getwd()
	excelDir := nowPath + string(os.PathSeparator) + config.ExcelDir

	exporter := service.NewExporter(excelDir + "videoList_ " + tmpTime + ".xlsx")

	// 导出数据到 Excel 文件
	excelPath, err := exporter.Export(data)
	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stat(filepath.Dir(excelPath))
	if err == nil {
		cmd := exec.Command("explorer", filepath.Dir(excelPath))
		err = cmd.Start()
		if err != nil {
			//return err
			log.Printf("err:%v", err)
		}
	} else if os.IsNotExist(err) {
		cmd := exec.Command("explorer", "C:\\")
		err = cmd.Start()
		if err != nil {
			//return err
			log.Printf("err:%v", err)
		}

	} else {
		return define.M{
			"code": -1,
			"msg":  fmt.Sprintf("Excel 导出失败，err: %v", err),
			"data": excelPath,
		}
	}

	fmt.Println("Excel 文件已导出到:", excelPath)
	return define.M{
		"code": 200,
		"msg":  fmt.Sprintf("Excel 文件已导出到: %v", excelPath),
		"data": excelPath,
	}
}

func (a *App) VideoManage(form string) define.H {
	log.Printf("VideoManage :form:%v ,Type:%T", form, form)

	data, err := a.Vm.VideoList2(form)

	if err != nil {
		log.Printf("%v", err)
		return define.M{
			"code": -1,
			"msg":  "ERROR:" + err.Error(),
		}
	}
	return define.M{
		"code": 200,
		"data": data,
	}
}

func (a *App) VideoCreate(form interface{}) define.H {
	log.Printf("VideoCreate :form:%v ,Type:%T", form, form)

	jsonStr, err := json.Marshal(form)
	if err != nil {
		log.Printf("err:%v", err)

		// 处理错误
	}

	var myStruct define.ItemJson
	err = json.Unmarshal(jsonStr, &myStruct)
	if err != nil {
		log.Printf("err:%v", err)
	}
	log.Printf("myStruct:%+v ,Type:%T", myStruct, myStruct)
	msg := "新建成功"
	if myStruct.Id != "" {
		msg = "修改成功"
	}

	err = a.Vm.CreateVideo(myStruct)
	if err != nil {
		return define.M{
			"code": -1,
			"msg":  "ERROR : " + err.Error(),
		}
	}
	return define.M{
		"code": 200,
		"msg":  msg,
	}
}

func (a *App) VideoDelete(form interface{}) define.H {
	log.Printf("delete: form:%v ,Type:%T", form, form)

	jsonStr, err := json.Marshal(form)
	if err != nil {
		log.Printf("err:%v", err)
	}

	var myStruct define.ItemJson
	err = json.Unmarshal(jsonStr, &myStruct)
	if err != nil {
		log.Printf("err:%v", err)
	}
	log.Printf("myStruct:%+v ,Type:%T", myStruct, myStruct)
	msg := "删除成功"

	err = a.Vm.DeleteVideo(myStruct)
	if err != nil {
		return define.M{
			"code": -1,
			"msg":  "ERROR : " + err.Error(),
		}
	}
	return define.M{
		"code": 200,
		"msg":  msg,
	}
}

// SearchKey define.D
func (a *App) SearchKey(form define.M) define.H {
	log.Printf("form:%+v ,Type:%T", form, form)
	//return define.M{}

	data, err := a.Vm.SearchKey(form)

	if err != nil {
		return define.M{
			"code": -1,
			"msg":  "ERROR:" + err.Error(),
		}
	}
	//log.Printf("data:%+v", data)

	return define.M{
		"code": 200,
		"data": data,
	}
}

// OpenFolder  打开本地文件夹
func (a *App) OpenFolder(form string) define.H {
	//log.Printf("form:%+v ,Type:%T", form, form)
	//return define.M{}

	err := service.OpenFolder(form)

	if err != nil {
		return define.M{
			"code": -1,
			"msg":  "ERROR:" + err.Error(),
		}
	}
	return define.M{
		"code": 200,
		"data": "",
	}
}

// ConfigEdit  更新配置信息
func (a *App) ConfigEdit(form define.M) define.H {
	//log.Printf("form:%+v ,Type:%T", form, form)
	//return define.M{}

	data, err := service.ConfigEdit(form)

	if err != nil {
		return define.M{
			"code": -1,
			"msg":  "ERROR:" + err.Error(),
		}
	}
	log.Printf("data:%+v", data)

	return define.M{
		"code": 200,
		"msg":  "配置成功",
		"data": data,
	}
}

func (a *App) VideoGetNumber(form interface{}) define.H {
	log.Printf("form:%+v ,Type:%T", form, form)
	//return define.M{}
	number := ""
	//根据 number 爬取资源
	if number != "" {

	} else {
		//分配 number
		data, err := a.Vm.GetVideoNumber()
		if err != nil {
			return define.M{
				"code": -1,
				"msg":  "ERROR:" + err.Error(),
			}
		}
		//log.Printf("data:%+v", data)

		return define.M{
			"code": 200,
			"data": data,
		}
	}

	return nil
}
