package service

import (
	"fmt"
	"github.com/evercyan/brick/xfile"
	"gopkg.in/yaml.v3"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
)

type DirInfo struct {
	Filed []string
	Count []int64
	Size  []string
}
type FileStat struct {
	Ext   string
	Count int64
	Size  int64
}

func formatSize(bytes int64) string {
	units := [7]string{" ", "K", "M", "G", "T", "P", "E"}
	if bytes < 1024 {
		return fmt.Sprintf("%v B", bytes)
	}
	z := 0
	v := float64(bytes)
	for v > 1024.0 {
		z++
		v /= 1024.0
	}
	return fmt.Sprintf("%.2f %siB", v, units[z])
}
func formatSize_MB(bytes int64) string {

	v := float64(bytes) / (1024 * 1024)
	return fmt.Sprintf("%.1f", v)
}

func PutDirInfoDesk(data DirInfo, dirPath string) error {
	b, _ := yaml.Marshal(data)
	if err := xfile.Write(dirPath, string(b)); err != nil {
		log.Printf("Set CfgDirInfo Write err: %v", err)
		return err
	}

	return nil
}

func GetDirInfo(directory string) DirInfo {
	var dirInfo DirInfo
	fileStats := make(map[string]FileStat)

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				fmt.Printf("Error: Permission denied to access %s\n", path)
				return nil
			} else {
				fmt.Printf("Error: %s\n", err)
				return err
			}
		}

		if !info.IsDir() {
			ext := filepath.Ext(path)
			size := info.Size()

			if _, ok := fileStats[ext]; !ok {
				fileStats[ext] = FileStat{Ext: ext, Count: 1, Size: size}
			} else {
				fileStats[ext] = FileStat{Ext: ext, Count: fileStats[ext].Count + 1, Size: fileStats[ext].Size + size}
			}
		}

		return nil
	})

	if err != nil {
		return dirInfo
	}

	// 将map转换为slice，方便排序
	stats := make([]FileStat, 0, len(fileStats))
	for _, stat := range fileStats {
		stats = append(stats, stat)
	}

	// 按照Count排序
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Count > stats[j].Count
	})
	var filed []string
	var count []int64
	var size []string
	// 只打印前十个数量的后缀名
	for i, stat := range stats {
		if i >= 10 {
			break
		}

		filed = append(filed, stat.Ext)
		count = append(count, stat.Count)
		size = append(size, formatSize_MB(stat.Size))
	}
	dirInfo = DirInfo{
		Filed: filed,
		Count: count,
		Size:  size,
	}
	return dirInfo
}

type MathData struct {
	max    float64
	min    float64
	avg    float64
	median float64
}

func getBasicData(numbers []float64) MathData {
	//numbers := []float64{1, 3, 5, 7, 9}
	if len(numbers) <= 0 {
		return MathData{}
	}
	m := MathData{
		// 计算最大值
		max: math.Ceil(max(numbers)),
		// 计算最小值
		min: math.Ceil(min(numbers)),
		// 计算平均值
		avg: math.Ceil(avg(numbers)),
		// 计算中位数
		median: math.Ceil(median(numbers)),
	}
	//fmt.Printf("最大值: %v, 最小值: %v, 平均值: %v, 中位数: %v\n", max, min, avg, median)
	return m
}

// 计算最大值
func max(numbers []float64) float64 {
	max := numbers[0]
	for _, num := range numbers {
		if num > max {
			max = num
		}
	}
	return max
}

// 计算最小值
func min(numbers []float64) float64 {
	min := numbers[0]
	for _, num := range numbers {
		if num < min {
			min = num
		}
	}
	return min
}

// 计算平均值
func avg(numbers []float64) float64 {
	sum := 0.0
	for _, num := range numbers {
		sum += num
	}
	return sum / float64(len(numbers))
}

// 计算中位数
func median(numbers []float64) float64 {
	sort.Float64s(numbers)
	n := len(numbers)
	if n%2 == 0 {
		return (numbers[n/2-1] + numbers[n/2]) / 2
	}
	return numbers[n/2]
}
