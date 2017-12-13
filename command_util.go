package main

import (
	"os/exec"
	"log"
)

func RunCommand(command string) string {
	output, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		log.Fatalf("do command %s failed %v,stderr=%s", command, err, string(err.(*exec.ExitError).Stderr))
	}

	return string(output)
}
