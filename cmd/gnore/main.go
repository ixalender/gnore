package main

import (
	"flag"
	"fmt"
	"gnore"
)

// TODO: use build -tags to split production and debug template path
const repoFolder = "/usr/local/share/gnore/templates"

const repoURL = "https://github.com/github/gitignore.git"

const usage = `
  Usage:
    gnore help                        	this usage info
    gnore update                        update available templates
    gnore list                          list available templates
    gnore get <template> <dest>         add .gitignore file to the <path> by the specific <template> name

  Example:

    add .gitignore file by template name to the same dir
    $ gnore get python .
`

func main() {
	flag.Parse()

	switch flag.Arg(0) {
	case "usage":
		fmt.Println(usage)
	case "list":
		gnore.ListTemplates()
	case "update":
		gnore.Info("Updating templates...")
		err := gnore.UpdateTemplates()
		gnore.CheckError(err)
		gnore.Info("Done.")
	case "get":
		if template := flag.Arg(1); template != "" {
			path := parseArg(flag.Args(), 2, ".")
			gnore.GetTemplate(template, path)
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
