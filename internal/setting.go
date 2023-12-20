package internal

import (
	"changeme/config"
	"changeme/internal/define"
	"changeme/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/evercyan/brick/xfile"
	"gopkg.in/yaml.v3"
	"log"
)

func (a *App) GetSettingList() define.H {

	data := a.set.Data
	log.Printf("GetSettingList :%v", data)
	log.Printf("GetSettingList -> a.Vm.ElasticIndex :%v", a.Vm.ElasticIndex)
	return define.M{
		"code": 200,
		"data": data,
	}
}

type SettingManage struct {
	Path string
	Data *SettingData
}
type SettingData struct {
	Bak BakSetting `json:"Bak"`
}
type BakSetting struct {
	EsIndex     string `json:"EsIndex"`
	LocalPath   string `json:"LocalPath"`
	ExcelDir    string `json:"ExcelDir"`
	ContainerID string `json:"ContainerID"`
}

func NewSet() SettingManage {
	var data SettingManage
	cfgPath := utils.GetCfgPath()

	data = SettingManage{
		Path: fmt.Sprintf(config.SettingFile, cfgPath),
		//Data: &SettingData{
		//	Bak: BakSetting{
		//		EsIndex:     config.EsIndex,
		//		LocalPath:   config.HomeInfoDir,
		//		ExcelDir:    config.ExcelDir,
		//		ContainerID: config.ContainerID,
		//	},
		//},
	}
	data.readJson()
	return data
}

func (s *SettingManage) readJson() {
	if xfile.IsExist(s.Path) {

		list := &SettingData{}
		if err := yaml.Unmarshal([]byte(xfile.Read(s.Path)), &list); err != nil {
			if err != nil {
				fmt.Printf("SettingManage err:%v", err)
			}
		} else {
			s.Data = list
		}
	} else {
		err := s.putJson()
		if err != nil {
			fmt.Printf("SettingManage err:%v", err)
		}
	}

}
func (s *SettingManage) putJson() error {
	var data *SettingData
	data = &SettingData{
		Bak: BakSetting{
			EsIndex:     config.EsIndex,
			LocalPath:   config.HomeInfoDir,
			ExcelDir:    config.ExcelDir,
			ContainerID: config.ContainerID,
		},
	}
	b, _ := yaml.Marshal(data)
	//log.Printf("putJson -> b:%v", b)
	if err := xfile.Write(s.Path, string(b)); err != nil {
		log.Printf("Set CfgDirInfo Write err: %v", err)
		return err
	}
	s.Data = data
	//log.Printf("putJson -> s.Data:%v", s.Data)

	return nil

}

func (s *SettingManage) UpdateSetting(form SettingData) error {
	//写文件
	b, _ := yaml.Marshal(form)
	//log.Printf("putJson -> b:%v", b)
	if err := xfile.Write(s.Path, string(b)); err != nil {
		log.Printf("Set CfgDirInfo Write err: %v", err)
		return err
	}

	//更新 SettingData
	s.Data = &form

	return nil
}

func (a *App) UpdateSettingList(form interface{}) define.H {
	log.Printf("VideoCreate :form:%v ,Type:%T", form, form)

	jsonStr, err := json.Marshal(form)
	if err != nil {
		log.Printf("err:%v", err)
	}

	var myStruct SettingData
	err = json.Unmarshal(jsonStr, &myStruct)
	if err != nil {
		log.Printf("err:%v", err)
	}
	log.Printf("myStruct:%+v ,Type:%T", myStruct, myStruct)

	err = a.set.UpdateSetting(myStruct)
	a.Vm.ElasticIndex = a.set.Data.Bak.EsIndex
	if err != nil {
		return define.M{
			"code": -1,
			"msg":  "ERROR : " + err.Error(),
		}
	}
	return define.M{
		"code": 200,
		"msg":  "更新成功",
	}
}
