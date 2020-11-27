package main

import (
    "flag"
    "fmt"
)

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

func main()  {
    flag.Parse()

    switch flag.Arg(0) {
    case "list":
        fmt.Println("list is here")
    case "update":
        fmt.Println("update is here")
    case "get":
        fmt.Println("get is here")
        
        if template := flag.Arg(1); template != "" {
            fmt.Printf("template '%s' is here\n", template)
            path := flag.String("path", ".", "destination path")
            fmt.Printf("dest path is '%s'\n", *path)
        }
        
    default:
        fmt.Println(usage)
    }
}

