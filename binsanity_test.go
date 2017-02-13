// binsanity_test.go - tests for binsanity
//
// NOTE: not using testify here in order to keep out non-core deps, but
// consequently the testing is a little weak.
package binsanity_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/biztos/binsanity"
)

func Test_Process_Basic(t *testing.T) {
	dir, err := ioutil.TempDir("test", "binsanity_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	pkg := "bstest"
	src := filepath.Join("test", "files")
	destfile := filepath.Join(dir, "bstest")

	// We don't actually know our import path so we might as well guess.
	// (No idea how this will play with CI, but Travis will let us know.)
	err = binsanity.Process(src, pkg, "", destfile)
	if err != nil {
		t.Fatal(err)
	}

}
