package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/guonaihong/gout"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
)

const PKGUrl = "https://unpkg.com"

type PageData struct {
	FileName       string `json:"filename"`
	PackageName    string `json:"packageName"`
	PackageVersion string `json:"packageVersion"`
	Target         struct {
		Path    string            `json:"path"`
		Type    string            `json:"type"`
		Details map[string]Detail `json:"details"`
	} `json:"target"`
}

type Detail struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

func ParsePage(name, version, filename string) (PageData, error) {
	url := fmt.Sprintf("%s/browse/%s@%s%s", PKGUrl, name, version, filename)
	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	html := ""
	var code int
	pageData := PageData{}
	err := gout.GET(url).BindBody(&html).Code(&code).Do()
	if err != nil {
		return pageData, err
	}
	if code == 404 {
		return pageData, errors.New("404 error " + url)
	}
	jsonRule := regexp.MustCompile(`window.__DATA__ =(.+?)</script>`)
	result := jsonRule.FindAllStringSubmatch(html, 1)
	if len(result) != 1 || len(result[0]) != 2 {
		return pageData, errors.New("解析网页数据失败,url=" + url)
	}

	err = json.Unmarshal([]byte(result[0][1]), &pageData)
	if err != nil {
		return pageData, err
	}
	return pageData, nil
}

func DownloadFile(name, version, filename string) error {
	log.Println(filename)
	url := fmt.Sprintf("%s/%s@%s%s", PKGUrl, name, version, filename)
	raw := make([]byte, 0)
	var code int
	err := gout.GET(url).BindBody(&raw).Code(&code).Do()
	if err != nil {
		return err
	}
	if code == 404 {
		return errors.New("404 error " + url)
	}
	fmt.Println(code)
	savePath := fmt.Sprintf("./%s@%s%s", name, version, filename)

	// 父目录
	_, err = os.Stat(path.Join(savePath, ".."))
	if err != nil {
		err = os.MkdirAll(path.Join(savePath, ".."), os.ModeDir)
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(savePath, raw, 0766)
}

func DownloadDir(name, version, filename string) error {
	pageData, err := ParsePage(name, version, filename)
	if err != nil {
		return err
	}

	// 强迫症排下序
	dirArr := make([]string, 0)
	fileArr := make([]string, 0)
	for key, item := range pageData.Target.Details {
		if item.Type == "directory" {
			dirArr = append(dirArr, key)
		}
		if item.Type == "file" && !strings.HasSuffix(item.Path, ".DS_Store") {
			fileArr = append(fileArr, key)
		}
	}
	sort.Strings(dirArr)
	sort.Strings(fileArr)

	for _, key := range dirArr {
		err = DownloadDir(name, version, pageData.Target.Details[key].Path)
		if err != nil {
			return err
		}
	}

	for _, key := range fileArr {
		err = DownloadFile(name, version, pageData.Target.Details[key].Path)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	info := ""
	filename := "/"

	flag.StringVar(&info, "i", info, "package info. vue@2.6.11")
	flag.StringVar(&filename, "f", filename, "download file or dir. /dist/ or /dist/vue.min.js")
	flag.Parse()

	info = strings.TrimSpace(info)

	arr := strings.Split(info, "@")
	if info == "" || len(arr) != 2 || arr[0] == "" || arr[1] == "" {
		flag.Usage()
		os.Exit(0)
	}

	var err error
	if strings.HasSuffix(filename, "/") {
		err = DownloadDir(arr[0], arr[1], filename)
	} else {
		err = DownloadFile(arr[0], arr[1], filename)
	}
	if err != nil {
		log.Fatal(err)
	}
}
