// cmd/binsanity/main.go -- the binsanity executable.
//
// NOTE: contrary to my general inclination I am not including anything
// here that's not in the standard library.  Because this might be more
// generally useful than my other work.

// The binsanity program converts files to Go source.
//
//    go get github.com/biztos/binsanity
//    cd go/src/my-project
//    binsanity my-assets
//    go test -cover
package main

import (
	"fmt"
)

func main() {
	fmt.Println("HERE")
}
