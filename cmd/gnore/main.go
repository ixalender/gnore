package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

const repoFolder = "templates"

const repoURL = "https://github.com/github/gitignore.git"

const usage = `
  Usage:
    gnore update                        update available templates
    gnore list                          list available templates
    gnore get <template> <dest>         add .gitignore file to the <path> by the specific <template> name

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
		fmt.Println("Updating templates...")
		updateTemplates()
		fmt.Println("Done.")
	case "get":
		if template := flag.Arg(1); template != "" {
			path := parseArg(flag.Args(), 2, ".")
			getTemplate(template, path)
		}

	default:
		fmt.Println(usage)
	}
}

func parseArg(args []string, pos int, defaultVal string) string {
	if len(args) > pos {
		return args[pos]
	}
	return defaultVal
}

func listTemplates() (err error) {
	var templates []string
	repoDir := getRepoDir(repoFolder)

	err = checkRepo(repoDir)
	checkError(err)

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
	repoDir := getRepoDir(repoFolder)

	if _, err := os.Stat(repoDir + "/.git"); os.IsNotExist(err) {
		err := os.MkdirAll(repoDir, 0755)
		if err != nil {
			return err
		}
		clone(repoDir)
	} else {
		pull(repoDir)
	}

	return
}

func clone(dest string) (err error) {
	_, cloneErr := git.PlainClone(dest, false, &git.CloneOptions{
		URL:               repoURL,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	checkError(cloneErr)
	return
}

func pull(dest string) (err error) {
	r, err := git.PlainOpen(dest)
	checkError(err)

	w, err := r.Worktree()
	checkError(err)

	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	checkError(err)
	return
}

func getTemplate(name string, path string) (err error) {
	repoDir := getRepoDir(repoFolder)

	templateFile, err := searchTemplate(name)
	if err != nil || len(templateFile) == 0 {
		checkError(fmt.Errorf("There's no template with name: %s", name))
	}

	err = checkRepo(repoDir)
	checkError(err)

	_, err = copyFile(templateFile, path+"/.gitignore")
	checkError(err)

	return
}

func searchTemplate(name string) (string, error) {
	var filePath string
	repoDir := getRepoDir(repoFolder)
	err := filepath.Walk(repoDir, func(
		path string,
		info os.FileInfo,
		err error) error {

		if strings.ToLower(info.Name()) == strings.ToLower(name+".gitignore") {

			filePath = path
			return nil
		}
		return nil
	})
	checkError(err)
	return filePath, nil
}

func checkRepo(path string) (err error) {
	if _, err := os.Stat(path + "/.git"); os.IsNotExist(err) {
		return fmt.Errorf("Here's no any templates yet. Please do 'gnore update'")
	}
	return
}

func getRepoDir(repoFolder string) string {
	selfDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	checkError(err)

	return fmt.Sprintf("%s/%s", selfDir, repoFolder)
}

func copyFile(src, dest string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dest)
	if err != nil {
		return 0, err
	}
	defer destination.Close()

	nBytes, err := io.Copy(destination, source)
	return nBytes, err
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
