// main.go - binsanity example
package main

import (
	"fmt"
)

func main() {
	fmt.Println(MustAssetString("foo"))
	fmt.Println(MustAssetString("bar"))
	fmt.Println(MustAssetString("baz/bat/bloopf"))
	_, err := Asset("doobie")
	fmt.Println("For doobie: ", err)
}
