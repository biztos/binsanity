// binsanity.go -- like bindoc but stupider and sane.

// Package binsanity encodes files into Go source, with testing.
//
// WARNING -- PRE-ALPHA SOFTWARE
//
// ** This package probably isn't working yet. **
//
// Inspired by the bindata package, binsanity aims to provide a minimally
// useful subset of features while also enabiling proper testing of the
// generated Go source code.
//
// For a much more featureful, but less testable approach see:
//
// https://godoc.org/github.com/jteeuwen/go-bindata
//
// One generally does not use this package directly, but rather the command
// binsanity.
//
// Differences From Bindata
//
// * Data is not compressed.
//
// * Only the AssetNames, Asset, and MustAsset functions are implemented.
//
// * Edge cases, probably numerous, have not been much considered.
//
// * Your test coverage will not be reduced. :-)
package binsanity

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Process converts all files in srcdir into accessible data in a Go file
// destfile belonging to pkg, imported as importpath.  The destfile
// must be an empty string, or have ".go" as its extension.  If it is the
// empty string then the default will be used: "binsanity.go" in the current
// directory (this is usually the desired behavior).  A test file
// ("filename_test.go") is also generated in the same directory.  If either
// file exists it is overwritten.
//
// If pkg is the empty string, the first "package" statement in the first
// go file in the destination directory is used to define the package; if none
// is found then the package is assumed to be "main".
//
// If pkg is "main" then no test file is written.
//
// If importpath is the empty string and pkg is not "main" then a guess is
// made based on the assumption that destfile's directory is a package
// directory in a standard Go source directory, e.g. $GOPATH/src. If this
// fails an error is returned.
//
// Paths are stripped of their prefixes up to the srcdir and converted to
// slash format when stored as asset names.
//
// The first error encountered is returned.
func Process(srcdir, pkg, importpath, destfile string) error {

	if destfile == "" {
		destfile = "binsanity.go"
	}
	if filepath.Ext(destfile) != ".go" {
		return fmt.Errorf("Destination file does not have .go extension: %s",
			destfile)
	}
	if srcdir == "" {
		srcdir = "."
	}

	info, err := os.Stat(srcdir)
	if err != nil {
		return fmt.Errorf("Source dir: %v", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("Source not a directory: %s", srcdir)
	}

	fileData := map[string][]byte{}

	walker := func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("Error reading %s: %v", path, err)
		}
		path = filepath.ToSlash(strings.TrimPrefix(path, srcdir))
		fileData[path] = b
		return nil
	}

	err = filepath.Walk(srcdir, walker)
	if err != nil {
		return fmt.Errorf("Error walking %s: %v", srcdir, err)
	}

	names := make([]string, len(fileData))
	idx := 0
	for k, _ := range fileData {
		names[idx] = k
		idx++
	}
	sort.Strings(names)

	pkg, err = getPkg(pkg, destfile)
	if err != nil {
		return fmt.Errorf("Error finding package name: %v", err)
	}
	importpath, err = getImportPath(importpath, destfile)
	if err != nil {
		return fmt.Errorf("Error finding import path: %v", err)
	}

	// Create the package file.
	pf := fmt.Sprintf("// %s - auto-generated; edit at thine own peril!\n",
		filepath.Base(destfile))
	pf += "//\n// More info: https://github.com/biztos/binsanity\n\n"
	pf += fmt.Sprintf("package %s\n\n", pkg)
	pf += "import \"fmt\"\n"
	pf += pkgFuncs
	pf += "// The names, sorted:\n"
	pf += "var names = []string{\n"
	for _, n := range names {
		pf += fmt.Sprintf("\t\"%s\",\n", n)
	}
	pf += "}\n\n"
	pf += "// The data itself (long lines ahead):\n"
	pf += "var data = map[string][]byte{\n"
	for _, n := range names {
		pf += fmt.Sprintf("\t\"%s\": %#v,\n", n, fileData[n])
	}
	pf += "}\n"

	// ...and write it.  (Doing it this way is obviously not very memory
	// efficient, but it makes testing this package fairly easy.)
	if err := ioutil.WriteFile(destfile, []byte(pf), os.ModePerm); err != nil {
		return fmt.Errorf("Error writing package file to %s: %v", destfile, err)
	}

	// No tests for main, it's not usually done that way.
	if pkg == "main" {
		return nil
	}

	// Create the test file and write it:
	tfile := strings.TrimSuffix(destfile, ".go") + "_test.go"
	tf := fmt.Sprintf("// %s - auto-generated; edit at thine own peril!\n",
		filepath.Base(tfile))
	tf += "//\n// More info: https://github.com/biztos/binsanity\n\n"
	tf += fmt.Sprintf("package %s_test\n\n", pkg)
	tf += "import (\n"
	tf += "\t\"fmt\"\n"
	tf += "\t\"testing\"\n"
	tf += "\n"
	tf += fmt.Sprintf("\t\"%s\"\n", importpath)
	tf += ")\n"
	tf += fmt.Sprintf(testFuncs, pkg, pkg, pkg, pkg, pkg, pkg, pkg)
	if err := ioutil.WriteFile(tfile, []byte(tf), os.ModePerm); err != nil {
		return fmt.Errorf("Error writing test file to %s: %v", tfile, err)
	}

	// Done... pending bug reports, of course, which are sort of inevitable
	// for something this hastily written.
	return nil

}

func getPkg(pkg, destfile string) (string, error) {

	if pkg != "" {
		return pkg, nil
	}

	dir := filepath.Dir(destfile)

	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}
	for _, info := range infos {
		file := info.Name()
		if !strings.HasSuffix(file, ".go") {
			continue
		}
		b, err := ioutil.ReadFile(filepath.Join(dir, file))
		if err != nil {
			return "", err
		}
		lines := strings.Split(string(b), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "package ") {
				pkg = strings.TrimSpace(strings.TrimPrefix(line, "package "))
				return pkg, nil
			}
		}
	}

	return "main", nil

}

func getImportPath(importpath, destfile string) (string, error) {

	if importpath != "" {
		return importpath, nil
	}

	abspath, _ := filepath.Abs(filepath.Dir(destfile))

	// Ideally we might just have this in our GOPATH (which can be many).
	gopaths := filepath.SplitList(os.Getenv("GOPATH"))
	for _, p := range gopaths {
		agp, _ := filepath.Abs(p)
		src := filepath.Join(agp, "src")
		if strings.HasPrefix(abspath, src) {
			ipath := strings.TrimPrefix(strings.TrimPrefix(abspath, src),
				string(filepath.Separator))
			return ipath, nil
		}
	}

	// Failing the GOPATH solution, what is worth trying?  We could walk
	// up the directory tree looking for "src" but then why would that
	// not be in your GOPATH?  It doesn't seem like there's a sane option
	// here that will not be wrong a lot, and sometimes-wrong is worse than
	// a predictable error.
	return "", fmt.Errorf("Path %s not found under $GOPATH %s",
		abspath, os.Getenv("GOPATH"))

}

var pkgFuncs = `
// Asset returns the byte content of the asset for the given name, or an error
// if no such asset is available.
func Asset(name string) ([]byte, error) {
	if b := data[name]; b != nil {
		return b, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset returns the byte content of the asset for the given name, or
// panics if no such asset is available.
func MustAsset(name string) []byte {
	b, err := Asset(name)
	if err != nil {
		panic(err.Error())
	}
	return b
}

// AssetNames returns the alpha-sorted names of the assets.
func AssetNames() []string {
	return names
}

`

// NOTE: this requires 7 substitutions in Sprintf, all of the pkg.
var testFuncs = `
func TestAssetNames(t *testing.T) {
	names := %s.AssetNames()
	t.Log(names)
}

func TestAsset(t *testing.T) {

	// Not found:
	missing := "---* no such asset we certainly hope *---"
	_, err := %s.Asset(missing)
	if err == nil {
		t.Fatal("No error for missing asset.")
	}
	if err.Error() != "Asset "+missing+" not found" {
		t.Fatal("Wrong error for missing asset: ", err.Error())
	}

	// Found (each one):
	for _, name := range %s.AssetNames() {
		// NOTE: it would be nice to test the non-zero sizes but it's possible
		// to have empty files, so...
		_, err := %s.Asset(name)
		if err != nil {
			t.Fatal(err.Error())
		}
	}
}

func TestMustAsset(t *testing.T) {

	// Not found:
	missing := "---* no such asset we certainly hope *---"
	exp := "Asset ---* no such asset we certainly hope *--- not found"
	panicky := func() { %s.MustAsset(missing) }
	AssertPanicsWith(t, panicky, exp, "panic for not found")

	// Found (each one):
	for _, name := range %s.AssetNames() {
		// NOTE: it would be nice to test the non-zero sizes but it's possible
		// to have empty files, so...
		_ = %s.MustAsset(name)
	}
}

// For a more useful version of this see: https://github.com/biztos/testig
func AssertPanicsWith(t *testing.T, f func(), exp string, msg string) {

	panicked := false
	got := ""
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
				got = fmt.Sprintf("%%s", r)
			}
		}()
		f()
	}()

	if !panicked {
		t.Fatalf("Function did not panic: %%s", msg)
	} else if got != exp {

		t.Fatalf("Panic not as expected: %%s\n  expected: %%s\n    actual: %%s",
			msg, exp, got)
	}

	// (In go testing, success is silent.)

}
`
