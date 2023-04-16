// binsmisc_test.go -- miscellaneous shared binsanity_test stuff.

package binsanity_test

import (
	"os"
	"path/filepath"

	"github.com/biztos/binsanity"
)

var ExampleDir = filepath.Join("testdata", "example")
var ExampleSubDir = filepath.Join(ExampleDir, "sub")
var ExampleAssetDir = filepath.Join(ExampleDir, "assets")
var NonesuchDir = filepath.Join("testdata", "nope")
var ScanDir = filepath.Join("testdata", "scan")

func RestoreDefaults() {

	binsanity.FilePathAbs = filepath.Abs
	binsanity.Args = os.Args
	binsanity.ExitFunc = os.Exit
	binsanity.OutWriter = os.Stdout
	binsanity.ErrWriter = os.Stderr

}
