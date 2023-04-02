// binsfind.go -- binsanity package-finding functions

package binsanity

import (
	"fmt"
	"go/scanner"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

// FindPackage attempts to figure out what package we are in, using a brazenly
// insufficient heuristic: the first package declaration in a go source file
// in the target directory; or the name of the directory; or "main" if that
// is not a valid identifier.
//
// The path provided is to the source file that will contain the package
// declaration when generated.
func FindPackage(path string) (string, error) {

	abspath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(abspath)
	files, err := GoFilesBySize(dir)
	if err != nil {
		return "", err
	}

	// We want to scan for the first package only.  If you have multiple
	// package declarations and have not yet cleaned them up, there's no way
	// for us to guess what you mean.
	for _, file := range files {
		pkg, err := ScanForPackage(file)
		if err != nil {
			return "", err
		}
		if pkg != "" {
			return pkg, err
		}
	}

	// Oh well, we tried.  The directory is *probably* a good choice.
	pkg := filepath.Base(dir)
	if !ValidIdent(pkg) {
		pkg = "main"
	}
	return pkg, nil

}

// ValidIdent returns true if s is a valid Go identifier according to
// the spec: https://go.dev/ref/spec#Identifiers
func ValidIdent(s string) bool {
	if len(s) == 0 || !unicode.IsLetter(rune(s[0])) {
		return false
	}
	for _, c := range s {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && c != '_' {
			return false
		}
	}
	return true
}

// GoFilesBySize returns a list of paths for Go source files paths in dir,
// sorted by file size.
func GoFilesBySize(dir string) ([]string, error) {

	// We want to look at files smallest-first because we could have very
	// large files... seeing as we are in the business of generating them!
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	infos := make([]os.FileInfo, 0, len(entries))
	for _, entry := range entries {
		// The DirEntry vs FileInfo problem is kinda necessary, insofar as
		// you don't want to stat every file when reading a dir, but it's not
		// very elegant.  Worth looking for well-tested convenience wrappers,
		// or making our own (ugh, cross-platform testing, NOPE).
		if !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}

		// The super annoying part of this race condition is that we will
		// presumably get nil back (since FileInfo is an interface, you can
		// do that) and so have to check for the race *ANYWAY* even if it's
		// almost impossible to test.  Gonna have to come up with some other
		// mockable function for ReadDir I guess.  Grrr...
		info, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("FS race condition on %s: %s", entry.Name(), err)
		}
		infos = append(infos, info)

		// Searching files smaller than len("package .") is a waste of time,
		// but not worth dealing with Yet Another Test Case for it as it is
		// unlikely to really happen, and doing full reads of big files is
		// much the greater inefficiency.  Cheese Frostman, get some effing
		// sleep...

	}
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Size() < infos[j].Size()
	})
	paths := make([]string, len(infos))
	for idx, info := range infos {
		paths[idx] = filepath.Join(dir, info.Name())
	}

	return paths, nil

}

// ScanForPackage scans the Go source of the file at path and returns the
// package found, if any.  No package found is not an error, as the file may
// be a work in progress.
func ScanForPackage(path string) (string, error) {

	src, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	// Scan it now... some of this stuff is just copied from the docs, it
	// presumably makes sense to people doing actual parsing. :-)
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil /* no error handler */, scanner.ScanComments)

	// Repeated calls to Scan yield the token sequence found in the input.
	want_pkg := false
	for {
		_, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		if want_pkg {
			// This item SHOULD be the package identifier.
			if tok == token.IDENT {
				// TODO: (maybe!) handle weird edge case for stuff like:
				// `package ไก่` which will currently get you "ไก" and
				// the following token will be token.ILLEGAL because the
				// tone mark is not allowed.
				//
				// In practice what will happen is your package will fail to
				// compile because it's named ไก่ and it is not at all clear
				// that naming the package ไก (different word!) will be a
				// good solution.  However: extreme edge case!
				//
				// So strange that Thai is partially supported, presumably
				// many other langs too. I think the root problem is that
				// in some languages, you do not construct Words out of
				// Letters but rather Letters and Marks (etc).
				return lit, nil
			}

			// Well unless it's malformed, that is... but we might yet
			// have one downstream so keep looking.
			want_pkg = false

		}
		if tok == token.PACKAGE {
			// Next item must be the package identifier.
			want_pkg = true
		}
	}

	return "", nil

}

// FindModule attempts to figure out what module we are in, by reading the
// go.mod file in the target directory or a parent.  If none is found, an
// error is returned: we can't make valid tests without an import string.
//
// If the go.mod is in a parent directory, we assume we are in a normal
// subdirectory of that, and join the paths appropriately.
//
// The path provided is to the source file that will contain the package
// declaration when generated.
func FindModule(path string) (string, error) {

	// don't use filepath for joins use regular path!

	return "foo", nil
}
