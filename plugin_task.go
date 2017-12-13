package main

import (
	"os"
	"strings"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	"bufio"
)

func InitIdea() {

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

func InitVscode() {
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



func Setup(moduleName string) {
	path := fmt.Sprintf("./src/%s", moduleName)
	mainPath := fmt.Sprintf("%s/main.go", path)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
		ioutil.WriteFile(mainPath, []byte("package main;\nfunc main()\n{\n}\n"), 0755)
	}
}
