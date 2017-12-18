package main

import (
	"fmt"
	"os"
	"log"
)


func AllInOne() {

	buildTask()

	chdirToSourceFolder()

	DockerTask(fmt.Sprintf("%s/%s:%s", DOCKERIMAGEPREFIX, os.Args[2], os.Getenv(CIBUILDREF)))

}

func chdirToSourceFolder() {
	if err := os.Chdir(fmt.Sprintf("%s/%s", SOURCEFOLDER, os.Args[2])); err != nil {
		log.Fatalf("chdir failed:%v", err)
	}
}
func buildTask() {
	command := fmt.Sprintf("%s;%s", BeforeScript(), BuildTaskCommand(os.Args[2]))
	output := RunCommand(command)
	log.Printf("AllInOne=%s,output=%s\n", command, output)
}
