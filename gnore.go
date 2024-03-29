package gnore

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
)

// TODO: use build -tags to split production and debug template path
const repoFolder = "/usr/local/share/gnore/templates"

const repoURL = "https://github.com/github/gitignore.git"

// ListTemplates shows a list of templates
func ListTemplates() (err error) {
	var templates []string
	repoDir := getRepoDir(repoFolder)

	err = checkRepo(repoDir)
	CheckError(err)

	err = filepath.Walk(repoDir, func(
		path string,
		info os.FileInfo,
		err error) error {

		if filepath.Ext(path) == ".gitignore" {
			templates = append(templates, strings.Split(info.Name(), ".")[0])
		}
		return nil
	})
	CheckError(err)

	sort.Slice(templates, func(i, j int) bool {
		return strings.ToLower(templates[i]) < strings.ToLower(templates[j])
	})
	for _, template := range templates {
		fmt.Println(template)
	}

	return
}

// UpdateTemplates updates templates with clone or pull action
func UpdateTemplates() (err error) {
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

	CheckError(cloneErr)
	return
}

func pull(dest string) (err error) {
	r, err := git.PlainOpen(dest)
	CheckError(err)

	w, err := r.Worktree()
	CheckError(err)

	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	CheckError(err)
	return
}

// GetTemplate puts the template as .gitignore file to the specific path
func GetTemplate(name string, path string) (err error) {
	repoDir := getRepoDir(repoFolder)

	err = checkRepo(repoDir)
	CheckError(err)

	templateFile, err := SearchTemplate(name)
	if err != nil || len(templateFile) == 0 {
		CheckError(fmt.Errorf("There's no template with the name: %s", name))
	}

	filePath := path + "/.gitignore"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		Info("Copy...")
		_, err = copyFile(templateFile, filePath)
		CheckError(err)
	} else {
		Info("Merge...")
		_, err = mergeFiles(templateFile, filePath)
		CheckError(err)
	}

	Info("Done.")

	return
}

// SearchTemplate searches .gitignore template by given name
func SearchTemplate(name string) (string, error) {
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
	CheckError(err)
	return filePath, nil
}

func checkRepo(path string) (err error) {
	if _, err := os.Stat(path + "/.git"); os.IsNotExist(err) {
		return fmt.Errorf("There are no templates yet. Please do 'gnore update'")
	}
	return
}

func getRepoDir(repoFolder string) string {
	return repoFolder
}

func copyFile(src, dest string) (int64, error) {
	_, err := checkFileRegular(src)
	if err != nil {
		return 0, err
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

func mergeFiles(from, to string) (int64, error) {
	_, err := checkFileRegular(from)
	if err != nil {
		return 0, err
	}

	_, err = checkFileRegular(to)
	if err != nil {
		return 0, err
	}

	sourceBuf, err := ioutil.ReadFile(from)
	if err != nil {
		return 0, fmt.Errorf("Error while reading %s", from)
	}

	dest, err := os.OpenFile(to, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}
	defer dest.Close()

	if _, err := dest.WriteString(fmt.Sprintf("\r\n\r\n%s", string(sourceBuf))); err != nil {
		return 0, fmt.Errorf("Error while merging %s to %s", from, to)
	}

	return 1, nil
}

func checkFileRegular(filePath string) (int64, error) {
	fileStat, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}

	if !fileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", filePath)
	}

	return 1, nil
}

// Info prints given information with params
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// CheckError checks given error
func CheckError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}
