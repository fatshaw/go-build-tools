package main

import (
	"os"
	"io/ioutil"
	"strings"
	"fmt"
	"log"
)

func InitDep() {

	if _, err := os.Stat("src"); os.IsNotExist(err) {
		log.Fatal("not src folder")
	}

	files, err := ioutil.ReadDir("src")
	if err != nil {
		log.Fatal(err)
	}

	command := fmt.Sprintf("%s;%s", InitGoEnvironmentCommand(), DownloadDepCommand())
	for _, f := range files {
		// ignore github.com source folder
		if strings.Contains(f.Name(), "github.com") {
			continue
		}

		log.Printf("dep for folder = %s\n", f.Name())
		command = fmt.Sprintf("%s;%s", command, DepTaskCommand(f.Name()))
	}

	output := RunCommand(command)
	log.Printf("depTask=%s,output=%s\n", command, output)

}

func BeforeScriptTask() {

	output := RunCommand(fmt.Sprintf("%s", BeforeScript()))
	log.Printf("BeforeScriptTask=%s,output=%s\n", fmt.Sprintf("%s", BeforeScript()), output)

}
