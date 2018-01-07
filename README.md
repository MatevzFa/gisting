# Gisting

`gisting` is a command line program for creating and downloading GitHub Gists.

## Installation
```
go get github.com/MatevzFa/gisting
```

## Usage

```
$ gisting --token $(cat ~/Documents/.gittoken) create -d "Some random gist" --private main.go
https://gist.github.com/dc68f61ef51f6a77fd4febb86ada885a
```

```
$ gisting help
usage: gisting [<flags>] <command> [<args> ...]

Flags:
      --help         Show context-sensitive help (also try --help-long and
                     --help-man).
  -t, --token=TOKEN  OAuth token for accessing the gist API

Commands:
  help [<command>...]
    Show help.

  download <id>
    Download a gist.

  create [<flags>] <files>...
    Create a gist.
```