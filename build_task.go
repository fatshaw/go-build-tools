package main

import (
	"fmt"
	"log"
)

func BuildTask(moduleName string) {
	output := RunCommand(fmt.Sprintf("%s", []string{"-c", InitGoEnvironmentCommand(), BuildTaskCommand(moduleName)}))

	log.Printf("buildTask=%soutput=%s\n", []string{"-c", InitGoEnvironmentCommand(), BuildTaskCommand(moduleName)}, output)

}
