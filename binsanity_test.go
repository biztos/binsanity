// binsanity_test.go - tests for binsanity
//
// NOTE: not using testify here in order to keep out non-core deps, but
// consequently the testing is a little weak.
package binsanity_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/biztos/binsanity"
)

func Test_Process_ErrorNotGoFile(t *testing.T) {

	err := binsanity.Process("any-src", "anypkg", "anyimport", "foo.txt")
	if err == nil {
		t.Fatal("no error returned")
	}
	if err.Error() != "Destination file does not have .go extension: foo.txt" {
		t.Fatalf("Wrong error returned: %v", err)
	}

}

func Test_Process_ErrorNotImportable(t *testing.T) {

	gopath := os.Getenv("GOPATH")
	defer os.Setenv("GOPATH", gopath)
	os.Setenv("GOPATH", "thisisnotyourregulargopath")
	dir, err := ioutil.TempDir("test", "binsanity_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	absdir, _ := filepath.Abs(dir)

	pkg := "bstest"
	src := filepath.Join("test", "files")
	destfile := filepath.Join(dir, "bstest.go")

	err = binsanity.Process(src, pkg, "", destfile)
	if err == nil {
		t.Fatal("no error returned")
	}

	exp := "Error finding import path: Path " + absdir +
		" not found under $GOPATH thisisnotyourregulargopath"

	if err.Error() != exp {
		t.Fatalf("Wrong error returned: %v", err)
	}
}

func Test_Process_ErrorNoSuchSourceDir(t *testing.T) {

	dir, err := ioutil.TempDir("test", "binsanity_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	pkg := "bstest"
	src := filepath.Join("test", "no-such-dir")
	destfile := filepath.Join(dir, "bstest.go")

	err = binsanity.Process(src, pkg, "", destfile)
	if err == nil {
		t.Fatal("no error returned")
	}
	exp := regexp.MustCompile("^Source dir: .* no such file or directory$")
	if !exp.MatchString(err.Error()) {
		t.Fatalf("Wrong error returned: %v", err)
	}

}

func Test_Process_ErrorSourceDirIsFile(t *testing.T) {

	dir, err := ioutil.TempDir("test", "binsanity_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	pkg := "bstest"
	src := "README.md"
	destfile := filepath.Join(dir, "bstest.go")

	err = binsanity.Process(src, pkg, "", destfile)
	if err == nil {
		t.Fatal("no error returned")
	}
	if err.Error() != "Source not a directory: README.md" {
		t.Fatalf("Wrong error returned: %v", err)
	}

}

func Test_Process_Basic(t *testing.T) {
	dir, err := ioutil.TempDir("test", "binsanity_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	pkg := "bstest"
	src := filepath.Join("test", "files")
	destfile := filepath.Join(dir, "bstest.go")

	// Hopefully this works:
	imp := "github.com/biztos/binsanity/test/" + filepath.Base(dir)

	err = binsanity.Process(src, pkg, imp, destfile)
	if err != nil {
		t.Fatal(err)
	}

	pinfo, err := os.Stat(destfile)
	if err != nil {
		t.Fatalf("Package file stat error for %v", err)
	}
	if pinfo.IsDir() {
		t.Fatalf("Package file is a dir... WTF?")
	}
	tinfo, err := os.Stat(filepath.Join(dir, "bstest_test.go"))
	if err != nil {
		t.Fatalf("Package file stat error for %v", err)
	}
	if tinfo.IsDir() {
		t.Fatalf("Test file is a dir... WTF?")
	}

	// Now the really crazy thing is to go in and run the tests.
	if err := RunTestsInDir(t, dir); err != nil {
		t.Fatalf("Bad test results: %v", err)
	}

}

func Test_Process_DefaultMain(t *testing.T) {

	// Our temp dir for cleanup:
	// (TODO: $ENV option to not clean up.)
	dir, err := ioutil.TempDir("test", "binsanity_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	src := filepath.Join("test", "files")
	destfile := filepath.Join(dir, "bstest.go")

	err = binsanity.Process(src, "main", "", destfile)
	if err != nil {
		t.Fatal(err)
	}

	// Let's make sure we have just the one file (no tests for main).
	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(infos) != 1 {
		t.Fatalf("Wrong number of files found: %#v", infos)
	}
	_, err = os.Stat(filepath.Join(dir, "bstest.go"))
	if err != nil {
		t.Fatalf("Package file error: %v", err)
	}

	// Check that it actually has the right package name.
	b, err := ioutil.ReadFile(filepath.Join(dir, "bstest.go"))
	if err != nil {
		t.Fatal(err)
	}
	have := false
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || line[0:2] == "//" {
			continue
		}
		if line == "package main" {
			have = true
			break
		}
		if strings.HasPrefix(line, "package ") {
			t.Fatalf("Wrong package: %s", line)
		}

	}
	if !have {
		t.Fatal("Package not found.")
	}
}

func Test_Process_AllDefaults(t *testing.T) {

	// Our temp dir for cleanup:
	// (TODO: $ENV option to not clean up.)
	dir, err := ioutil.TempDir("test", "binsanity_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	absdir, _ := filepath.Abs(dir)

	// Default dir is current dir, so let's go there.
	cwd, err := os.Getwd()
	if err != nil {
		panic("failed to Getwd: " + err.Error())
	}
	if err := os.Chdir(dir); err != nil {
		panic("failed to Chdir: " + err.Error())
	}
	defer os.Chdir(cwd)

	// While it would probably not be much use, it's totally possible to
	// build from an empty directory.  Use case might be that you're just
	// setting up your workflow but don't have any files yet.
	err = binsanity.Process("", "", "", "")
	if err != nil {
		t.Fatal(err)
	}

	// Let's make sure we have just the one file (no tests for main).
	infos, err := ioutil.ReadDir(absdir)
	if err != nil {
		t.Fatal(err)
	}
	if len(infos) != 1 {
		t.Fatalf("Wrong number of files found: %#v", infos)
	}
	_, err = os.Stat(filepath.Join(absdir, "binsanity.go"))
	if err != nil {
		t.Fatalf("Package file error: %v", err)
	}

	// Check that it actually has the right package name.
	b, err := ioutil.ReadFile(filepath.Join(absdir, "binsanity.go"))
	if err != nil {
		t.Fatal(err)
	}
	have := false
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || line[0:2] == "//" {
			continue
		}
		if line == "package main" {
			have = true
			break
		}
		if strings.HasPrefix(line, "package ") {
			t.Fatalf("Wrong package: %s", line)
		}

	}
	if !have {
		t.Fatal("Package not found.")
	}
}

func Test_Process_FindPackage(t *testing.T) {

	// Our temp dir for cleanup:
	// (TODO: $ENV option to not clean up.)
	dir, err := ioutil.TempDir("test", "binsanity_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	src := filepath.Join("test", "files")
	destfile := filepath.Join(dir, "bstest.go")

	// Put in a file we can look at to get the package name.
	gofile := filepath.Join(dir, "other.go")
	gobytes := []byte("// For instance:\npackage foobar\n")
	if err := ioutil.WriteFile(gofile, gobytes, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	err = binsanity.Process(src, "", "", destfile)
	if err != nil {
		t.Fatal(err)
	}

	// Check that it actually has the right package name.
	b, err := ioutil.ReadFile(filepath.Join(dir, "bstest.go"))
	if err != nil {
		t.Fatal(err)
	}
	have := false
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || line[0:2] == "//" {
			continue
		}
		if line == "package foobar" {
			have = true
			break
		}
		if strings.HasPrefix(line, "package ") {
			t.Fatalf("Wrong package: %s", line)
		}

	}
	if !have {
		t.Fatal("Package not found.")
	}
}

// Run tests with coverage.  Anything other than 100% is an error.
func RunTestsInDir(t *testing.T, dir string) error {

	cwd, err := os.Getwd()
	if err != nil {
		panic("failed to Getwd: " + err.Error())
	}
	if err := os.Chdir(dir); err != nil {
		panic("failed to Chdir: " + err.Error())
	}
	defer os.Chdir(cwd)

	cmd := "go"
	args := []string{"test", "--cover"}
	output, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return err
	}

	// We expect something like this:
	// PASS
	// coverage: 100% of statements
	// ok      github.com/biztos/binsanity    0.010s
	//
	// However we probably don't want to assume too much about newlines. Hack:
	f := strings.Fields(string(output))
	if len(f) < 3 || f[0] != "PASS" {
		return fmt.Errorf("Tests did not PASS: %s", string(output))
	}
	if f[1] != "coverage:" || f[2] != "100.0%" {
		return fmt.Errorf("Coverage seems low: %s", string(output))
	}

	return nil

}
