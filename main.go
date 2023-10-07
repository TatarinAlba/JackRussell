package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// TODO: Change structure for appropriate data for parsing
type Commit struct {
	commitId, action, scope, body string
}

func generateChangelogForV(ver string, commits []Commit, outputFile *os.File) {
	var outputStr bytes.Buffer
	outputStr.WriteString(fmt.Sprintf("## [%s](http://google.com)\n", ver))
	outputStr.WriteString(fmt.Sprintf("### Bug Fixes\n"))
	for i := len(commits) - 1; i >= 0; i-- {
		if commits[i].action == "fix" {
			if commits[i].scope == "" {
				outputStr.WriteString(fmt.Sprintf("* %s\n", commits[i].body))
			} else {
				outputStr.WriteString(fmt.Sprintf("* **%s**: %s\n", commits[i].scope, commits[i].body))
			}
		}
	}
	outputFile.WriteString(outputStr.String())
}

/*
Function for parsing commits from the given repository from one tag to anoter and
converting it to the list of self-defined structure that have all need information about commit

Parameters:

fromV: string - Tag given in string format from which to start parsing

toV: string - Tag given in string format to which parse information

repoPth: string - Path to the folder with .git folder (repository path)

Returns:

Slice of Commit structure
*/
func parseCommits(toV, fromV, repPth string) []Commit {
	out, err := exec.Command("git", "-C", repPth, "log", "--pretty=oneline", fmt.Sprintf("%s...%s", toV, fromV)).Output()
	if err != nil {
		fmt.Printf("During program execution error occured: %s\n", err)
		os.Exit(1)
	}
	// Parsing commits and inserting into array
	reg := regexp.MustCompile("(.+) (build|ci|docs|feat|fix|perf|refactor|style|test)(\\(.+\\)|): (.+)")
	match := reg.FindAllStringSubmatch(string(out), -1)
	var commits []Commit = make([]Commit, len(match))
	for i, el := range match {
		commits[i] = Commit{el[1], el[2], strings.Trim(el[3], "()"), el[4]}
	}
	return commits
}

func main() {
	var path string
	flag.StringVar(&path, "path", "", "Path to the repo. Absolute or from the binary")
	flag.Parse()
	_, err := os.Stat(path + "/.git")
	if err != nil {
		fmt.Printf("Error: There is no .git folder by the current path\n")
		os.Exit(1)
	}
	outputFile, err := os.Create("CHANGELOG.md")
	if err != nil {
		fmt.Printf("Error: Cannot create the file CHANGELOG.md\n")
		os.Exit(1)
	}
	defer outputFile.Close()
	tagsB, err := exec.Command("git", "-C", path, "tag", "--sort=-creatordate").Output()
	if err != nil {
		fmt.Printf("Error: Durring command execution error is occured %s", err)
	}
	tags := strings.Split(string(tagsB), "\n")
	tags = tags[:len(tags)-1]
	for i := 1; i <= len(tags)-1; i++ {
		commits := parseCommits(tags[i-1], tags[i], path)
		generateChangelogForV(tags[i-1], commits, outputFile)
	}
}
