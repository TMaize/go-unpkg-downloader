package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/guonaihong/gout"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type PageData struct {
	FileName       string `json:"filename"`
	PackageName    string `json:"packageName"`
	PackageVersion string `json:"packageVersion"`
	Target         Target `json:"target"`
}

type Target struct {
	Path    string                 `json:"path"`
	Type    string                 `json:"type"`
	Details map[string]interface{} `json:"details"`
}

func (data PageData) isFile() bool {
	return data.Target.Type == "file"
}

func (data PageData) getBrowseUrl() string {
	return "https://unpkg.com/browse/" + data.PackageName + "@" + data.PackageVersion + data.FileName
}

func (data PageData) getDownloadUrl() string {
	return "https://unpkg.com/" + data.PackageName + "@" + data.PackageVersion + data.FileName
}

func (data PageData) getChildren() []PageData {
	list := make([]PageData, 0)
	if data.Target.Details == nil {
		return list
	}
	keys := make([]string, 0)
	sort.Strings(keys)
	for key := range data.Target.Details {
		keys = append(keys, key)
	}

	for _, key := range keys {

		var item = data.Target.Details[key].(map[string]interface{})
		var fileType = item["type"].(string)
		var filePath = item["path"].(string)
		var fileName = item["path"].(string)
		if fileType == "directory" {
			fileName = fileName + "/"
		}
		list = append(list, PageData{
			FileName:       fileName,
			PackageName:    data.PackageName,
			PackageVersion: data.PackageVersion,
			Target:         Target{Path: filePath, Type: fileType, Details: nil},
		})
	}
	return list
}

//---------------------------------------------

func ParsePage(url string) PageData {
	html := ""
	var code int
	pageData := PageData{}
	err := gout.GET(url).BindBody(&html).Code(&code).Do()
	if err != nil {
		panic(err)
	}
	jsonRule := regexp.MustCompile(`window.__DATA__ =(.+?)</script>`)
	result := jsonRule.FindAllStringSubmatch(html, 1)
	if len(result) != 1 || len(result[0]) != 2 {
		panic("not found __DATA__ at " + url)
	}
	err = json.Unmarshal([]byte(result[0][1]), &pageData)
	if err != nil {
		panic("unmarshal __DATA__ error")
	}
	return pageData
}

func DownloadDir(browseUrl, dest string) {
	data := ParsePage(browseUrl)
	children := data.getChildren()
	for _, item := range children {
		if item.isFile() {
			DownloadFile(item.getDownloadUrl(), dest)
		} else {
			DownloadDir(item.getBrowseUrl(), dest)
		}
	}
}

func DownloadFile(downloadUrl, dest string) {
	parse, err := url.Parse(downloadUrl)
	if err != nil {
		panic(err)
	}
	dest, err = filepath.Abs(dest)
	if err != nil {
		panic(err)
	}
	savePath := filepath.Join(dest, parse.Path)

	fmt.Printf("%s\n%s\n\n", downloadUrl, strings.ReplaceAll(savePath, "\\", "/"))

	raw := make([]byte, 0)
	var code int
	err = gout.GET(downloadUrl).BindBody(&raw).Code(&code).Do()
	if err != nil {
		panic(err)
	}
	if code != 200 {
		panic(fmt.Sprint(code, "error", downloadUrl))
	}

	_, err = os.Stat(filepath.Join(savePath, ".."))
	if err != nil {
		err = os.MkdirAll(filepath.Join(savePath, ".."), os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	err = ioutil.WriteFile(savePath, raw, 0766)
	if err != nil {
		panic(err)
	}
}

func GetStdUrl(pkg string) string {
	pkg = strings.Replace(pkg, "http:", "", 1)
	pkg = strings.Replace(pkg, "https:", "", 1)
	pkg = strings.Replace(pkg, "unpkg.com/browse", "", 1)
	pkg = strings.Replace(pkg, "unpkg.com", "", 1)
	reg := regexp.MustCompile("/{2,}")
	pkg = string(reg.ReplaceAll([]byte(pkg), []byte("/")))
	if strings.HasPrefix(pkg, "/") {
		pkg = pkg[1:]
	}

	browseUrl := "https://unpkg.com/browse/" + pkg

	var statusCode int
	err := gout.HEAD(browseUrl).Code(&statusCode).Do()
	if err != nil {
		panic(err)
	}

	if statusCode != 200 && !strings.HasSuffix(browseUrl, "/") {
		browseUrl = browseUrl + "/"
		err = gout.HEAD(browseUrl).Code(&statusCode).Do()
		if err != nil {
			panic(err)
		}

		if statusCode != 200 {
			panic(fmt.Sprintf("can't access %s %d", browseUrl, statusCode))
		}
	}
	return browseUrl
}

func main() {
	var dist = "./dist"

	command := &cobra.Command{
		Use:   "go-unpkg-downloader pkg [flags]",
		Short: "unpkg download tool",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("require pkg address")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			pkg := args[0]
			data := ParsePage(GetStdUrl(pkg))

			fmt.Printf("%s@%s\n\n", data.PackageName, data.PackageVersion)

			if data.isFile() {
				DownloadFile(data.getDownloadUrl(), dist)
			} else {
				DownloadDir(data.getBrowseUrl(), dist)
			}
		},
	}
	command.Flags().StringVarP(&dist, "dist", "d", dist, "download save path")

	if err := command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
