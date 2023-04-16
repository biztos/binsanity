// cmd/binsanity/main.go -- the binsanity executable.

// The binsanity program converts asset files to Go source.
//
//	go get github.com/biztos/binsanity/...
//	cd path/to/my-project
//	binsanity my-assets # binsanity --help for more options
//	go test -cover
package main

import (
	"github.com/biztos/binsanity"
)

func main() {
	binsanity.RunApp(binsanity.Args) // <-- os.Args by default
}
