// binsapp_test.go - tests for stuff in binsapp.go (just the app wrapper)
package binsanity_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/biztos/binsanity"
)

func TestRunAppErrNoDirArg(t *testing.T) {

	assert := assert.New(t)

	args := []string{}
	exited := false
	exit_code := 0
	exit := func(c int) {
		exited = true
		exit_code = c
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	binsanity.ExitFunc = exit
	binsanity.OutWriter = stdout
	binsanity.ErrWriter = stderr

	binsanity.RunApp(args)
	assert.True(exited, "exited")
	assert.Equal(1, exit_code, "exit 1")
	assert.Equal("", stdout.String(), "stdout")
	assert.Equal("Single arg required: ASSET_DIR\n", stderr.String(), "stderr")

}

func TestRunAppErrDirNotExist(t *testing.T) {

	assert := assert.New(t)

	args := []string{"appname", NonesuchDir}
	exited := false
	exit_code := 0
	exit := func(c int) {
		exited = true
		exit_code = c
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	binsanity.ExitFunc = exit
	binsanity.OutWriter = stdout
	binsanity.ErrWriter = stderr

	binsanity.RunApp(args)
	assert.True(exited, "exited")
	assert.Equal(1, exit_code, "exit 1")
	assert.Equal("", stdout.String(), "stdout")
	assert.Regexp(regexp.MustCompile("Asset dir: .*"), stderr, "stderr")

}

func TestRunAppOkVersion(t *testing.T) {

	// really just testing the cli rigging here because it triggers early,
	// BUT it found me a couple bugs so nice regression to keep.
	assert := assert.New(t)

	args := []string{"appname", "--version"}
	exited := false
	exit := func(c int) {
		exited = true
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	binsanity.ExitFunc = exit
	binsanity.OutWriter = stdout
	binsanity.ErrWriter = stderr

	binsanity.RunApp(args)
	assert.False(exited, "did not exit through func")
	assert.Equal("binsanity version v1.0.0\n", stdout.String(), "stdout")
	assert.Equal("", stderr.String(), "stderr")
}

func TestRunAppOk(t *testing.T) {

	assert := assert.New(t)

	// By passing in the right args we can recreate the example file 1:1 in
	// a different location.
	tdir := t.TempDir()
	file := filepath.Join(tdir, "binsanity.go")
	args := []string{
		"appname",
		"--package=main",
		"--module=biztos.com/example",
		"--output=" + file,
		ExampleAssetDir,
	}

	exited := false
	exit := func(c int) {
		exited = true
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	binsanity.ExitFunc = exit
	binsanity.OutWriter = stdout
	binsanity.ErrWriter = stderr

	binsanity.RunApp(args)
	assert.False(exited, "did not exit through func")
	assert.Equal("files: 3, bytes: 46\n", stdout.String(), "stdout")
	assert.Equal("", stderr.String(), "stderr")

	// Check the code file.
	if assert.FileExists(file) {
		assert.Equal("match",
			difflines(file, filepath.Join(ExampleDir, "binsanity.go")))

	}

	// Check the test file.
	tfile := filepath.Join(tdir, "binsanity_test.go")
	if assert.FileExists(tfile) {
		assert.Equal("match",
			difflines(tfile, filepath.Join(ExampleDir, "binsanity_test.go")))

	}
}

// return a very hacky "diff" of two files by line.  files must exist.
// (not worth the diffmatchpatch complexity for just this)
func difflines(f1, f2 string) string {

	b1, err := os.ReadFile(f1)
	if err != nil {
		panic(err)
	}
	b2, err := os.ReadFile(f2)
	if err != nil {
		panic(err)
	}

	s1 := string(b1)
	s2 := string(b2)
	if s1 == s2 {
		return "match"
	}

	lines1 := strings.Split(s1, "\n") // working on 'doze?
	lines2 := strings.Split(s2, "\n")
	for i, line := range lines1 {
		if len(lines2) < i+1 {
			return fmt.Sprintf("file 1 continues past file 2 at line %d: %s",
				i+1, line)
		}
		if line != lines2[i] {
			return fmt.Sprintf("files diverge at line %d: %s <-!-> %s",
				i+1, line, lines2[i])
		}
	}
	// match so far but #2 could be longer.
	if len(lines2) > len(lines1) {
		return fmt.Sprintf("file 2 continues past file 1 at line %d: %s",
			len(lines1)+1, lines2[len(lines1)])
	}

	return "oops programmer error"
}
