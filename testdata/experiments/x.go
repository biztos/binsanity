// temp x.go test shit

package main

import (
    "errors"
    "fmt"

    "os"
    "path"
    "path/filepath"
    "text/template"
    "time"

    "golang.org/x/mod/modfile"
)

type GenCode struct {
    Timestamp   time.Time
    CodeFile    string
    TestFile    string
    Package     string
    Module      string
    Names       []string
    DataSums    []string
    DataStrings []string
}

func main() {

    gc := &GenCode{
        Timestamp:   time.Now(),
        CodeFile:    "boobers.go",
        TestFile:    "boobers_test.go",
        Package:     "baz",
        Module:      "foo.com/bar/baz",
        Names:       []string{"one", "two"},
        DataStrings: []string{"one-data", "two-data"},
    }

    tmpl, err := template.New("test").Parse(codeFileTemplate)
    if err != nil {
        panic(err)
    }
    err = tmpl.Execute(os.Stdout, gc)
    if err != nil {
        panic(err)
    }

}

var codeFileTemplate = `package {{.Package}}

// gen: {{.Timestamp}}

import (
    "fmt"
    "{{.Module}}"
)

var names = []string{
{{range .Names}}    "{{.}}",
{{end}}}
var data = []string{
{{range .DataStrings}}    "{{.}}",
{{end}}}
}

package {{.Package}}_test

func foo() {
    fmt.Println({{.Package}})
}
`

// FindImportPath attempts to figure out what module we are in, by reading the
// go.mod file in the target directory or a parent, and returns an import path
// usable by the test files, or an error if no go.mod is found.
//
// If the go.mod is in a parent directory, we assume we are in a normal
// subdirectory of that, and join the paths appropriately.
//
// The path provided is to the source file that will contain the package
// declaration when generated. Note that the package name and the import path
// do not need to agree.
func FindImportPath(file string) (string, error) {

    // Find our nearest go.mod if we have one.
    abspath, err := filepath.Abs(file)
    if err != nil {
        // yeah... this is just setting ourselves up for testing misery.
        // how is this ever triggered?! gonna need to look that up in src...
        return "", err
    }
    dir := filepath.Dir(abspath)

    // Want to climb up the tree until we find a go.mod or we hit an error
    // other than "not found."
    var b []byte

    // WHAT IS THE PLATFORM-INDEPENDENT ROOT DIRECTORY??!?!!
    // should be a const in os...
    mdir := dir
    mpath := filepath.Join(mdir, "go.mod")
    for b == nil && mdir != "/" {
        fmt.Println(mpath)
        // Same overhead as Stat anyway for failure case:
        b, err = os.ReadFile(mpath)
        if os.IsNotExist(err) {
            err = nil
            pdir := filepath.Dir(mdir)
            if pdir == mdir {
                // We are at the root directory, nowhere left to go.
                break
            }
            mdir = pdir
            mpath = filepath.Join(mdir, "go.mod")
        } else if err != nil {
            // oops other err.  don't go further.
            break
        } else if b == nil {
            // uh-oh, empty file!
            break
        }
    }
    if b == nil {
        // No luck, bummer.
        return "", errors.New("No go.mod file found.")
    }

    // Now parse that mod file.
    // TODO: verify we are doing the right thing for versions.  It looks like
    // we are, in that Go will figure that foo.com/bar/v5/baz means the "baz"
    // subdir in the foo.com/bar repo at v5.
    mfile, err := modfile.Parse(mpath, b, nil)
    if err != nil {
        return "", err
    }
    mod := mfile.Module.Mod.Path

    // And reassemble based on the non-module path we have left. Not taking
    // any chances with path separators here.
    dir = filepath.ToSlash(dir)
    mdir = filepath.ToSlash(mdir)
    subs := dir[len(mdir):]

    import_path := path.Join(mod, subs)

    return import_path, nil

}
