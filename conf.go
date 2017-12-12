package main

import "fmt"

func GetCommand(moduleName string) string {

	return fmt.Sprintf("echo \"depping...,please wait...\";export http_proxy=cow.98.cn:7777;"+
		"export https_proxy=cow.98.cn:7777;"+
		"export GOPATH=`pwd`;"+
		"export PATH=`pwd`/bin:$PATH;"+
		"go get -v -u github.com/golang/dep/cmd/dep;"+
		"pushd src/%s;"+
		"dep init && dep ensure;"+
		"popd", moduleName)
}
