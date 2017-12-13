package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"os/exec"
	"io"
	"text/tabwriter"
	"log"
	"encoding/base64"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"context"
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
			"go-build-tools dep [modulename]",
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
	}

	if len(os.Args) == 2 && os.Args[1] == "dep" {

		initDep()

	} else if len(os.Args) == 3 {

		if os.Args[1] == "setup" {
			setup(os.Args[2])
		}

		if os.Args[1] == "init" {

			if os.Args[2] == "idea" {
				initIdea()
			}
			if os.Args[2] == "vscode" {
				initVscode()
			}
		}

		if os.Args[1] == "build" {
			buildTask(os.Args[2])
		}

		if os.Args[1] == "pushImage" {
			dockerTask(fmt.Sprintf("daocloud.io/baidao/%s:%s", os.Args[2], os.Getenv("CI_BUILD_REF")))
		}

	} else {
		usage(os.Stdout, example)
	}
}

func setup(moduleName string) {
	path := fmt.Sprintf("./src/%s", moduleName)
	mainPath := fmt.Sprintf("%s/main.go", path)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
		ioutil.WriteFile(mainPath, []byte("package main;\nfunc main()\n{\n}\n"), 0755)
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

func initDep() {

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

	output := runCommand(fmt.Sprintf("%s", command))
	log.Printf("depTask=%s\noutput=%s\n", fmt.Sprintf("%s", command), output)

}

func runCommand(command string) string {
	output, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		log.Fatalf("do command %s failed %v", command, err)
	}

	return string(output)

}

func initIdea() {

	bytes, _ := ioutil.ReadFile(".idea/workspace.xml")
	content := string(bytes)

	var perm os.FileMode = 0755
	if !strings.Contains(content, "name=\"GoLibraries") {
		lines := file2lines(".idea/workspace.xml")
		lines[len(lines)-1] = fmt.Sprintf("<component name=\"GoLibraries\">\n<option name=\"urls\">" +
			"<list>\n<option value=\"file://$PROJECT_DIR$\" /></list></option></component>")
		lines = append(lines, "</project>")
		writeContent := make([]byte, 0)
		for _, line := range lines {
			writeContent = append(writeContent, line...)
			writeContent = append(writeContent, "\n"...)
		}
		ioutil.WriteFile(".idea/workspace.xml", writeContent, perm)
	}
}

func initVscode() {
	var perm os.FileMode = 0755

	os.MkdirAll(".vscode", perm)

	dat := make(map[string]string)
	if _, err := os.Stat(".vscode/settings.json"); err == nil {
		// path/to/whatever does exist
		content, _ := ioutil.ReadFile(".vscode/settings.json")
		json.Unmarshal(content, &dat)
	}

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dat["go.gopath"] = dir

	content, _ := json.Marshal(dat)
	ioutil.WriteFile(".vscode/settings.json", content, perm)
}

func file2lines(filePath string) []string {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return lines
}

func buildTask(moduleName string) {
	output := runCommand(fmt.Sprintf("%s", []string{"-c", InitGoEnvironmentCommand(), BuildTask(moduleName)}))

	log.Printf("buildTask=%soutput=%s\n", []string{"-c", InitGoEnvironmentCommand(), BuildTask(moduleName)}, output)

}

func dockerTask(imageName string) {

	defaultHeaders := map[string]string{"User-Agent": "ego-v-0.0.1"}
	cli, _ := client.NewClient("unix:///var/run/docker.sock", "v1.24", nil, defaultHeaders)

	authConfig := types.AuthConfig{
		Username:      "developer@baidao.com",
		Password:      "65T-Tvq-sVc-BDR",
		ServerAddress: "daocloud.io",
	}
	auth, err := cli.RegistryLogin(context.Background(), authConfig)

	if err != nil {
		panic(err)
	}

	log.Printf("login response=%s", auth.Status)

	buildOptions := types.ImageBuildOptions{
		Tags:           []string{imageName},
		Dockerfile:     "simulator/Dockerfile",
		SuppressOutput: false,
		Remove:         true,
		ForceRemove:    true,
		PullParent:     true,
	}

	Tar(GetCurrentDirectory(), "repo.tar")
	dockerBuildContext, err := os.Open("repo.tar")
	defer dockerBuildContext.Close()

	if err != nil {
		log.Fatalf("build context failed:%v", err)
	}

	buildResponse, err := cli.ImageBuild(context.Background(), dockerBuildContext, buildOptions)
	if err != nil {
		log.Fatalf("buildImage=%s failed:%v", imageName, err)
	}

	defer buildResponse.Body.Close()
	io.Copy(os.Stdout, buildResponse.Body)


	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	r, err := cli.ImagePush(context.Background(), imageName, types.ImagePushOptions{RegistryAuth: authStr})
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, r)
	defer r.Close()

}
