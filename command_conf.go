package main

import "fmt"

func InitGoEnvironmentCommand() string {

	return fmt.Sprint("export http_proxy=cow.98.cn:7777;" +
		"export https_proxy=cow.98.cn:7777;" +
		"export GOPATH=`pwd`;" +
		"export PATH=`pwd`/bin:$PATH")
}

func DownloadDepCommand() string {
	return fmt.Sprint("go get -v -u github.com/golang/dep/cmd/dep")
}

func DepTaskCommand(folderName string) string {
	return fmt.Sprintf("pushd src/%s;dep init;dep ensure;popd", folderName)
}

func BuildTaskCommand(moduleName string) string {
	return fmt.Sprintf("go test $(go list %s... | grep -v vendor|tr \"\\n\" \" \");"+
		"go build -o output/%s %s/main;cp -f Dockerfile output/;echo \"result=\"$?", moduleName, moduleName, moduleName)
}
