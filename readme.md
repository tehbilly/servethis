`servethis` is a simple cross-platform file serving utility written in [go](http://golang.org/). It provides a quick and easy way to share any files on your filesystem via http.

## Installation
Assuming you have go [installed](http://golang.org/doc/install/) you can simply run `go get github.com/tehbilly/servethis`.

## Usage
Run `servethis -h` to see the plethora of options available to you!.

Without arguments `servethis` will serve the contents of the current working directory on port 8000. To change either the `path` that's served or the `port` that data is served on use the flags with the same names. For example: `servethis -path c:/share/stuff -port 9191`.