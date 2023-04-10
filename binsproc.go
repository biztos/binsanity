// binsproc.go -- binsanity main processing func

package binsanity

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"
)

const DummyDataString = "H4sIAAAAAAAA/8rP5gIEAAD//30OFtoDAAAA"
const DummyDataSum = "dc51b8c96c2d745df3bd5590d990230a482fd247123599548e0632fdbf97fc22"

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

// GenData holds the data injected into the templates when generating files.
type GenData struct {
	Timestamp         time.Time
	CodeFile          string
	TestFile          string
	Package           string
	Module            string
	Names             []string
	DataSums          []string
	DataStrings       []string
	ExistingAssetName string
	ExistingAssetSum  string
	MissingAssetName  string
	AssetsEmpty       bool
}

// Config holds the values used in Process, in order to avoid confusion.
type Config struct {
	Dir     string
	Package string
	File    string
	Module  string
}

// Process converts all files in cfg.Dir into readable data in a Go file
// cfg.File belonging to package cfg.Package, and writes tests for the
// generated code.  The test file imports cfg.Module.
//
// The file defaults to "binsanity.go" and if not empty, must have ".go"
// as its extension.
//
// The package and module values are guessed if not provided; a standard Go
// Module environment is assumed.
//
// The test file is named "binsanity_test.go" or the equivalent for the
// code file, and provides full coverage of the generated functions.
//
// If either file exists it is overwritten.
//
// Paths are stripped of their prefixes up to the dir and converted to
// slash format when stored as asset names.
//
// In the rare case of *no* assets found in the directory, a single special
// asset is created in order to achieve test coverage.  Its name is randomized
// and should not conflict with any real-world data as it begins with 256
// underscores.
//
// This asset is *not* returned by the AssetNames function.
//
// The first error encountered is returned.
func Process(cfg *Config) (*Result, error) {

	// must.. resist... edit-in-place... temptation... :-)
	dir := cfg.Dir
	mod := cfg.Module
	pkg := cfg.Package
	file := cfg.File

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
		mod, err = FindImportPath(file)
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

	// Grab filenames.
	paths := []string{}
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

		paths = append(paths, path)

		return nil
	}
	err = filepath.Walk(dir, walker)
	if err != nil {
		return nil, fmt.Errorf("Error walking %s: %v", dir, err)
	}
	sort.Strings(paths)

	// Get data for generating the files.
	tfile := file[:len(file)-3] + "_test.go"
	gen := &GenData{
		CodeFile:    filepath.Base(file),
		TestFile:    filepath.Base(tfile),
		Package:     pkg,
		Module:      mod,
		Names:       make([]string, len(paths)),
		DataSums:    make([]string, len(paths)),
		DataStrings: make([]string, len(paths)),
	}
	total_bytes := 0
	for idx, path := range paths {
		// name is cleaned version of path.
		gen.Names[idx] = strings.TrimPrefix(
			filepath.ToSlash(strings.TrimPrefix(path, dir)),
			string(filepath.Separator),
		)
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("Error reading %s: %v", path, err)
		}
		total_bytes += len(b)

		// sum is of raw bytes.
		gen.DataSums[idx] = fmt.Sprintf("%x", sha256.Sum256(b))

		// data is compressed.
		var buf bytes.Buffer
		writer := gzip.NewWriter(&buf)

		// TODO: (as a general task) figure out how to test this crap and
		// write it up on a blog or something.  It will keep coming up, and
		// it's non-idiomatic to ignore the error but also really not obvious
		// how TF you are supposed to test it without doing some very weird
		// gymnastics.  Presumably a "compress" function that takes an
		// interface as its arg, right?  And you write an implementation that
		// blows up, just for testing.  Yay.  Fuck.  So far so good BUT you
		// still need to trap that error unless you make it panic... which I
		// guess is legit in this case, but again not idiomatic... also maybe
		// worth checking the gzip implementation and see if the writer here
		// actually *can* return an error, right?
		if _, err := writer.Write(b); err != nil {
			return nil, fmt.Errorf("Error compressing asset %s: %v", path, err)
		}
		if err := writer.Close(); err != nil {
			return nil, err
		}

		gen.DataStrings[idx] = base64.StdEncoding.EncodeToString(buf.Bytes())

	}

	// Special case for empty assets -- you might want to have an empty set
	// of assets, but we still want test coverage.
	if len(paths) == 0 {

		b := make([]byte, 256)
		rand.Read(b) // more untestable fun...

		name := fmt.Sprintf("%s%x", strings.Repeat("_", 256), b)

		gen.AssetsEmpty = true
		gen.Names = []string{name}
		gen.DataStrings = []string{DummyDataString}
		gen.DataSums = []string{DummyDataSum}

	}

	// Create the code file.
	ctmpl, err := template.New("t").Parse(MustAssetString("code.tmpl"))
	// ctmpl, err := template.New("t").Parse(CTMP)
	if err != nil {
		panic(err) // testable how? probably not at all....
	}
	cwriter, err := os.Create(file)
	if err != nil {
		return nil, err
	}
	defer cwriter.Close()
	err = ctmpl.Execute(cwriter, gen)
	if err != nil {
		panic(err)
	}

	// Some special sauce for the test file:
	test_idx := int(len(gen.Names) / 2)
	gen.ExistingAssetName = gen.Names[test_idx]
	gen.ExistingAssetSum = gen.DataSums[test_idx]
	gen.MissingAssetName = gen.Names[len(gen.Names)-1] + "--NOPE"

	// Create the test file.
	ttmpl, err := template.New("t").Parse(MustAssetString("tests.tmpl"))
	// ttmpl, err := template.New("t").Parse(TTMP)
	if err != nil {
		panic(err) // testable????
	}
	twriter, err := os.Create(tfile)
	if err != nil {
		return nil, err
	}
	defer twriter.Close()
	err = ttmpl.Execute(twriter, gen)
	if err != nil {
		panic(err)
	}

	// Done... pending bug reports, of course, which are sort of inevitable
	// for something this hastily written.
	res := &Result{
		Files: len(paths),
		Bytes: total_bytes,
	}
	return res, nil

}
