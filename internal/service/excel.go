package service

import (
	"changeme/internal/define"
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"log"
	"os"
	"path/filepath"
	"time"
	//"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type Exporter struct {
	filePath string
}

func NewExporter(filePath string) *Exporter {
	return &Exporter{
		filePath: filePath,
	}
}

func (e *Exporter) Export(data define.SearchResult) (string, error) {
	// 创建一个新的 Excel 文件
	file := excelize.NewFile()

	os.MkdirAll(filepath.Dir(e.filePath), os.ModePerm)
	// 解析数据结构体并将数据写入 Excel 文件
	json, _ := json.Marshal(data)

	err := e.writeData(file, json)
	if err != nil {
		return "", err
	}

	// 保存 Excel 文件
	err = file.SaveAs(e.filePath)
	if err != nil {
		return "", err
	}

	return e.filePath, nil
}

type SearchResult struct {
	Hits     int64         `json:"Hits"`
	Start    int           `json:"Start"`
	Query    string        `json:"Query"`
	PrevFrom int           `json:"PrevFrom"`
	NextFrom int           `json:"NextFrom"`
	Items    []interface{} `json:"Items"`
}

func (e *Exporter) writeData(file *excelize.File, data []byte) error {

	jsonData := data
	var searchResult SearchResult
	err := json.Unmarshal(jsonData, &searchResult)
	if err != nil {
		//log.Fatal(err)
		return err
	}

	//file := excelize.NewFile()
	sheetName := "Sheet1"
	//sheetName := "videoList"

	// 设置表头
	file.SetCellValue(sheetName, "A1", "Url")
	file.SetCellValue(sheetName, "B1", "Type")
	file.SetCellValue(sheetName, "C1", "Id")
	file.SetCellValue(sheetName, "D1", "Number")
	file.SetCellValue(sheetName, "E1", "Pic")
	file.SetCellValue(sheetName, "F1", "Title")
	file.SetCellValue(sheetName, "G1", "VideoFormal")
	file.SetCellValue(sheetName, "H1", "VideoPreview")
	file.SetCellValue(sheetName, "I1", "VideoShort")
	file.SetCellValue(sheetName, "J1", "ShortLocal")
	file.SetCellValue(sheetName, "K1", "PreviewLocal")
	file.SetCellValue(sheetName, "L1", "FormalLocal")
	file.SetCellValue(sheetName, "M1", "ImgLocal")

	// 写入数据
	for i, item := range searchResult.Items {
		itemMap := item.(map[string]interface{})
		payload := itemMap["Payload"].(map[string]interface{})

		row := i + 2
		file.SetCellValue(sheetName, fmt.Sprintf("A%d", row), itemMap["Url"])
		file.SetCellValue(sheetName, fmt.Sprintf("B%d", row), itemMap["Type"])
		file.SetCellValue(sheetName, fmt.Sprintf("C%d", row), itemMap["Id"])
		file.SetCellValue(sheetName, fmt.Sprintf("D%d", row), payload["Number"])
		file.SetCellValue(sheetName, fmt.Sprintf("E%d", row), payload["Pic"])
		file.SetCellValue(sheetName, fmt.Sprintf("F%d", row), payload["Title"])
		file.SetCellValue(sheetName, fmt.Sprintf("G%d", row), payload["VideoFormal"])
		file.SetCellValue(sheetName, fmt.Sprintf("H%d", row), payload["VideoPreview"])
		file.SetCellValue(sheetName, fmt.Sprintf("I%d", row), payload["VideoShort"])

		file.SetCellValue(sheetName, fmt.Sprintf("J%d", row), payload["ShortLocal"])
		file.SetCellValue(sheetName, fmt.Sprintf("K%d", row), payload["PreviewLocal"])
		file.SetCellValue(sheetName, fmt.Sprintf("L%d", row), payload["FormalLocal"])
		file.SetCellValue(sheetName, fmt.Sprintf("M%d", row), payload["ImgLocal"])
	}

	return nil
}

func test() {
	// 创建导出器实例
	tmpTime := fmt.Sprintf("%v", time.Now().Unix())
	exporter := NewExporter("./data/excel/videoList_ " + tmpTime + ".xlsx")

	// 示例数据
	data := define.SearchResult{}

	// 导出数据到 Excel 文件
	excelPath, err := exporter.Export(data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Excel 文件已导出到:", excelPath)
}
