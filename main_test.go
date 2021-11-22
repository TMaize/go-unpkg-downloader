package main

import (
	"fmt"
	"testing"
)

func TestGetStdUrl(t *testing.T) {
	fmt.Println(GetStdUrl("vue"))
	fmt.Println(GetStdUrl("vue@2.6.14"))
	fmt.Println(GetStdUrl("vue@2.6.11/types"))
	fmt.Println(GetStdUrl("https://unpkg.com/browse/vue@2.6.11/types/"))

	fmt.Println(GetStdUrl("@koa/router"))
	fmt.Println(GetStdUrl("https://unpkg.com/@koa/router@10.1.1/package.json"))
}

func TestParsePage(t *testing.T) {
	data := ParsePage("https://unpkg.com/browse/vue@2.6.11/types/")
	fmt.Println(data.PackageName, data.PackageVersion, data.FileName, data.isFile())

	data = ParsePage("https://unpkg.com/browse/@koa/router@10.1.1/package.json")
	fmt.Println(data.PackageName, data.PackageVersion, data.FileName, data.isFile())
}
