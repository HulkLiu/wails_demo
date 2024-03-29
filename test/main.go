package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileTypes 包含我们要查找的文件类型关键字
type FileTypes struct {
	Number       string
	Title        string
	Img          []string
	VideoFormal  []string
	VideoPreview []string
	VideoShort   []string
}

// scanDirectory 遍历指定目录，寻找包含关键字的文件，并返回文件名
func scanDirectory(dir string) (map[string]FileTypes, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	results := make(map[string]FileTypes)
	for _, file := range files {
		if file.IsDir() {
			subDir := filepath.Join(dir, file.Name())
			subFiles, err := os.ReadDir(subDir)
			if err != nil {
				return nil, err
			}
			fileTypes := FileTypes{}
			for _, subFile := range subFiles {
				filename := subFile.Name()
				if strings.Contains(filename, "jpg") {
					fileTypes.Img = append(fileTypes.Img, filename)
				} else if strings.Contains(filename, "png") {
					fileTypes.Img = append(fileTypes.Img, filename)
				} else if strings.Contains(filename, "VideoFormal") {
					fileTypes.VideoFormal = append(fileTypes.VideoFormal, filename)
				} else if strings.Contains(filename, "VideoPreview") {
					fileTypes.VideoPreview = append(fileTypes.VideoPreview, filename)
				} else if strings.Contains(filename, "VideoShort") {
					fileTypes.VideoShort = append(fileTypes.VideoShort, filename)
				}
				fileTypes.Title = strings.ReplaceAll(filepath.Base(filename), filepath.Ext(filename), "")
			}
			fileTypes.Number = file.Name()
			results[file.Name()] = fileTypes
		}
	}
	return results, nil
}

// printResults 打印结果
func printResults(results map[string]FileTypes) {
	for dir, fileTypes := range results {
		fmt.Printf("%s:\n", dir)
		fmt.Printf("  Number: %s\n", fileTypes.Number)
		fmt.Printf("  Title: %s\n", fileTypes.Title)
		fmt.Printf("  Img: %s\n", strings.Join(fileTypes.Img, ", "))
		fmt.Printf("  VideoFormal: %s\n", strings.Join(fileTypes.VideoFormal, ", "))
		fmt.Printf("  VideoPreview: %s\n", strings.Join(fileTypes.VideoPreview, ", "))
		fmt.Printf("  VideoShort: %s\n", strings.Join(fileTypes.VideoShort, ", "))
	}
}

func main() {

	// 指定需要扫描的目录
	dir := `G:\Videos\`

	// 扫描目录并获取结果
	results, err := scanDirectory(dir)
	if err != nil {
		fmt.Println("Error scanning directory:", err)
		os.Exit(1)
	}

	// 打印结果
	printResults(results)
}
