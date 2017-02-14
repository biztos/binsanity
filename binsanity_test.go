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
	"strings"
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
	destfile := filepath.Join(dir, "bstest.go")

	// We don't actually know our import path so we might as well guess.
	// (No idea how this will play with CI, but Travis will let us know.)
	err = binsanity.Process(src, pkg, "", destfile)
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

func Test_Process_AllDefaults(t *testing.T) {

	// Our temp dir for cleanup:
	// (TODO: $ENV option to not clean up.)
	dir, err := ioutil.TempDir("test", "binsanity_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

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
	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(infos) != 1 {
		t.Fatalf("Wrong number of files found: %#v", infos)
	}
	_, err = os.Stat(filepath.Join(dir, "binsanity.go"))
	if err != nil {
		t.Fatalf("Package file error: %v", err)
	}

	// Check that it actually has the right package name.
	b, err := ioutil.ReadFile(filepath.Join(dir, "binsanity.go"))
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
