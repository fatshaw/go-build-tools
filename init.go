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
	"bytes"
	"log"
)

func main() {

	example := [][2]string{
		{
			"go-dep-tools setup [modulename]",
			"set up a new project with name=[modulename]",
		},
		{
			"go-dep-tools init [vscode|idea]",
			"init go environment for vscode or idea",
		},
		{
			"go-dep-tools dep [modulename]",
			"install the project's dependencies",
		},
	}

	if len(os.Args) == 3 {

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

		if os.Args[1] == "dep" {
			initDep(os.Args[2])
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

func initDep(moduleName string) {

	cmd := exec.Command("sh", "-c", fmt.Sprintf("%s", GetCommand(moduleName)))

	out, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Print(err)
		return
	}

	err = cmd.Start()
	if err != nil {
		fmt.Print(err)
		return
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(out)

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("output=%s\n", buf.String())
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
