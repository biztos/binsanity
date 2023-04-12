// binsfind_test.go - tests for stuff in binsfind.go
package binsanity_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/biztos/binsanity"
)

var ExampleDir = filepath.Join("testdata", "example")
var ExampleSubDir = filepath.Join(ExampleDir, "sub")
var ExampleAssetDir = filepath.Join(ExampleDir, "assets")
var NonesuchDir = filepath.Join("testdata", "nope")
var ScanDir = filepath.Join("testdata", "scan")

func TestFindImportPathErrRoot(t *testing.T) {

	assert := assert.New(t)

	// Why is there no filepath.Root?!
	// This *should* work... I think...
	path := filepath.Join(string(filepath.Separator), "foo.go")
	_, err := binsanity.FindImportPath(path)
	assert.ErrorContains(err, "No go.mod file found.")

}

func TestFindImportPathOk(t *testing.T) {

	assert := assert.New(t)
	mod, err := binsanity.FindImportPath(filepath.Join(ExampleSubDir, "foo.go"))
	assert.Nil(err, "no error")
	assert.Equal("biztos.com/example/sub", mod, "path as expected")
}

func TestFindPackageErrReadDir(t *testing.T) {

	assert := assert.New(t)
	_, err := binsanity.FindPackage(filepath.Join(NonesuchDir, "foo.go"))
	assert.IsType(err, &os.PathError{})

}

func TestFindPackageOkTrueMain(t *testing.T) {

	assert := assert.New(t)
	pkg, err := binsanity.FindPackage(filepath.Join(ExampleDir, "foo.go"))
	assert.Nil(err, "no error")
	assert.Equal("main", pkg, "package as expected")
}

func TestFindPackageOkSubDir(t *testing.T) {

	assert := assert.New(t)
	pkg, err := binsanity.FindPackage(filepath.Join(ExampleSubDir, "foo.go"))
	assert.Nil(err, "no error")
	assert.Equal("sub", pkg, "package as expected")
}

func TestFindPackageOkAssetDir(t *testing.T) {

	// the idea here is that if you have a dir with nothing in it (yet) but
	// the binsanity file, it should use the name of the dir as its package.
	assert := assert.New(t)
	pkg, err := binsanity.FindPackage(filepath.Join(ExampleAssetDir, "foo.go"))
	assert.Nil(err, "no error")
	assert.Equal("assets", pkg, "package as expected")
}

func TestFindPackageOkMainFallback(t *testing.T) {

	// obviously don't use the dir if it's not a valid package name!
	spacedir := filepath.Join(ExampleDir, "has space")
	assert := assert.New(t)
	pkg, err := binsanity.FindPackage(filepath.Join(spacedir, "foo.go"))
	assert.Nil(err, "no error")
	assert.Equal("main", pkg, "package as expected")

}

func TestScanForPackagePathErr(t *testing.T) {

	assert := assert.New(t)
	_, err := binsanity.ScanForPackage(NonesuchDir)
	assert.IsType(err, &os.PathError{})

}

func TestScanForPackageOk(t *testing.T) {

	assert := assert.New(t)
	pkg, err := binsanity.ScanForPackage(filepath.Join(ScanDir, "good.go"))
	assert.Nil(err, "no error")
	assert.Equal("good", pkg, "package as expected") // the good place!

}

func TestScanForPackageTricky(t *testing.T) {

	assert := assert.New(t)
	pkg, err := binsanity.ScanForPackage(filepath.Join(ScanDir, "baddish.go"))
	assert.Nil(err, "no error")
	assert.Equal("baddish", pkg, "package as expected")

}

func TestScanForPackageBad(t *testing.T) {

	assert := assert.New(t)
	pkg, err := binsanity.ScanForPackage(filepath.Join(ScanDir, "bad.go"))
	assert.Nil(err, "no error for package not found")
	assert.Equal("", pkg, "package is empty string")

}

func TestGoFilesBySizeDirError(t *testing.T) {

	assert := assert.New(t)
	_, err := binsanity.GoFilesBySize(NonesuchDir)
	assert.IsType(err, &os.PathError{})

}

func TestGoFilesBySizeDirOk(t *testing.T) {

	assert := assert.New(t)
	files, err := binsanity.GoFilesBySize(ExampleDir)
	assert.Nil(err, "no error")
	exp := []string{
		filepath.Join(ExampleDir, "small.go"),
		filepath.Join(ExampleDir, "main.go"),
		filepath.Join(ExampleDir, "binsanity.go"),
		filepath.Join(ExampleDir, "medium.go"),
		filepath.Join(ExampleDir, "big.go"),
	}
	assert.EqualValues(exp, files, "expected files in size order")

}

func TestValidIdentZeroLengthFalse(t *testing.T) {

	assert := assert.New(t)

	assert.False(binsanity.ValidIdent(""))

}

func TestValidIdentNonLetterStartFalse(t *testing.T) {

	assert := assert.New(t)

	assert.False(binsanity.ValidIdent("/this"))

}

func TestValidIdentNonLetterFinishFalse(t *testing.T) {

	assert := assert.New(t)

	assert.False(binsanity.ValidIdent("this/"))

}

func TestValidIdentWeirdButAllowedTrue(t *testing.T) {

	assert := assert.New(t)

	// yay effing ident spec...
	assert.True(binsanity.ValidIdent("ok12ßßมาก"))

}
