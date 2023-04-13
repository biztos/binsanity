// binsproc_test.go - tests for stuff in binsproc.go
package binsanity_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/biztos/binsanity"
)

func TestResultString(t *testing.T) {

	assert := assert.New(t)

	res := &binsanity.Result{
		Files: 1234,
		Bytes: 5678,
	}
	assert.Equal("files: 1234, bytes: 5678", res.String())

}

func TestProcessErrNoAssetDir(t *testing.T) {

	assert := assert.New(t)

	cfg := &binsanity.Config{}

	_, err := binsanity.Process(cfg)
	assert.ErrorContains(err, "Source file not specified.")

}

func TestProcessErrFileNotGo(t *testing.T) {

	assert := assert.New(t)

	cfg := &binsanity.Config{Dir: ExampleAssetDir, File: "foo.txt"}

	_, err := binsanity.Process(cfg)
	assert.ErrorContains(err, "Output must be to a .go file.")

}

func TestProcessErrAssetDirNotExist(t *testing.T) {

	assert := assert.New(t)

	cfg := &binsanity.Config{Dir: ExampleAssetDir + "nopers", File: "foo.go"}

	_, err := binsanity.Process(cfg)
	assert.ErrorContains(err, "Asset dir")

}

func TestProcessErrAssetDirNotDir(t *testing.T) {

	assert := assert.New(t)

	cfg := &binsanity.Config{
		Dir:  filepath.Join(ExampleDir, "main.go"),
		File: "foo.go",
	}

	_, err := binsanity.Process(cfg)
	assert.ErrorContains(err, "ot a directory")

}
