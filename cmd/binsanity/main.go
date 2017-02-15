// cmd/binsanity/main.go -- the binsanity executable.
//
// NOTE: contrary to my general inclination I am not including anything
// here that's not in the standard library, as this program's utility
// is pretty broad (if it exists at all).
//
// TODO: more options! Tests!

// The binsanity program converts asset files to Go source.
//
//    go get github.com/biztos/binsanity/...
//    cd go/src/my-project
//    binsanity my-assets # binsanity --help for more options
//    go test -cover
package main

import (
	"github.com/biztos/binsanity"
	"os"
)

var args = os.Args
var exit = os.Exit
var stdout = os.Stdout
var stderr = os.Stderr

func main() {
	binsanity.RunApp(args, exit, stdout, stderr)
}
