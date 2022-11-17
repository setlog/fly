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
	scriptFilePath := filepath.FromSlash(path.Join("src/main/migration", nextFlywayScriptName("src/main/migration", f.VersionIncrementMethod, f.ScriptPrefix+".sql")))
	panik.OnError(ioutil.WriteFile(scriptFilePath, []byte("USE ${"+extractSchemaName(f.ScriptPrefix)+"};\n\n\n"), 0600))
	openInVsCode(scriptFilePath)
}

func nextFlywayScriptName(folderPath string, incrementMethod flags.VersionIncrementMethod, append string) string {
	var major, minor string
	previousScriptFileName, wasMinorIncrement := latestFlywayScriptFileName(folderPath)
	if incrementMethod == flags.IncrementAsk {
		major, minor = nextFlywayScriptVersion(previousScriptFileName, wasMinorIncrement, flags.IncrementMajor)
		option1 := fmt.Sprintf("V%s.%s__%s", major, minor, append)
		major, minor = nextFlywayScriptVersion(previousScriptFileName, wasMinorIncrement, flags.IncrementMinor)
		option2 := fmt.Sprintf("V%s.%s__%s", major, minor, append)
		options := []string{
			option1, option2,
		}
		var title string
		if previousScriptFileName != "" {
			title = fmt.Sprintf("Choose new file name to follow %s:", previousScriptFileName)
		} else {
			title = "Choose name for first script file:"
		}
		return options[selectOption(title, options)]
	} else {
		major, minor = nextFlywayScriptVersion(previousScriptFileName, wasMinorIncrement, incrementMethod)
		return fmt.Sprintf("V%s.%s__%s", major, minor, append)
	}
}

func extractSchemaName(scriptPrefix string) string {
	adminIndex := strings.Index(scriptPrefix, "_admin")
	if adminIndex == 3 {
		return "schema_admin"
	}
	if adminIndex > 3 {
		return "schema_client_admin"
	}
	return "schema_client"
}

func nextFlywayScriptVersion(previousScriptFileName string, wasMinorIncrement bool, incrementMethod flags.VersionIncrementMethod) (major, minor string) {
	incrementMinor := wasMinorIncrement
	if incrementMethod == flags.IncrementMajor {
		incrementMinor = false
	} else if incrementMethod == flags.IncrementMinor {
		incrementMinor = true
	}
	if previousScriptFileName == "" {
		return incrementFlywayScriptVersion("000", "000", incrementMinor)
	}
	prevMajor, prevMinor, _ := getFlywayScriptVersion(previousScriptFileName)
	return incrementFlywayScriptVersion(prevMajor, prevMinor, incrementMinor)
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
