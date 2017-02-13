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
	dir, err := ioutil.TempDir("", "binsanity_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	src := filepath.Join("test", "files")
	// err = binsanity.Process(src, "bstest", dir)
	err = binsanity.Process(src, "bstest", "", "test/bstest.go")
	if err != nil {
		t.Fatal(err)
	}

}
