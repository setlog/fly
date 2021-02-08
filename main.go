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
	"strings"

	"github.com/setlog/fly/flags"
	"github.com/setlog/panik"
)

func main() {
	defer panik.ExitTraceTo(os.Stderr)
	f := flags.Parse(os.Args[1:])
	panik.OnError(os.MkdirAll(filepath.FromSlash("src/main/migration"), 0700))
	scriptFilePath := filepath.FromSlash(path.Join("src/main/migration", nextFlywayScriptPrefix("src/main/migration", f.VersionIncrementMethod)+f.ScriptPrefix+".sql"))
	panik.OnError(ioutil.WriteFile(scriptFilePath, []byte("USE ${"+extractSchemaName(f.ScriptPrefix)+"};\n\n\n"), 0600))
	openInVsCode(scriptFilePath)
}

func nextFlywayScriptPrefix(folderPath string, incrementMethod flags.VersionIncrementMethod) string {
	major, minor := nextFlywayScriptVersion(folderPath, incrementMethod)
	return fmt.Sprintf("V%s.%s__", major, minor)
}

func extractSchemaName(scriptPrefix string) string {
	if strings.Index(scriptPrefix, "_admin") == 3 {
		return "schema_admin"
	}
	if strings.HasSuffix(scriptPrefix, "_admin") {
		return "schema_client_admin"
	}
	return "schema_client"
}

func nextFlywayScriptVersion(folderPath string, incrementMethod flags.VersionIncrementMethod) (major, minor string) {
	latestScriptFileName, wasMinorIncrement := latestFlywayScriptFileName(folderPath)
	incrementMinor := wasMinorIncrement
	if incrementMethod == flags.IncrementMajor {
		incrementMinor = false
	} else if incrementMethod == flags.IncrementMinor {
		incrementMinor = true
	}
	if latestScriptFileName == "" {
		return incrementFlywayScriptVersion("000", "000", incrementMinor)
	}
	major, minor, _ = getFlywayScriptVersion(latestScriptFileName)
	return incrementFlywayScriptVersion(major, minor, incrementMinor)
}

func latestFlywayScriptFileName(folderPath string) (scriptFileName string, wasMinorIncrement bool) {
	largestMajor, largestMinor := 1, 0
	infos, err := ioutil.ReadDir(filepath.FromSlash(folderPath))
	panik.OnError(err)
	for _, info := range infos {
		if !info.IsDir() {
			major, minor, err := getFlywayScriptVersion(info.Name())
			if err == nil {
				if atoi(major) > largestMajor {
					largestMajor, largestMinor = atoi(major), atoi(minor)
					scriptFileName = info.Name()
					wasMinorIncrement = false
				} else if atoi(major) == largestMajor && atoi(minor) > largestMinor {
					largestMinor = atoi(minor)
					scriptFileName = info.Name()
					wasMinorIncrement = true
				}
			}
		}
	}
	return scriptFileName, wasMinorIncrement
}

func getFlywayScriptVersion(fileName string) (major, minor string, err error) {
	r := regexp.MustCompile(`^(V|v)([0-9]+)\.([0-9]+)__.+\.sql$`)
	if !r.Match([]byte(fileName)) {
		return "", "", fmt.Errorf("regexp mismatch")
	}
	versions := regexp.MustCompile("[0-9]+").FindAllString(fileName, 2)
	return versions[0], versions[1], nil
}

func incrementFlywayScriptVersion(major, minor string, incrementMinor bool) (newMajor, newMinor string) {
	if incrementMinor {
		return major, fmt.Sprintf(fmt.Sprintf("%%0%dd", len(minor)), atoi(minor)+1)
	}
	return fmt.Sprintf(fmt.Sprintf("%%0%dd", len(major)), atoi(major)+1), fmt.Sprintf(fmt.Sprintf("%%0%dd", len(minor)), 1)
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
