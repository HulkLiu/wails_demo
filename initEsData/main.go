package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/evercyan/brick/xfile"
	"github.com/olivere/elastic/v7"
	"gopkg.in/yaml.v3"
	"log"
	"reflect"
	"regexp"
	"time"
)

type PayloadJson struct {
	Title        string `json:"Title"`
	VideoShort   string `json:"VideoShort"`
	Pic          string `json:"Pic"`
	Url          string `json:"Url"`
	VideoPreview string `json:"VideoPreview"`
	VideoFormal  string `json:"VideoFormal"`
	Number       string `json:"Number"`
	ShortLocal   string `json:"ShortLocal"`
	PreviewLocal string `json:"PreviewLocal"`
	FormalLocal  string `json:"FormalLocal"`
	ImgLocal     string `json:"ImgLocal"`
	CreateTime   string `json:"CreateTime"`
	UserInput    string `json:"UserInput"`
	VideoType    string `json:"VideoType"`
	Description  string `json:"Description"`
}

type ItemJson struct {
	Url     string      `json:"Url"`
	Type    string      `json:"Type"`
	Id      string      `json:"Id"`
	Payload PayloadJson `json:"Payload"`
}
type SearchResult struct {
	Hits     int64
	Start    int
	Query    string
	PrevFrom int
	NextFrom int
	Items    []interface{}
}
type Item struct {
	Url     string
	Type    string
	Id      string
	Payload interface{}
}

var (
	VideoDefaultJson = "./VideoDefaultJson.conf"
	err              error
	ElasticIndex     = "vjshi10" // vjshi10  test10
	client           *elastic.Client
	size             = 10000 // 查询数量 默认10
)

func main() {
	err = initES()
	if err != nil {
		log.Printf("err:%v", err)
	}
	if checkIndex(ElasticIndex) {
		err = writeDefJson(VideoDefaultJson)
		if err != nil {
			log.Printf("err:%v", err)
		}
	} else {
		err = readDefJson(VideoDefaultJson)
		if err != nil {
			log.Printf("err:%v", err)

		}
	}

}

func checkIndex(index string) bool {

	// 检查索引是否存在
	exists, err := client.IndexExists(index).Do(context.Background())
	if err != nil {
		fmt.Println("Failed to check index existence:", err)
		return false
	}

	if exists {
		fmt.Println("ElasticIndex exists")
		return true
	} else {
		fmt.Println("ElasticIndex does not exist")
		return false

	}
}

func initES() error {
	client, err = elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

func writeDefJson(path string) error {
	resp, err := client.Search(ElasticIndex).
		Query(elastic.NewQueryStringQuery(rewriteQueryString(ElasticIndex))).
		Size(size).
		Do(context.Background())
	if err != nil {
		return err
	}
	data := resp.Each(reflect.TypeOf(Item{})) // 可换 ItemJson

	b, _ := yaml.Marshal(data)
	if err := xfile.Write(path, string(b)); err != nil {
		return err
	}
	return nil
}
func readDefJson(path string) error {

	if !xfile.IsExist(path) {
		return errors.New("defEsJson 文件不存在")
	}

	var list []ItemJson
	if err := yaml.Unmarshal([]byte(xfile.Read(path)), &list); err != nil {
		if err != nil {
			return err
		}
	}
	for _, item := range list {

		indexServer := client.Index().
			Index(ElasticIndex).
			BodyJson(item)
		indexServer.Id(item.Id)

		_, err = indexServer.Do(context.Background())

		if err != nil {
			return err
		}
	}
	time.Sleep(time.Millisecond)
	// 刷新ES中的数据
	_, err = client.Refresh().Index(ElasticIndex).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
func rewriteQueryString(q string) string {
	re := regexp.MustCompile(`([A-Z][a-z]*):`)
	return re.ReplaceAllString(q, "Payload.$1:")
}
