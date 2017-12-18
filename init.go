package main

import (
	"fmt"
	"os"
	"io"
	"text/tabwriter"
	"log"
)

func main() {

	example := [][2]string{
		{
			"go-build-tools setup [modulename]",
			"set up a new project with name=[modulename]",
		},
		{
			"go-build-tools init [vscode|idea]",
			"init go environment for vscode or idea",
		},
		{
			"go-build-tools dep",
			"install the project's dependencies",
		},
		{
			"go-build-tools build [module_name]",
			"build module",
		},
		{
			"go-build-tools pushImage [module_name]",
			"build docker image and push",
		},
		{
			"go-build-tools allinone [module_name]",
			"allinone work",
		},
	}

	if len(os.Args) == 2 && os.Args[1] == "dep" {

		InitDep()

	} else if len(os.Args) == 3 {

		if os.Args[1] == "setup" {
			Setup(os.Args[2])
		}

		if os.Args[1] == "init" {

			if os.Args[2] == "idea" {
				InitIdea()
			}
			if os.Args[2] == "vscode" {
				InitVscode()
			}
		}

		if os.Args[1] == "build" {
			BuildTask(os.Args[2])
		}

		if os.Args[1] == "pushImage" {
			DockerTask(fmt.Sprintf("daocloud.io/baidao/%s:%s", os.Args[2], os.Getenv("CI_BUILD_REF")))
		}

		if os.Args[1] == "allinone" {
			AllInOne()
		}

	} else {
		usage(os.Stdout, example)
	}
}

func usage(w io.Writer, examples [][2] string) {
	fmt.Fprintln(w, "go-dep-tools is a tool to manage go dep and init work")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Usage: \"go-dep-tools [command]\"")
	fmt.Fprintln(w)
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w)
	fmt.Fprintln(w, "Examples:")
	for _, example := range examples {
		fmt.Fprintf(tw, "\t%s\t%s\n", example[0], example[1])
	}
	tw.Flush()
	fmt.Fprintln(w)
}

func AllInOne() {

	command := fmt.Sprintf("%s;%s", BeforeScript(), BuildTaskCommand(os.Args[2]))
	output := RunCommand(command)
	log.Printf("AllInOne=%s,output=%s\n", command, output)

	if err := os.Chdir(fmt.Sprintf("/root/go/src/ytx/futures/go/%s", os.Args[2])); err != nil {
		log.Fatalf("chdir failed:%v", err)
	}

	DockerTask(fmt.Sprintf("daocloud.io/baidao/%s:%s", os.Args[2], os.Getenv("CI_BUILD_REF")))

}
