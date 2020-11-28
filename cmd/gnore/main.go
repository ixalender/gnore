package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

const repoDir = "./templates"

const repoURL = "https://github.com/github/gitignore.git"

const usage = `
  Usage:
    gnore update                        update available templates
    gnore list                          list available templates
    gnore get <template> [--path dest]  add .gitignore file by specific template

  Options:
    -p, --path dest   config file to load [default: robo.yml]

  Examples:

    add .gitignore file by template name to the same dir
    $ gnore get python .
`

func main() {
	flag.Parse()

	switch flag.Arg(0) {
	case "list":
		listTemplates()
	case "update":
		fmt.Println("Updateing templates...")
		updateTemplates()
		fmt.Println("Done.")
	case "get":
		if template := flag.Arg(1); template != "" {
			path := flag.String("path", ".", "destination path")
			getTemplate(template, *path)
		}

	default:
		fmt.Println(usage)
	}
}

func listTemplates() (err error) {
	var templates []string
	err = filepath.Walk(repoDir, func(
		path string,
		info os.FileInfo,
		err error) error {

		if filepath.Ext(path) == ".gitignore" {
			templates = append(templates, strings.Split(info.Name(), ".")[0])
		}
		return nil
	})
	checkError(err)

	for _, template := range templates {
		fmt.Println(template)
	}

	return
}

func updateTemplates() (err error) {
	if _, err := os.Stat(repoDir + "/.git"); os.IsNotExist(err) {
		err := os.MkdirAll(repoDir, 0755)
		if err != nil {
			return err
		}
		clone()
	} else {
		pull()
	}

	return
}

func clone() (err error) {
	_, cloneErr := git.PlainClone(repoDir, false, &git.CloneOptions{
		URL:               repoURL,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	checkError(cloneErr)
	return
}

func pull() (err error) {
	r, err := git.PlainOpen(repoDir)
	checkError(err)

	w, err := r.Worktree()
	checkError(err)

	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	checkError(err)
	return
}

func getTemplate(name string, path string) (err error) {
	info("getting template is not implemented...")
	return
}

func info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

func checkError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}
