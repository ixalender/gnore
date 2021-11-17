<p align="center"><img src="hand.svg" style="width: 70px;" width="70"></p>

# Gnore
The command-line tool allows getting an appropriate .gitignore file from github's [gitignore repo](https://github.com/github/gitignore.git).

## Usage

Update available templates
```
gnore update
```

List available templates
```
gnore list
```

Add .gitignore file to the destination directory by the specific template name. If **.gitignore** exists, **gnore** will add specified template to the end of **.gitignore** file.
```
gnore get <template> <dest>
```

Examples
```
gnore get python .

gnore get vue ./some-project-folder
```
## License
MIT

