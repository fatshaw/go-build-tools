package main

import "fmt"

func GetCommand(moduleName string) string {

	return fmt.Sprintf("export http_proxy=192.168.18.80:7777;"+
		"export https_proxy=192.168.18.80:7777;"+
		"export GOPATH=`pwd`;"+
		"export PATH=`pwd`/bin:$PATH;"+
		"go get -v -u github.com/golang/dep/cmd/dep"+
		"pushd src/%s;"+
		"dep init && dep ensure;"+
		"popd", moduleName)
}
