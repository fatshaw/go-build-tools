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

	command := []string{"-c", InitGoEnvironmentCommand(), DownloadDepCommand()}
	for _, f := range files {
		// ignore github.com source folder
		if strings.Contains(f.Name(), "github.com") {
			continue
		}

		log.Printf("dep for folder = %s\n", f.Name())
		command = append(command, DepTaskCommand(f.Name()))
	}

	output := RunCommand(fmt.Sprintf("%s", command))
	log.Printf("depTask=%s\noutput=%s\n", fmt.Sprintf("%s", command), output)

}