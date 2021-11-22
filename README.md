# go-unpkg-downloader

unpkg.com 批量下载工具，[releases](https://github.com/TMaize/go-unpkg-downloader/releases)

```
unpkg download tool

Usage:
  go-unpkg-downloader pkg [flags]

Flags:
  -d, --dist string   download save path (default "./dist")
  -h, --help          help for go-unpkg-downloader

```

```sh
go-unpkg-downloader vue
go-unpkg-downloader vue@2.6.11
go-unpkg-downloader vue@2.6.11/types
go-unpkg-downloader https://unpkg.com/browse/vue@2.6.11/types/
go-unpkg-downloader @koa/router
go-unpkg-downloader https://unpkg.com/@koa/router@10.1.1/package.json
```
