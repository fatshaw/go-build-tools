package main

import (
	"fmt"
	"log"
)

func BuildTask(moduleName string) {
	output := RunCommand(fmt.Sprintf("%s && %s", InitGoEnvironmentCommand(), BuildTaskCommand(moduleName)))

	log.Printf("buildTask=%s,output=%s\n", fmt.Sprintf("%s && %s", InitGoEnvironmentCommand(), BuildTaskCommand(moduleName)), output)

}
