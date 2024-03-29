package service

import (
	"changeme/config"
	"changeme/internal/define"
	"changeme/internal/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"net"
	"sync"
	"time"

	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

var logger = utils.NewLogger("video")

type VideoManage struct {
	ElasticIndex     string
	Client           *elastic.Client
	Err              error
	PageSize         int
	VideoDefaultJson string
	TotalInfo        VideoTotalInfo
	wg               sync.Mutex
}
type VideoTotalInfo struct {
	Video int64
	Type  map[string]int
}

func checkPort(host string, port int) bool {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
func CmdEsRun() (string, error) {
	//var res string
	cmd := exec.Command("cmd", "/C", "docker start "+config.ContainerID)
	//cmd := exec.Command("docker", "start", config.ContainerID)
	res, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	logger.Printf("%v", "docker start initEsData success")
	return string(res), nil

}

func (v *VideoManage) GetInfo() VideoTotalInfo {
	var total = VideoTotalInfo{}
	resp, err := v.Client.Search(v.ElasticIndex).
		Query(elastic.NewQueryStringQuery(
			rewriteQueryString(v.ElasticIndex))).
		Size(10000).
		Do(context.Background())

	if err != nil {
		return total
	}

	total.Video = resp.TotalHits()
	data := resp.Each(reflect.TypeOf(define.Item{}))

	jsonData, _ := json.Marshal(data)
	var items []define.ItemJson
	_ = json.Unmarshal(jsonData, &items)

	videoTypeCount := make(map[string]int)

	for _, item := range items {
		videoType := item.Payload.VideoType
		videoTypeCount[videoType]++
	}
	total.Type = videoTypeCount

	return total
}

func ConnectionList() ([]*define.Connection, error) {
	nowPath, _ := os.Getwd()
	data, err := ioutil.ReadFile(nowPath + string(os.PathSeparator) + define.ConfigName)
	if errors.Is(err, os.ErrNotExist) {
		return []*define.Connection{{
			Identity: "1",
			Name:     "1",
			Addr:     "1",
			Port:     "1",
			UserName: "1",
			PassWord: "1",
		}}, nil
	}
	conf := new(define.Config)
	err = json.Unmarshal(data, conf)
	if err != nil {
		return nil, err
	}
	return conf.Connections, nil
}

var (
	err error
)

func interfaceToMap(i interface{}) (map[string]interface{}, error) {
	if m, ok := i.(map[string]interface{}); ok {
		return m, nil
	}
	return map[string]interface{}{}, fmt.Errorf("unable to convert interface to string")

}
func interfaceToSlice(i interface{}) ([]interface{}, error) {
	if m, ok := i.([]interface{}); ok {
		return m, nil
	}
	return []interface{}{}, fmt.Errorf("unable to convert interface to string")
}
func interfaceToString(i interface{}) (string, error) {
	if str, ok := i.(string); ok {
		return str, nil
	}
	return "", fmt.Errorf("unable to convert interface to string")
}
func interfaceToInt(i interface{}) (int, error) {
	if num, ok := i.(float64); ok {
		return int(num), nil
	}
	return 0, fmt.Errorf("unable to convert interface to int")
}

func ConfigEdit(form define.M) (define.LoginInfo, error) {
	var users []define.LoginInfo
	var res define.LoginInfo
	log.Printf("form:%+v,T:%T", form, form)
	con, _ := interfaceToMap(form["config"])
	log.Printf("con:%v", con)
	convertUser, _ := interfaceToString(con["username"])
	convertPass, _ := interfaceToString(con["password"])

	u := define.LoginInfo{
		UserName: "admin",
		PassWord: "123456",
	}
	users = append(users, u)
	log.Printf("user:%v,pass:%v", convertUser, convertPass)
	for _, v := range users {
		if v.UserName == convertUser && v.PassWord == convertPass {
			return u, nil
		}
	}

	return res, errors.New("查无此人")
}

func (v *VideoManage) SearchKey(form define.M) (define.SearchResult, error) {
	var result define.SearchResult

	q := form["query"]
	f := form["page"]
	//log.Printf("form[\"query\"]:%+v,T:%T", q, q)
	//log.Printf("form[\"nextFrom\"]:%+v,T:%T", f, f)

	convertQ, _ := interfaceToString(q)
	convertF, _ := interfaceToInt(f)

	if v.Client == nil {
		data, _ := v.DefaultList(convertQ, convertF)
		return data, nil
	}

	//convertF = 10
	result, err := v.getSearchResult(convertQ, convertF)
	if err != nil {
		return result, err
	}
	//log.Printf("result.Start :%v,result.Hits:%v,result.PrevFrom:%v,result.NextFrom:%v,result.Query :%v", result.Start, result.Hits, result.PrevFrom, result.NextFrom, result.Query)
	return result, nil
}

func (v *VideoManage) DefaultList(q string, from int) (define.SearchResult, error) {
	var result define.SearchResult

	nowPath, err := os.Getwd()
	if err != nil {
		return result, err
	}
	data, _ := ioutil.ReadFile(nowPath + string(os.PathSeparator) + v.VideoDefaultJson)
	err = json.Unmarshal(data, &result)

	result.Start = from
	info := result.Items
	//log.Printf("result.Items:%v ,type:%T", result.Items, result.Items)
	l := int64(len(result.Items))
	var indices []int
	if q != "" && q != "vjshi6" {

		for k, item := range result.Items {
			//log.Printf("item:%v ,type:%T", item, item)

			va, ok := item.(map[string]interface{})
			//va, ok := item.(define.M)
			if !ok {
				return define.SearchResult{}, nil
			}

			title := va["Payload"].(map[string]interface{})["Title"].(string)

			if !strings.Contains(title, q) {
				indices = append(indices, k)
			} else {
				l--
			}
		}

		for i := len(indices) - 1; i >= 0; i-- {
			index := indices[i]
			info = append(info[:index], info[index+1:]...)
		}
	}
	result.Items = info
	result.Hits = l

	if result.Start == 0 {
		result.PrevFrom = -1
	} else {
		result.PrevFrom = (result.Start - 1) / v.PageSize * v.PageSize
	}
	result.NextFrom = result.Start + len(result.Items)

	if err != nil {
		return result, err
	}

	return result, nil
}

func (v *VideoManage) VideoList() (define.SearchResult, error) {
	var result define.SearchResult

	from := 0
	q := ""

	if v.Client == nil {
		data, _ := v.DefaultList(q, from)
		return data, nil
	}

	result, err := v.getSearchResult(q, from)
	if err != nil {
		return result, err
	}
	//log.Printf("%v", result)
	return result, nil
	//return []define.AutoGenerated{{}}, nil

}

func (v *VideoManage) VideoList2(keyword string) (define.SearchResult, error) {
	var result define.SearchResult

	from := 0
	q := keyword

	_, err = v.Client.Refresh().Index(v.ElasticIndex).Do(context.Background())
	if err != nil {
		return result, err
	}
	if err != nil {
		data, _ := v.DefaultList(q, from)
		return data, nil
	}
	result, err = v.getSearchResult2(q, from)
	if err != nil {
		return result, err
	}
	return result, nil

}

func (v *VideoManage) getSearchResult2(q string, from int) (define.SearchResult, error) {
	var result define.SearchResult
	if q == "" {
		q = v.ElasticIndex
	}
	result.Query = q

	resp, err := v.Client.
		Search(v.ElasticIndex).
		Query(elastic.NewQueryStringQuery(
			rewriteQueryString(q))).
		Size(10000).
		Do(context.Background())

	if err != nil {
		return result, err
	}

	result.Hits = resp.TotalHits()

	result.Items = resp.Each(reflect.TypeOf(define.Item{}))

	//分组统计
	if false {
		jsonData, _ := json.Marshal(result.Items)
		var items []define.ItemJson
		_ = json.Unmarshal(jsonData, &items)

		videoTypeCount := make(map[string]int)
		log.Printf("GetInfo - len(items):%v", len(items))
		for _, item := range items {
			videoType := item.Payload.VideoType
			videoTypeCount[videoType]++
		}
		log.Printf("GetInfo - videoTypeCount: %v", videoTypeCount)
	}

	log.Printf("result.Start :%v,result.Hits:%v,result.PrevFrom:%v,result.NextFrom:%v,result.Query :%v,result.Items:%v", result.Start, result.Hits, result.PrevFrom, result.NextFrom, result.Query, len(result.Items))

	return result, nil
}

// FileTypes 包含我们要查找的文件类型关键字
type FileTypes struct {
	Number       string
	Title        string
	Img          []string
	VideoFormal  []string
	VideoPreview []string
	VideoShort   []string
}

// ScanDirectory 遍历指定目录，寻找包含关键字的文件，并返回文件名
func (v *VideoManage) ScanDirectory(dir string) (map[string]FileTypes, error) {
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
				fileTypes.Number = file.Name()
				fileTypes.Title = filepath.Base(filename)

			}
			results[file.Name()] = fileTypes
		}
	}
	return results, nil
}

func (v *VideoManage) getSearchResult(q string, from int) (define.SearchResult, error) {
	var result define.SearchResult
	if q == "" {
		q = v.ElasticIndex
	}
	result.Query = q

	resp, err := v.Client.
		Search(v.ElasticIndex).
		Query(elastic.NewQueryStringQuery(
			rewriteQueryString(q))).
		From(from).
		Do(context.Background())

	if err != nil {
		return result, err
	}

	result.Hits = resp.TotalHits()
	result.Start = from
	result.Items = resp.Each(reflect.TypeOf(define.Item{}))
	if result.Start == 0 {
		result.PrevFrom = -1
	} else {
		result.PrevFrom = (result.Start - 1) / v.PageSize * v.PageSize
	}
	result.NextFrom = result.Start + len(result.Items)

	log.Printf("result.Start :%v,result.Hits:%v,result.PrevFrom:%v,result.NextFrom:%v,result.Query :%v", result.Start, result.Hits, result.PrevFrom, result.NextFrom, result.Query)

	return result, nil
}

func rewriteQueryString(q string) string {
	re := regexp.MustCompile(`([A-Z][a-z]*):`)
	return re.ReplaceAllString(q, "Payload.$1:")
}

func OpenFolder(form string) error {
	dir := form

	dir = strings.ReplaceAll(dir, "\\\\", "\\")
	dir = strings.ReplaceAll(dir, "\"", "")
	dir = strings.TrimSpace(dir)
	dir = filepath.Dir(dir)
	log.Println(dir)
	_, err = os.Stat(dir)
	if err == nil {
		cmd := exec.Command("explorer", dir)
		err = cmd.Start()
		if err != nil {
			return err
		}
	} else if os.IsNotExist(err) {
		cmd := exec.Command("explorer", "D:\\")
		err = cmd.Start()
		if err != nil {
			return err
		}

	} else {
		return errors.New("发生错误" + fmt.Sprintf("%s", err))
	}

	return nil
}

func (v *VideoManage) CreateVideo(form define.ItemJson) error {
	//log.Printf("form:%+v ,Type:%T", form, form)
	//return nil

	item := form
	if item.Id == "" {
		item.Id = item.Payload.Number
	}
	if item.Type == "" {
		item.Type = v.ElasticIndex
	}
	log.Printf("item:%+v ", item)

	indexServer := v.Client.Index().
		Index(v.ElasticIndex).
		BodyJson(item)
	indexServer.Id(item.Id)

	_, err = indexServer.Do(context.Background())

	if err != nil {
		return err
	}
	time.Sleep(time.Millisecond)
	// 刷新ES中的数据
	_, err = v.Client.Refresh().Index(v.ElasticIndex).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (v *VideoManage) DeleteVideo(form define.ItemJson) error {
	//log.Printf("form:%+v ,Type:%T", form, form)
	//return nil

	item := form
	if item.Id == "" {
		item.Id = item.Payload.Number
	}
	log.Printf("delete:item:%+v ", item)

	_, err := v.Client.Delete().
		Index(v.ElasticIndex).
		Id(item.Id).
		Do(context.Background())

	if err != nil {
		return err
	}
	// 刷新ES中的数据
	_, err = v.Client.Refresh().Index(v.ElasticIndex).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

var num int

func (v *VideoManage) GetVideoNumber() (define.ItemJson, error) {
	v.wg.Lock()
	defer func() {
		v.wg.Unlock()
	}()
	var res define.ItemJson

	t := "10" + time.Now().Format(config.TimeFormat) + fmt.Sprintf("%v", num)
	num++
	res = define.ItemJson{
		Payload: define.PayloadJson{
			Number: t,
			//ShortLocal: config.HomeInfoDir + "G:\Videos/5874001/红军长征路线图视频素材_5874001_VideoShort.mp4",
			//PreviewLocal: config.HomeInfoDir + "G:\Videos/5874001/红军长征路线图视频素材_5874001_VideoPreview.mp4",
			//FormalLocal: config.HomeInfoDir + "G:\Videos/5874001/红军长征路线图视频素材_5874001_VideoFormal.mp4",
			//ImgLocal: config.HomeInfoDir + "G:\Videos/5874001/红军长征路线图视频素材_5874001.jpg",
		},
	}

	return res, nil
}
