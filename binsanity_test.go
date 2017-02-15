// binsanity_test.go - tests for binsanity
//
// NOTE: not using testify here in order to keep out non-core deps, but
// consequently the testing is a little weak.
package binsanity_test

import (
	"bytes"
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

func Test_Process_ExplicitMain(t *testing.T) {

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

// Of course I find a bug the first time I run the thing outside my dev host.
// A symlink to a dir looks like a file, but if you try to read it (the dir)
// as a file you'll get an error.
func Test_Process_SymLinkDir(t *testing.T) {

	// Our temp dir:
	// (TODO: $ENV option to not clean up.)
	dir, err := ioutil.TempDir("test", "binsanity_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Make another dir inside it:
	dir2, err := ioutil.TempDir(dir, "target_")
	if err != nil {
		t.Fatal(err)
	}
	linkSrc, err := filepath.Abs(dir2)
	if err != nil {
		t.Fatal(err)
	}

	// And symlink that:
	if err := os.Symlink(linkSrc, filepath.Join(dir, "linky")); err != nil {
		t.Fatal(err)
	}

	// It's OK to have the source and dest be in the same place, although
	// it's probably not something you'd do in real life.
	src := dir
	destfile := filepath.Join(dir, "bstest.go")

	// Now we can run it and it should give an error until the bug is fixed,
	// after which it should not.
	//
	// Something like this:
	//     Error walking test/binsanity_test_713708718: Error reading
	//     test/binsanity_test_713708718/linky: read
	//     test/binsanity_test_713708718/linky: is a directory
	err = binsanity.Process(src, "bstest", "", destfile)
	if err != nil {
		t.Fatal(err)
	}

	// Still here? Let's run the tests!
	RunTestsInDir(t, dir)

}

// Another goodie: if the symlink is not found we should get a useful
// error, not a random one.
func Test_Process_SymLinkNotFound(t *testing.T) {

	// Our temp dir:
	// (TODO: $ENV option to not clean up.)
	dir, err := ioutil.TempDir("test", "binsanity_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Symlink something that will not be found:
	if err := os.Symlink("NOTHING", filepath.Join(dir, "linky")); err != nil {
		t.Fatal(err)
	}

	// It's OK to have the source and dest be in the same place, although
	// it's probably not something you'd do in real life.
	src := dir
	destfile := filepath.Join(dir, "bstest.go")

	// Now we can run it and it should give a useful error.
	err = binsanity.Process(src, "bstest", "", destfile)
	if err == nil {
		t.Fatal("No error returned for missing symlink.")
	}
	if !strings.HasSuffix(err.Error(), "no such file or directory") {
		t.Fatalf("Error not as expected: %v", err)
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

func Test_ParseArgs_Version(t *testing.T) {

	var out bytes.Buffer
	var err bytes.Buffer
	exited := -1
	exit := func(c int) { exited = c }
	args := []string{"p", "-v"}

	binsanity.ParseArgs(args, exit, &out, &err)
	if exited != 0 {
		t.Fatalf("Wrong exit code: %d", exited)
	}

	if errs := err.String(); errs != "" {
		t.Fatalf("Got standard error: %s", errs)
	}
	outs := out.String()
	if outs != "binsanity version 0.1.0\n" {
		t.Fatalf("Got unexpected standard output: %s", outs)
	}
}

func Test_ParseArgs_Help(t *testing.T) {

	var out bytes.Buffer
	var err bytes.Buffer
	exited := -1
	exit := func(c int) { exited = c }
	args := []string{"p", "-h"}

	binsanity.ParseArgs(args, exit, &out, &err)
	if exited != 0 {
		t.Fatalf("Wrong exit code: %d", exited)
	}

	if errs := err.String(); errs != "" {
		t.Fatalf("Got standard error: %s", errs)
	}

	// Let's not actually test the help text content right now, hm?
	outs := out.String()
	if outs == "" {
		t.Fatal("Got nothing on standard output.")
	}
}

func Test_ParseArgs_NoArgs(t *testing.T) {

	var out bytes.Buffer
	var err bytes.Buffer
	exited := -1
	exit := func(c int) { exited = c }
	args := []string{"p"}

	binsanity.ParseArgs(args, exit, &out, &err)
	if exited != 1 {
		t.Fatalf("Wrong exit code: %d", exited)
	}
	if outs := out.String(); outs != "" {
		t.Fatalf("Got standard output: %s", outs)
	}
	errs := err.String()
	if errs != "Source dir required. Usage: binsanity [options] ASSETDIR\n" {
		t.Fatalf("Got unexpected standard error: %s", errs)
	}
}

func Test_ParseArgs_BadArg(t *testing.T) {

	var out bytes.Buffer
	var err bytes.Buffer
	exited := -1
	exit := func(c int) { exited = c }
	args := []string{"p", "--nope"}

	binsanity.ParseArgs(args, exit, &out, &err)
	if exited != 1 {
		t.Fatalf("Wrong exit code: %d", exited)
	}
	if outs := out.String(); outs != "" {
		t.Fatalf("Got standard output: %s", outs)
	}
	errs := err.String()
	if errs != "Bad option: --nope. Usage: binsanity [options] ASSETDIR\n" {
		t.Fatalf("Got unexpected standard error: %s", errs)
	}
}

func Test_ParseArgs_TwoDirs(t *testing.T) {

	var out bytes.Buffer
	var err bytes.Buffer
	exited := -1
	exit := func(c int) { exited = c }
	args := []string{"p", "first", "second"}

	binsanity.ParseArgs(args, exit, &out, &err)
	if exited != 1 {
		t.Fatalf("Wrong exit code: %d", exited)
	}
	if outs := out.String(); outs != "" {
		t.Fatalf("Got standard output: %s", outs)
	}
	errs := err.String()
	if errs != "Too many source dirs. Usage: binsanity [options] ASSETDIR\n" {
		t.Fatalf("Got unexpected standard error: %s", errs)
	}
}

func Test_ParseArgs_AllSet(t *testing.T) {

	var out bytes.Buffer
	var err bytes.Buffer
	exited := -1
	exit := func(c int) { exited = c }
	args := []string{"p", "--import=imp", "--package=pkg", "--output=out",
		"srcdir"}

	res := binsanity.ParseArgs(args, exit, &out, &err)
	if exited != -1 {
		t.Fatalf("Wrong exit code: %d", exited)
	}
	if outs := out.String(); outs != "" {
		t.Fatalf("Got standard output: %s", outs)
	}
	if errs := err.String(); errs != "" {
		t.Fatalf("Got standard error: %s", errs)
	}
	if len(res) != 4 {
		t.Fatalf("Got wrong-sized return: %v", res)
	}
	if res[0] != "srcdir" {
		t.Fatalf("Got wrong srcdir to pass to Process: %s", res[0])
	}
	if res[1] != "pkg" {
		t.Fatalf("Got wrong package to pass to Process: %s", res[1])
	}
	if res[2] != "imp" {
		t.Fatalf("Got wrong import to pass to Process: %s", res[2])
	}
	if res[3] != "out" {
		t.Fatalf("Got wrong destfile to pass to Process: %s", res[3])
	}
}

func Test_ParseArgs_SrcOnly(t *testing.T) {

	var out bytes.Buffer
	var err bytes.Buffer
	exited := -1
	exit := func(c int) { exited = c }
	args := []string{"p", "srcdir"}

	res := binsanity.ParseArgs(args, exit, &out, &err)
	if exited != -1 {
		t.Fatalf("Wrong exit code: %d", exited)
	}
	if outs := out.String(); outs != "" {
		t.Fatalf("Got standard output: %s", outs)
	}
	if errs := err.String(); errs != "" {
		t.Fatalf("Got standard error: %s", errs)
	}
	if len(res) != 4 {
		t.Fatalf("Got wrong-sized return: %v", res)
	}
	if res[0] != "srcdir" {
		t.Fatalf("Got wrong srcdir to pass to Process: %s", res[0])
	}
	if res[1] != "" || res[2] != "" || res[3] != "" {
		t.Fatalf("Got unexpected data to pass to Process: %v", res)
	}
}
