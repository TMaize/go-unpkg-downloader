package main

import "testing"

func TestParse(t *testing.T) {
	data, err := ParsePage("vue", "2.6.11", "/")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data.PackageName, data.PackageVersion)
	for _, v := range data.Target.Details {
		t.Log(v.Type, v.Path)
	}
}

func TestDownloadFile(t *testing.T) {
	err := DownloadFile("vue", "2.6.11", "/dist/vue.min.js")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDownloadDir(t *testing.T) {
	err := DownloadDir("vue", "2.6.11", "/src/core")
	if err != nil {
		t.Fatal(err)
	}
}
