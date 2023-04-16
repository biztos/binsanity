// main_test.go -- because main wants 100% test coverage too.
//
// (Seriously? Seriously. Design for testing.)

package main

import (
	"bytes"
	"testing"

	"github.com/biztos/binsanity"
)

func TestMainVersion(t *testing.T) {

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	binsanity.OutWriter = stdout
	binsanity.ErrWriter = stderr

	binsanity.Args = []string{"appname", "--version"}

	main()

	got_out := stdout.String()
	got_err := stderr.String()
	exp_out := "binsanity version v1.0.0\n"
	exp_err := ""
	if got_out != exp_out {
		t.Errorf("stdout:\ngot: %s\nexp: %s", got_out, exp_out)
	}
	if got_err != exp_err {
		t.Errorf("stdout:\ngot: %s\nexp: %s", got_err, exp_err)
	}

}
