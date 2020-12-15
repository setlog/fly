package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/setlog/fly/flags"
	"github.com/setlog/panik"
)

func main() {
	defer panik.ExitTraceTo(os.Stderr)
	f := flags.Parse(os.Args[1:])
	panik.OnError(os.MkdirAll(filepath.FromSlash("src/main/migration"), 0700))
	scriptFilePath := filepath.FromSlash(path.Join("src/main/migration", nextFlywayScriptPrefix("src/main/migration")+f.ScriptPrefix+".sql"))
	panik.OnError(ioutil.WriteFile(scriptFilePath, []byte("USE ${schema_client};\n\n\n"), 0600))
	openInVsCode(scriptFilePath)
}

func nextFlywayScriptPrefix(folderPath string) string {
	major, minor := nextFlywayScriptVersion(folderPath)
	return fmt.Sprintf("V%03d.%03d__", major, minor)
}

func nextFlywayScriptVersion(folderPath string) (major, minor int) {
	largestMajor, largestMinor := 1, 0
	infos, err := ioutil.ReadDir(filepath.FromSlash(folderPath))
	panik.OnError(err)
	for _, info := range infos {
		if !info.IsDir() {
			major, minor, err := getFlywayScriptVersion(info.Name())
			if err == nil {
				if major > largestMajor {
					largestMajor, largestMinor = major, minor
				} else if major == largestMajor && minor > largestMinor {
					largestMinor = minor
				}
			}
		}
	}
	return largestMajor, largestMinor + 1
}

func getFlywayScriptVersion(fileName string) (major, minor int, err error) {
	r := regexp.MustCompile(`^(V|v)([0-9]+)\.([0-9]+)__.+\.sql$`)
	if !r.Match([]byte(fileName)) {
		return 0, 0, fmt.Errorf("regexp mismatch")
	}
	versions := regexp.MustCompile("[0-9]+").FindAllString(fileName, 2)
	return atoi(versions[0]), atoi(versions[1]), nil
}

func atoi(s string) int {
	n, err := strconv.Atoi(s)
	panik.OnError(err)
	return n
}

func openInVsCode(filePath string) {
	cmd := exec.Command("code", filePath)
	defer func() {
		cmd.Process.Release()
	}()
	cmd.Start()
}
