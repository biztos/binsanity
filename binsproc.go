// binsproc.go -- binsanity main processing func

package binsanity

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Result is returned by Process and records the number of files and total
// bytes processed.
type Result struct {
	Files int
	Bytes int
}

// String returns the pretty-print version of Result.
func (r *Result) String() string {
	return fmt.Sprintf("files: %d, bytes: %d", r.Files, r.Bytes)
}

// Process converts all files in dir into readable data in a Go file
// file belonging to package pkg, within the module path mod ("main"
// being special); and writes tests for the generated code.
//
// Note that mod is the module to be *imported* in the test file, while
// pkg is the package name to be *declared* in the source file.
//
// The file defaults to "binsanity.go" and if not empty, must have ".go"
// as its extension.
//
// The values of pkg and mod are guessed by default, assuming a Go Module
// environment.  If this fails an error is returned.
//
// The test file is named "binsanity_test.go" or the equivalent for the
// specified file, and provides full coverage of the generated functions.
//
// If bintesting is true, simple tests for the content are included.  This
// will not affect coverage but it will slow testing down.  Its benefit is
// some small protection against accidental direct edits to the generated
// files (presumably one would more easily notice a change to the test file!)
// and, arguably, against pedantic unit-test fanatics like the author. ;-)
//
// If either file exists it is overwritten.
//
// Paths are stripped of their prefixes up to the dir and converted to
// slash format when stored as asset names.
//
// The first error encountered is returned.
func Process(dir, pkg, mod, file string, bintesting bool) (*Result, error) {

	var err error
	if dir == "" {
		// Don't just hoover up whatever's here, it has to be explicit.
		// (If you *want* do `binsanity /tmp` then fine, but say so.)
		return nil, errors.New("Source file not specified.")
	}
	if file == "" {
		file = "binsanity.go"
	}
	if filepath.Ext(file) != ".go" {
		return nil, errors.New("Output must be to a .go file.")
	}
	if pkg == "" {
		pkg, err = FindPackage(file)
		if err != nil {
			return nil, err
		}
	}
	if mod == "" {
		mod, err = FindModule(file)
		if err != nil {
			return nil, err
		}
	}

	info, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("Asset dir: %v", err)
	}
	if !info.IsDir() {
		// Any case where we want to allow a single file? Maybe.  It's a very
		// bad practice though, so not supporting it unless asked.
		return nil, fmt.Errorf("Not a directory: %s", dir)
	}

	// Inputs check out, let's go.
	res := &Result{}
	fileData := map[string][]byte{}

	// On the one hand, it's silly to not get the paths first in order and
	// then get the content for each path one by one so we don't have to
	// fill up with the memory of every damned thing.  On the other hand,
	// that's exactly what the generated source file is going to make you do
	// every single time, so we should be even.
	walker := func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		// It might be a link to a dir, or something missing...
		realInfo, err := os.Stat(path)
		if err != nil {
			return err
		}
		if realInfo.IsDir() {
			return nil
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("Error reading %s: %v", path, err)
		}
		path = filepath.ToSlash(strings.TrimPrefix(path, dir))
		path = strings.TrimPrefix(path, string(filepath.Separator))
		fileData[path] = b
		return nil
	}

	err = filepath.Walk(dir, walker)
	if err != nil {
		return nil, fmt.Errorf("Error walking %s: %v", dir, err)
	}

	// We want sorted names for consistency, we have no idea what order the
	// file system is going to use.  Also, it gives us a neat way to do the
	// lookups.
	names := make([]string, len(fileData))
	idx := 0
	for k := range fileData {
		names[idx] = k
		idx++
	}
	sort.Strings(names)

	// Create the package file.
	basename := filepath.Base(file)
	pf := fmt.Sprintf(fileStart, basename, pkg)
	for _, n := range names {
		pf += fmt.Sprintf("\t%q,\n", n)
	}
	pf += fileMid
	for _, n := range names {

		var buf bytes.Buffer
		writer := gzip.NewWriter(&buf)

		// TODO: figure out how to test this stuff.  Gonna need to wrap the
		// whole effing FS somehow right?  Because we need to be able to say
		// e.g. first read succeeds, second read fails, or even third.
		if _, err := writer.Write(fileData[n]); err != nil {
			return nil, fmt.Errorf("Error compressing asset %s: %v", n, err)
		}

		if err := writer.Close(); err != nil {
			return nil, err
		}

		encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

		pf += fmt.Sprintf("\t%q,\n", encoded)
	}
	pf += fmt.Sprintf(fileEnd, basename)

	// ...and write it.  (Doing it this way is obviously not very memory
	// efficient, but it makes testing this package fairly easy.)
	if err := os.WriteFile(file, []byte(pf), os.ModePerm); err != nil {
		return nil, fmt.Errorf("Error writing source file to %s: %v", file, err)
	}

	// Create the test file and write it:
	// REWORK ONCE EXAMPLE AT 100!
	tfile := strings.TrimSuffix(file, ".go") + "_test.go"
	tf := fmt.Sprintf("// %s - auto-generated; edit at thine own peril!\n",
		filepath.Base(tfile))
	tf += "//\n// More info: https://github.com/biztos/binsanity\n\n"
	tf += fmt.Sprintf("package %s_test\n\n", pkg)
	tf += "import (\n"
	tf += "\t\"fmt\"\n"
	tf += "\t\"testing\"\n"
	tf += "\n"
	tf += fmt.Sprintf("\t\"%s\"\n", mod)
	tf += ")\n"
	tf += fmt.Sprintf(testFuncs, pkg, pkg, pkg, pkg, pkg, pkg, pkg, pkg, pkg, pkg)
	if err := os.WriteFile(tfile, []byte(tf), os.ModePerm); err != nil {
		return nil, fmt.Errorf("Error writing test file to %s: %v", tfile, err)
	}

	// Done... pending bug reports, of course, which are sort of inevitable
	// for something this hastily written.
	return res, nil

}

const fileStart = `/* binsanity.go - auto-generated; edit at thine own peril!

More info: https://github.com/biztos/binsanity

*/

package %s

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"sort"

	%q
)

// Asset returns the byte content of the asset for the given name, or an error
// if no such asset is available.
func Asset(name string) ([]byte, error) {

	_, found := binsanity_cache[name]
	if !found {
		i := sort.SearchStrings(binsanity_names, name)
		if i == len(binsanity_names) || binsanity_names[i] != name {
			return nil, errors.New("Asset not found.")
		}

		// Not cached, so decode and cache it.
		decoded, err := base64.StdEncoding.DecodeString(binsanity_data[i])
		if err != nil {
			return nil, err
		}
		buf := bytes.NewReader(decoded)
		gzr, err := gzip.NewReader(buf)
		if err != nil {
			return nil, err
		}
		data, err := ioutil.ReadAll(gzr)
		if err != nil {
			return nil, err
		}
		binsanity_cache[name] = data
	}
	return binsanity_cache[name], nil

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

// MustAssetString returns the string content of the asset for the given name,
// or panics if no such asset is available.  This is a convenience function
// for string(MustAsset(name)).
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetNames returns the sorted names of the assets.
func AssetNames() []string {
	return binsanity_names
}

// this must remain sorted or everything breaks!
var binsanity_names = []string{
`

const fileMid = `}

// only decode once per asset.
var binsanity_cache = map[string][]byte{}

// assets are gzipped and base64 encoded. WARNING: long lines may follow!
var binsanity_data = []string{
`

const fileEnd = `}
// END OF %s
`

// NOTE: this requires 7 substitutions in Sprintf, all of the package name.
// NOTE: arguably pretty useful to test the actual data, maybe by checksum?
const testFuncs = `
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

func TestMustAssetString(t *testing.T) {

	// Not found:
	missing := "---* no such asset we certainly hope *---"
	exp := "Asset ---* no such asset we certainly hope *--- not found"
	panicky := func() { %s.MustAssetString(missing) }
	AssertPanicsWith(t, panicky, exp, "panic for not found")

	// Found (each one):
	for _, name := range %s.AssetNames() {
		// NOTE: it would be nice to test the non-zero sizes but it's possible
		// to have empty files, so...
		_ = %s.MustAssetString(name)
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
