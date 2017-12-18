package main

import "fmt"

func InitGoEnvironmentCommand() string {

	return fmt.Sprint("export GOPATH=`pwd`;" +
		"export PATH=`pwd`/bin:$PATH")
}

func DownloadDepCommand() string {
	return fmt.Sprint("go get -v -u github.com/golang/dep/cmd/dep")
}

func DepTaskCommand(folderName string) string {
	return fmt.Sprintf("pushd src/%s;dep ensure;popd", folderName)
}

func BuildTaskCommand(moduleName string) string {
	return fmt.Sprintf("go test $(go list ytx/futures/go/%s... | grep -v vendor|tr \"\\n\" \" \");"+
		"go build -o output/%s ;cp -f Dockerfile output/;echo \"result=\"$?", moduleName, moduleName)
}

func BeforeScript() string {
	return fmt.Sprint("export GOPATH=/root/go;" +
		"mkdir -p $GOPATH/src/ytx/futures/go;" +
		"cd $GOPATH/src/ytx/futures/go" +
		"rm -fr $CI_PROJECT_NAME;" +
		"cp -fr /root/$CI_PROJECT_DIR .;" +
		"cd $CI_PROJECT_NAME" +
		"go get -u github.com/golang/dep/cmd/dep;" +
		"$GOPATH/bin/dep ensure -update")
}
