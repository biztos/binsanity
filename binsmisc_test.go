// binsmisc_test.go -- miscellaneous shared binsanity_test stuff.

package binsanity_test

import (
	"path/filepath"
)

var ExampleDir = filepath.Join("testdata", "example")
var ExampleSubDir = filepath.Join(ExampleDir, "sub")
var ExampleAssetDir = filepath.Join(ExampleDir, "assets")
var NonesuchDir = filepath.Join("testdata", "nope")
var ScanDir = filepath.Join("testdata", "scan")
